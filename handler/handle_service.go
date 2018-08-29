package handler

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HandleService applies `Kind: Service`
func (h *Handler) HandleService(cmd *cobra.Command, args []string) {
	// Set Namespace
	namespace := cmd.Flag("namespace").Value.String()
	if namespace == "" {
		namespace = v1.NamespaceDefault
	}
	// Handle kubeconfig file
	kubeconfig := h.newKubeConfigHandler(cmd.Flag("kubeconfig").Value.String())
	// Set Service Client
	servicesClient := kubeconfig.Clientset.CoreV1().Services(namespace)
	// Decode Service file
	d := newDecoderHandler(cmd.Flag("file").Value.String())
	// Set Service object from file
	service := d.Obj.(*v1.Service)
	// Set labels
	service.Labels = labelHandler(cmd.Flag("labels").Value.String())
	// Get existing Service by Name
	existS, err := servicesClient.Get(service.GetObjectMeta().GetName(), metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		// Create service
		fmt.Println("Creating service...")
		result, err := servicesClient.Create(service)
		if err != nil {
			// Print create err
			fmt.Println(err)
		}
		fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	} else {
		// Update service
		fmt.Println("Updating service...")
		service.ResourceVersion = existS.ResourceVersion
		service.Spec.ClusterIP = existS.Spec.ClusterIP
		result, err := servicesClient.Update(service)
		if err != nil {
			// Panic on update err
			panic(err)
		}
		fmt.Printf("Updated service %q.\n", result.GetObjectMeta().GetName())
	}
}
