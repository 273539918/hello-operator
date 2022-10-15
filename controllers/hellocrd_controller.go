/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	demogroupv1 "demo/hello-operator/api/v1"
	resources "demo/hello-operator/pkg/resources"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// HellocrdReconciler reconciles a Hellocrd object
type HellocrdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demogroup.demo,resources=hellocrds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demogroup.demo,resources=hellocrds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demogroup.demo,resources=hellocrds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Hellocrd object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *HellocrdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	// 获取集群中对应crd的cr资源
	//log.Info(fmt.Sprintf("req is : %s", req.String()))
	var helloCrdResource = &demogroupv1.Hellocrd{}
	if err := r.Get(ctx, req.NamespacedName, helloCrdResource); err != nil {
		log.Error(err, "unable to fetch client")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//备份当前的cr信息
	helloCrdResourceOld := helloCrdResource.DeepCopy()

	//如果cr当前的status.hellostatus字段为空，则设置为Pending
	if helloCrdResource.Status.HelloStatus == "" {
		helloCrdResource.Status.HelloStatus = "Pending"
	}

	//根据hellostatus的当前状态执行：
	//1、如果当前状态是Pending： 修改为running，并更新到集群
	//2、如果当前状态是Running:  创建一个pod模版，检查它是否已存在
	//2.1 如果pod不存在：创建pod，重新入队等待处理
	//2.2 如果pod已存在：
	//2.2.1  pod的状态为success或failed，设置cr状态为创建中
	//2.2.2  pod的状态为running，更新cr的lastpodname
	//2.2.3  pod的状态为pending，重新入队等待处理
	//2.3 如果cr的status在执行完上述逻辑之后有变化：用最新的状态覆盖
	switch helloCrdResource.Status.HelloStatus {
	case "Pending":
		helloCrdResource.Status.HelloStatus = "Running"
		err := r.Status().Update(context.TODO(), helloCrdResource)
		if err != nil {
			log.Error(err, fmt.Sprintf("failed to update %s status", req.Name))
			return ctrl.Result{}, err
		} else {
			log.Info(fmt.Sprintf("updated %s status : %s", req.Name, helloCrdResource.Status.HelloStatus))
			//重新入队等待处理
			return ctrl.Result{Requeue: true}, nil
		}
	case "Running":
		//createpod方法：定义了crd创建pod的模版
		pod := resources.CreatePod(helloCrdResource)
		query := &corev1.Pod{}

		err := r.Client.Get(ctx, client.ObjectKey{Namespace: pod.Namespace, Name: pod.ObjectMeta.Name}, query)
		//如果返回pod不存在的错误
		if err != nil && errors.IsNotFound(err) {
			if helloCrdResource.Status.LastPodName == "" {
				err = ctrl.SetControllerReference(helloCrdResource, pod, r.Scheme)
				if err != nil {
					return ctrl.Result{}, err
				}
				//创建Pod
				err = r.Create(context.TODO(), pod)
				if err != nil {
					return ctrl.Result{}, err
				}
				log.Info("pod created successfully", "name", pod.Name)
				return ctrl.Result{Requeue: true}, nil
			}
		} else if err != nil {
			// 找不到pod
			log.Error(err, "cannot get pod")
			return ctrl.Result{}, err
		} else if query.Status.Phase == corev1.PodFailed ||
			query.Status.Phase == corev1.PodSucceeded {
			log.Info("container terminated", "reason", query.Status.Reason, "message", query.Status.Message)
			helloCrdResource.Status.HelloStatus = "Cleaning"
		} else if query.Status.Phase == corev1.PodRunning {
			//log.Info("Client last pod name: " + helloCrdResource.Status.LastPodName)
			if helloCrdResource.Status.LastPodName != helloCrdResource.Spec.ContainerImage+helloCrdResource.Spec.ContainerTag {
				if query.Status.ContainerStatuses[0].Ready {
					log.Info("Container is ready")
					helloCrdResource.Status.LastPodName = helloCrdResource.Spec.ContainerImage + helloCrdResource.Spec.ContainerTag
				} else {
					log.Info("Container not ready")
					return ctrl.Result{Requeue: true}, err
				}

				log.Info("Client last pod name: " + helloCrdResource.Status.LastPodName)
				log.Info("Pod is running.")
			}
		} else if query.Status.Phase == corev1.PodPending {
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{Requeue: true}, err
		}
		// 如果cr的status有变化，用当前cr的状态覆盖
		if !reflect.DeepEqual(helloCrdResourceOld.Status, helloCrdResource.Status) {
			err = r.Status().Update(context.TODO(), helloCrdResource)
			if err != nil {
				log.Error(err, "failed to update  status from running")
				return ctrl.Result{}, err
			} else {
				log.Info("updated  status RUNNING -> " + helloCrdResource.Status.HelloStatus)
				return ctrl.Result{Requeue: true}, nil
			}
		}
	default:

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HellocrdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demogroupv1.Hellocrd{}).
		Complete(r)
}
