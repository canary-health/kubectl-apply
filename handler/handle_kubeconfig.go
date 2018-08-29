package handler

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kubeConfigHandler struct {
	Clientset *kubernetes.Clientset
}

func (h *Handler) newKubeConfigHandler(kubeconfig string) *kubeConfigHandler {
	if kubeconfig == "" {
		// Look for kube config in os var
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = os.Getenv("KUBECONFIG") + "/config"
			if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
				fmt.Println(err)

				// Look for kube config in standard path
				home, err := homedir.Dir()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				kubeconfig = home + "/.kube/config"
				if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
					panic(err)
				}
			}
		}
	}

	kc := h.Aws.kmsDecrypt(kubeconfig)
	defer os.Remove(kc) // delete decrypted temp file
	config, err := clientcmd.BuildConfigFromFlags("", kc)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return &kubeConfigHandler{clientset}
}
