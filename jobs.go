package main

import (
	"context"
	"fmt"
	"log"

	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Define the start command
func newStartCmd(namespace, natsSubject *string) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start subscribing to Kubernetes job events and publish them to a NATS subject",
		Run: func(cmd *cobra.Command, args []string) {
			Start(*namespace, *natsSubject)
		},
	}

	// Flags for the start command
	startCmd.Flags().StringVarP(namespace, "namespace", "n", "default", "Namespace to subscribe to job events")
	startCmd.Flags().StringVarP(natsSubject, "subject", "s", "job.events", "NATS subject to publish job events")

	return startCmd
}

// Start function integrates the Kubernetes and NATS connections and subscribes to job events
func Start(namespace string, natsSubject string) {
	// Connect to Kubernetes Cluster
	clientset, err := ConnectToCluster()
	if err != nil {
		log.Fatalf("Failed to connect to Kubernetes: %v", err)
	}

	// Connect to NATS
	natsConn, err := ConnectToNATS()
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	// Subscribe to job events and publish them to NATS
	err = SubscribeToJobEvents(clientset, namespace, natsConn, natsSubject)
	if err != nil {
		log.Fatalf("Error subscribing to job events: %v", err)
	}
}

// Function to watch job events in a namespace and publish them to a NATS subject
func SubscribeToJobEvents(clientset *kubernetes.Clientset, namespace string, natsConn *nats.Conn, subject string) error {
	watcher, err := clientset.BatchV1().Jobs(namespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error watching job events: %v", err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		jobEvent, ok := event.Object.(*v1.Job)
		if !ok {
			log.Printf("unexpected type, ignoring: %T", event.Object)
			continue
		}

		// Prepare the event data
		eventData := map[string]interface{}{
			"type":   event.Type,
			"job":    jobEvent.Name,
			"status": jobEvent.Status,
		}

		// Convert the event data to JSON
		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			log.Printf("error marshalling job event to JSON: %v", err)
			continue
		}

		// Publish the event to the specified NATS subject
		err = natsConn.Publish(subject, eventJSON)
		if err != nil {
			log.Printf("error publishing job event to NATS: %v", err)
			continue
		}

		log.Printf("Published job event: %s to subject: %s", jobEvent.Name, subject)
	}

	return nil
}
