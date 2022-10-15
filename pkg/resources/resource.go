package resources

import (
	demogroupv1 "demo/hello-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getLabels(helloCrdResource *demogroupv1.Hellocrd) map[string]string {
	return map[string]string{
		"app":     helloCrdResource.Spec.ContainerImage,
		"version": helloCrdResource.Spec.ContainerTag,
	}
}

func CreatePod(helloCrdResource *demogroupv1.Hellocrd) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helloCrdResource.Spec.ContainerImage + helloCrdResource.Spec.ContainerTag,
			Namespace: helloCrdResource.Namespace,
			Labels:    getLabels(helloCrdResource),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  helloCrdResource.Spec.ContainerImage,
					Image: helloCrdResource.Spec.ContainerImageNamespace + "/" + helloCrdResource.Spec.ContainerImage + ":" + helloCrdResource.Spec.ContainerTag,
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}
