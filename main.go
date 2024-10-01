package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ConnectToCluster connects to a Kubernetes cluster, either from inside a pod or locally.
func ConnectToCluster() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Check if we're inside a Kubernetes pod
	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// If not in a pod, use the local kubeconfig file
		home := os.Getenv("HOME")
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

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

	// Example: View a specific job
	jobName := "example-job"
	ViewJob(clientset, namespace, jobName)

	// Example: Delete a job
	DeleteJob(clientset, namespace, jobName)

	// You can implement job creation based on your use case
	// Example: CreateJob(clientset, namespace, job)
}
