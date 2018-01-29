# Dynamic Flexvolume Provisioner

A Provisioner that creates local PVs. Make sure that the localflex driver is in your kubelet plugin directory.

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
EOF
```

```bash
kubectl create -f - <<EOF
kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: localflex
provisioner: monostream.com/localflex-provisioner
parameters:
  path: /mnt/disks
  affinity: "yes"
EOF
```

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