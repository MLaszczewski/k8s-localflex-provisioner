# Flexvolume Driver

Move compiled binary to `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/monostream.com~localflex` on your node.
Also adds node-affinity alpha annotation to the PV.

Example PV:
```json
{
  "kind": "PersistentVolume",
  "apiVersion": "v1",
  "metadata": {
    "name": "my-volume",
    "annotations": {
      "monostream.com/node": "kube-node-1",
      "monostream.com/path": "/mnt/disks/my-volume",
      "pv.kubernetes.io/provisioned-by": "monostream.com/localflex-provisioner",
      "volume.alpha.kubernetes.io/node-affinity": "{\"requiredDuringSchedulingIgnoredDuringExecution\":{\"nodeSelectorTerms\":[{\"matchExpressions\":[{\"key\":\"kubernetes.io/hostname\",\"operator\":\"In\",\"values\":[\"kube-node-2\"]}]}]}}"
    }
  },
  "spec": {
    "capacity": {
      "storage": "8Gi"
    },
    "flexVolume": {
      "driver": "monostream.com/localflex",
      "options": {
        "directory": "/mnt/disks",
        "name": "my-volume",
        "path": "/mnt/disks/my-volume"
      }
    },
    "accessModes": [
      "ReadWriteOnce"
    ],
    "persistentVolumeReclaimPolicy": "Delete",
    "storageClassName": "localflex"
  },
}
```