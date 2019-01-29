package handler

import (
	"fmt"

	"github.com/spf13/cobra"
	v1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
)

// HandleCronJob applies `Kind: CronJob`
func (h *Handler) HandleCronJob(cmd *cobra.Command, args []string) {
	// Set Namespace
	namespace := setNamespace(cmd.Flag("namespace").Value.String())
	// Handle kubeconfig file
	kubeconfig := h.newKubeConfigHandler(cmd.Flag("kubeconfig").Value.String())
	// Set CronJob Client
	cronjobsClient := kubeconfig.Clientset.BatchV1beta1().CronJobs(namespace)
	// Decode CronJob file
	d := newDecoderHandler(cmd.Flag("file").Value.String())
	// Set CronJob object from file
	cronjob := d.Obj.(*v1beta1.CronJob)
	// Set labels
	cronjob.Labels = labelHandler(cmd.Flag("labels").Value.String())
	// Set image if provided
	image := setImage(cmd.Flag("image").Value.String())
	// Set env params
	var paramPath string
	if cmd.Flag("paramPath").Value.String() != "" {
		paramPath = cmd.Flag("paramPath").Value.String()
	} else {
		paramPath = "/" + namespace + "/" + cronjob.Name
	}
	envars := h.Aws.ssmByPathToK8sEnvVar(paramPath)

	cs := cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers
	var ncs []v1.Container
	// Range Containers and set env vars
	for _, c := range cs {
		c.Image = image
		// Set SecurityContext if nil
		if c.SecurityContext == nil {
			v1sc := v1.SecurityContext{
				RunAsNonRoot:           setTrue(),
				RunAsUser:              ptrInt64(1001),
				ReadOnlyRootFilesystem: setTrue(),
			}
			c.SecurityContext = &v1sc
		}
		c.ImagePullPolicy = v1.PullAlways
		c.Env = envars
		ncs = append(ncs, c)
	}
	cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers = ncs

	// Applying cronjob
	fmt.Println("Applying cronjob...")
	result, err := cronjobsClient.Create(cronjob)
	if err != nil {
		// Print create err
		fmt.Println(err)
		fmt.Println("Updating cronjob...")
		result, err := cronjobsClient.Update(cronjob)
		if err != nil {
			// Panic on update err
			panic(err)
		}
		fmt.Printf("Updated cronjob %q.\n", result.GetObjectMeta().GetName())
	} else {
		fmt.Printf("Created cronjob %q.\n", result.GetObjectMeta().GetName())
	}
}
