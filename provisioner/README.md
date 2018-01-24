# Out-of-tree Dynamic Provisioner

A little demonstration of a Provisioner that creates local PVs. Upon delete it moves the created PV into an archived folder.

```bash
kubectl create -f - <<EOF
kind: Pod
apiVersion: v1
metadata:
  name: localvolume-provisioner
spec:
  containers:
    - name: localvolume-provisioner
      image: 937400120367.dkr.ecr.eu-west-1.amazonaws.com/localvolume-provisioner:latest
      imagePullPolicy: "IfNotPresent"
      env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
      volumeMounts:
        - name: pv-volume
          mountPath: /mnt/disks
  imagePullSecrets:
   - name: monoregistry
  volumes:
    - name: pv-volume
      hostPath:
        path: /mnt/disks
EOF
```

```bash
kubectl create -f - <<EOF
kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: example-localvolume
provisioner: monostream.com/localvolume-provisioner
parameters:
  path: /mnt/disks
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
  storageClassName: example-localvolume
EOF
```