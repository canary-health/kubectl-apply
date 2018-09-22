package handler

import (
	"fmt"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

// HandleDeployment applies `Kind: Deployment`
func (h *Handler) HandleDeployment(cmd *cobra.Command, args []string) {
	// Handle kubeconfig file
	kubeconfig := h.newKubeConfigHandler(cmd.Flag("kubeconfig").Value.String())
	// Set Namespace
	namespace := setNamespace(cmd.Flag("namespace").Value.String())
	// Set Deployment Client
	deploymentsClient := kubeconfig.Clientset.AppsV1().Deployments(namespace)
	// Decode deployment file
	d := newDecoderHandler(cmd.Flag("file").Value.String())
	// Set deployment object
	deployment := d.Obj.(*appsv1.Deployment)
	// Set labels
	deployment.Labels = labelHandler(cmd.Flag("labels").Value.String())
	// Set default replicas per env
	if deployment.Spec.Replicas == nil {
		if namespace == "development" {
			deployment.Spec.Replicas = ptrInt32(1)
		}
		if namespace == "acceptance" {
			deployment.Spec.Replicas = ptrInt32(2)
		}
		if namespace == "production" {
			deployment.Spec.Replicas = ptrInt32(3)
		}
	}
	// Set image if provided
	image := setImage(cmd.Flag("image").Value.String())
	// Set env params
	var paramPath string
	if cmd.Flag("paramPath").Value.String() != "" {
		paramPath = cmd.Flag("paramPath").Value.String()
	} else {
		paramPath = "/" + namespace + "/" + deployment.Name
	}
	envars := h.Aws.ssmByPathToK8sEnvVar(paramPath)
	cs := deployment.Spec.Template.Spec.Containers
	var ncs []v1.Container
	// Range Containers and set env vars
	for _, c := range cs {
		c.Image = image
		// Set SecurityContext if nil
		if c.SecurityContext == nil {
			v1sc := v1.SecurityContext{
				RunAsNonRoot:           setTrue(),
				ReadOnlyRootFilesystem: setTrue(),
			}
			c.SecurityContext = &v1sc
		}
		c.ImagePullPolicy = v1.PullAlways
		c.Env = envars
		ncs = append(ncs, c)
	}
	deployment.Spec.Template.Spec.Containers = ncs
	fmt.Println(deployment)
	// Applying deployment
	fmt.Println("Applying deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		// Print create err
		fmt.Println(err)
		fmt.Println("Updating deployment...")
		result, err := deploymentsClient.Update(deployment)
		if err != nil {
			// Panic on update err
			panic(err)
		}
		fmt.Printf("Updated deployment %q.\n", result.GetObjectMeta().GetName())
	} else {
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	}
}
