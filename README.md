## 推送crd到集群
make install

## 推送cr到集群
kubectl apply -f config/sample/

## 启动controller（本地）
make run 

## operator 部署到集群

### 镜像打包
make docker-build docker-push IMG=canghong/hello-operator:alphav1

### 集群启动镜像
```
$make deploy IMG=canghong/hello-operator:alphav1
$kubectl get pod hello-operator-controller-manager-777bc5986d-tn9ww -n hello-operator-system
$kubectl logs -f  hello-operator-controller-manager-777bc5986d-tn9ww -n hello-operator-system
...
1.6659124929325175e+09	INFO	Starting workers	{"controller": "hellocrd", "controllerGroup": "demogroup.demo", "controllerKind": "Hellocrd", "worker count": 1}
```
此时看到controller日志正常输出


## 卸载

### 卸载crd
make uninstall

### 下载controller
make undeploy



