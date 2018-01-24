package main

import (
	"os"
	"flag"
	"path"
	"errors"
	"syscall"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/util/wait"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	provisionerName =       "monostream.com/localflex-provisioner"
	pathAnnotation =        "monostream.com/path"
	nodeAnnotation =        "monostream.com/provisioner-node"
	flexDriver =            "monostream.com/localflex"
	diskPath =              "/tmp"
)

type localFlexProvisioner struct {
	// the directory to create pv-backing directories in
	pvDir string

	// define if node affinity annotation should be used
	affinity string

	// node of this localFlexProvisioner, set to NODE-NAME name. used to identify "this" provisioner's pvs.
	nodeName string
}

// NewLocalFlexProvisioner creates a new localflex provisioner
func NewLocalFlexProvisioner() controller.Provisioner {
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		glog.Fatal("env variable NODE_NAME must be set so that this provisioner can identify itself")
	}

	return &localFlexProvisioner{
		pvDir:    diskPath,
		nodeName: nodeName,
	}
}

var _ controller.Provisioner = &localFlexProvisioner{}

// Provision creates a storage asset and returns a pv object representing it.
func (p *localFlexProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {

	// if path parameter is passed use it
	for key, value := range options.Parameters {
		switch key {
		case "path":
			p.pvDir = value
		case "affinity":
			p.affinity = value
		}
	}
	pvPath := path.Join(p.pvDir, options.PVName)

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
			Annotations: map[string]string{
				nodeAnnotation: p.nodeName,
				pathAnnotation: pvPath,
				//v1.AlphaStorageNodeAffinityAnnotation: "",
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				FlexVolume: &v1.FlexVolumeSource{
					Driver: flexDriver,
					Options: map[string]string{
						"path": pvPath,
						"directory": p.pvDir,
						"name": options.PVName,
						"affinity": p.affinity,
					},
				},
			},
		},
	}

	return pv, nil
}

// Delete moves the storage asset that was created by Provision represented by the given pv.
func (p *localFlexProvisioner) Delete(volume *v1.PersistentVolume) error {

	// check provisioner node
	nodeName, ok := volume.Annotations[nodeAnnotation]
	if !ok {
		return errors.New("node annotation not found on PV")
	}
	if nodeName != p.nodeName {
		return &controller.IgnoredError{Reason: nodeName + ": node annotation on PV does not match " + p.nodeName}
	}

	// NOTE: this doesn't work in multi node setups
	// get volume path information from annotation
	/*pvPath, ok := volume.Annotations[pathAnnotation]
	if ok {
		pvDir := path.Dir(pvPath)
		archivedPath := path.Join(pvDir, "archived-"+volume.Name)

		// archive the volume
		if err := os.Rename(pvPath, archivedPath); err != nil {
			return err
		}
	}*/
	return nil
}

func main() {
	syscall.Umask(0)

	flag.Parse()
	flag.Set("logtostderr", "true")

	// create an InClusterConfig and use it to create a client for the controller to use to communicate with k8s
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatalf("Failed to create config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}

	// the controller needs to know what the server version is because out-of-tree provisioners aren't officially supported until 1.5
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		glog.Fatalf("Error getting server version: %v", err)
	}

	// create the provisioner: it implements the provisioner interface expected by the controller
	localFlexProvisioner := NewLocalFlexProvisioner()

	// start the provision controller which will dynamically provision local storage pvs
	pc := controller.NewProvisionController(
		clientset,
		provisionerName,
		localFlexProvisioner,
		serverVersion.GitVersion,
	)

	pc.Run(wait.NeverStop)
}