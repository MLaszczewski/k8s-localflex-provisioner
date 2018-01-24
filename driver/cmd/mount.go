// Copyright © 2018 munzli <manuel@monostream.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"flag"
	"errors"
	"encoding/json"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/monostream/k8s-localflex-provisioner/driver/helper"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var MountDir string

// mountCmd represents the mount command
var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Creates a directory",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least 2 args")
		}
		return nil
	},
	Long: `Creates a directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		var name string
		var affinity string
		var _ string

		// get the json options
		var options interface{}
		json.Unmarshal([]byte(args[1]), &options)
		mappedOptions := options.(map[string]interface{})

		for k, v := range mappedOptions {
			switch k {
			case "path":
				path = v.(string)
			case "name":
				name = v.(string)
			case "affinity":
				affinity = v.(string)
			case "directory":
				_ = v.(string)
			}
		}

		// if the target directory doesn't exist, create it
		if _, err := os.Stat(path); os.IsNotExist(err) {
			errDir := os.MkdirAll(path, 0755)
			if errDir != nil {
				helper.Handle(helper.Response{
					Status:  helper.StatusFailure,
					Message: errDir.Error(),
				})
				return
			}
		}

		// delete path generated by flexvolume and replace with symlink
		os.Remove(args[0])
		errLink := os.Symlink(path, args[0])
		if errLink != nil {
			helper.Handle(helper.Response{
				Status:  helper.StatusFailure,
				Message: errLink.Error(),
			})
			return
		}

		// update PV if affinity is set
		if affinity != "no" {
			err := updatePersistentVolume(name)
			if err != nil {
				helper.Handle(helper.Response{
					Status:  helper.StatusFailure,
					Message: err.Error(),
				})
				return
			}
		}

		helper.Handle(helper.Response{
			Status:  helper.StatusSuccess,
			Message: "successfully created the volume",
		})
	},
}

func init() {
	rootCmd.AddCommand(mountCmd)
}

func updatePersistentVolume(name string) error {
	// out of cluster config
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "/root/.kube/config", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return errors.New("build config: " + *kubeconfig + ": " + err.Error())
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New("clientset: " + err.Error())
	}

	nodeName, err := os.Hostname()
	if nodeName == "" {
		return errors.New("hostname: " + err.Error())
	}

	volumesClient := clientSet.CoreV1().PersistentVolumes()
	pv, err := volumesClient.Get(name, metav1.GetOptions{})
	if err != nil {
		return errors.New("get pv: " + err.Error())
	}

	// update affinity annotation
	annotations := pv.GetAnnotations()
	annotations[v1.AlphaStorageNodeAffinityAnnotation] = "{\"requiredDuringSchedulingIgnoredDuringExecution\":{\"nodeSelectorTerms\":[{\"matchExpressions\":[{\"key\":\"kubernetes.io/hostname\",\"operator\":\"In\",\"values\":[\"" + nodeName + "\"]}]}]}}"
	pv.SetAnnotations(annotations)

	_, error := volumesClient.Update(pv)
	if error != nil {
		return errors.New("update pv: " + error.Error())
	}

	// everything worked
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}