package main

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetJobs lists all jobs in the specified namespace.
func GetJobs(clientset *kubernetes.Clientset, namespace string) {
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error getting jobs: %v\n", err)
		return
	}
	for _, job := range jobs.Items {
		fmt.Printf("Job Name: %s, Status: %v\n", job.GetName(), job.Status)
	}
}

// CreateJob creates a job in the specified namespace using a provided Job object.
func CreateJob(clientset *kubernetes.Clientset, namespace string, job *batchv1.Job) {
	result, err := clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job %s\n", result.GetObjectMeta().GetName())
}

// DeleteJob deletes a job by name in the specified namespace.
func DeleteJob(clientset *kubernetes.Clientset, namespace, jobName string) {
	err := clientset.BatchV1().Jobs(namespace).Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("Error deleting job: %v\n", err)
		return
	}
	fmt.Printf("Deleted job %s\n", jobName)
}

// ViewJob prints the status of a specific job by name.
func ViewJob(clientset *kubernetes.Clientset, namespace, jobName string) {
	job, err := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error viewing job: %v\n", err)
		return
	}
	fmt.Printf("Job Name: %s, Active: %d, Succeeded: %d, Failed: %d\n", job.GetName(), job.Status.Active, job.Status.Succeeded, job.Status.Failed)
}

func main() {
	// Connect to the Kubernetes cluster
	clientset, err := ConnectToCluster()
	if err != nil {
		fmt.Printf("Error connecting to Kubernetes: %v\n", err)
		return
	}
	fmt.Println("Connected to Kubernetes")

	// Specify the namespace
	namespace := "magicbox-christianiabpos"

	// Example: Get all jobs in the default namespace
	GetJobs(clientset, namespace)

	CreateSelfExecutingJob("kokskat-job", namespace)
	// // Example: View a specific job
	// jobName := "example-job"
	// ViewJob(clientset, namespace, jobName)

	// // Example: Delete a job
	// DeleteJob(clientset, namespace, jobName)

	// You can implement job creation based on your use case
	// Example: CreateJob(clientset, namespace, job)
}
