# Kubernetes Flexvolume Provisioner and Driver

NOTE: This solution is a proof of concept and will be replaced once Kubernetes can dynamically provision local-storage.

This Provisioner and Flexvolume Driver are able to create dynamic local-storage for Kubernetes.

The Provisioner creates and empty PV with a Flexvolume Driver definition which is then called by kubelet and creates a local folder which is linked into the pod volume folder on the node.
The node-affinity annotation is then updated to ensure node affinity if a pod is restarted.

Usage:

Compile provisioner and driver with `make`.

Create the kubelet plugin directory on each k8s node
```bash
mkdir /usr/libexec/kubernetes/kubelet-plugins/volume/exec/monostream.com~localflex
```
and move the flexvolume driver created under `driver/localflex` there.

Make sure that `$HOME/.kube/config` exists on each node. Config is available when running `kubectl config view`.

Create the provisioner pod
```bash
kubectl create -f - <<EOF
kind: Pod
apiVersion: v1
metadata:
  name: localflex-provisioner
spec:
  containers:
    - name: localflex-provisioner
      image: monostream/localflex-provisioner:latest
      imagePullPolicy: "Always"
      env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
  imagePullSecrets:
   - name: monoregistry
EOF
```

Create the StorageClass with correct path to local-storage
```bash
kubectl create -f - <<EOF
kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: localflex
provisioner: monostream.com/localflex-provisioner
parameters:
  path: /mnt/disks
EOF
```

Test
```bash
kubectl create -f - <<EOF
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: local-test
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Mi
  storageClassName: localflex
EOF
```