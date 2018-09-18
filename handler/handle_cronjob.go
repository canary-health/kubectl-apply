package handler

import (
	"fmt"

	"github.com/spf13/cobra"
	v1beta1 "k8s.io/api/batch/v1beta1"
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
