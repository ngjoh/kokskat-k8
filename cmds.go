package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newShowJobCmd(namespace *string) *cobra.Command {
	showJobCmd := &cobra.Command{
		Use:   "showjob",
		Short: "Show all jobs in the namespace",
		Run: func(cmd *cobra.Command, args []string) {
			handleShowJobCmd(*namespace)
		},
	}
	showJobCmd.Flags().StringVarP(namespace, "namespace", "n", "default", "Namespace to list jobs from")
	return showJobCmd
}

func handleShowJobCmd(namespace string) {
	clientset, err := ConnectToCluster()
	if err != nil {
		fmt.Printf("Error connecting to Kubernetes: %v\n", err)
		return
	}
	GetJobs(clientset, namespace)
}

func newSelfExecuteCmd(jobName, namespace *string) *cobra.Command {
	selfExecuteCmd := &cobra.Command{
		Use:   "selfexecute",
		Short: "Create a self-executing job",
		Run: func(cmd *cobra.Command, args []string) {
			handleSelfExecuteCmd(*jobName, *namespace)
		},
	}
	selfExecuteCmd.Flags().StringVarP(jobName, "name", "j", "", "Name of the job")
	selfExecuteCmd.Flags().StringVarP(namespace, "namespace", "n", "default", "Namespace to create the job in")
	selfExecuteCmd.MarkFlagRequired("name")
	return selfExecuteCmd
}

func handleSelfExecuteCmd(jobName, namespace string) {
	CreateSelfExecutingJob(jobName, namespace)
	fmt.Printf("Created self-executing job %s in namespace %s\n", jobName, namespace)
}
