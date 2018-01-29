# Kubernetes Flexvolume Provisioner and Driver

**NOTE: This solution is a proof of concept and will be replaced once k8s can dynamically provision local-storage.**

***Suggestions for improvements and Pull Requests highly welcome!***

The Provisioner and Flexvolume Driver are able to create dynamic local-storage for k8s.

The Provisioner creates and empty PV with a Flexvolume Driver definition which is then called by kubelet and creates a local folder which is linked into the pod volume folder on the node.
The node-affinity annotation is then updated to ensure node affinity if a pod is restarted.
On deletion the PVs are deleted in k8s but not on the filesystem.

### Usage:

Create a DaemonSet which deploys the flexvolume driver (these pods will be kept in an endless loop so that k8s doesn't restart them)
```bash
kubectl create -f deployment/daemonset.yaml
```

Make sure that `$HOME/.kube/config` exists on each node. This is needed so that the driver can update the PV with the node affinity annotation. Config is available when running `kubectl config view`.
For example (on the nodes itself):
```bash
mkdir -p $HOME/.kube && cat > $HOME/.kube/config <<EOF
apiVersion: v1
clusters: 
- cluster:
    insecure-skip-tls-verify: true
    server: kube-master:8080
contexts: []
current-context: ""
kind: Config
preferences: {}
users: []
EOF
```
Replace `kube-master:8080` with correct master IP and port if necessary.

Create the provisioner deployment
```bash
kubectl create -f deployment/deployment.yaml
```

Create the StorageClass with correct path to local-storage. You can define if the node-affinity alpha annotation should be used (feature-gates need to be set `PersistentLocalVolumes=true,VolumeScheduling=true,MountPropagation=true`)
```bash
kubectl create -f deployment/storageclass.yaml
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