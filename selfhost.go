package main

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateSelfExecutingJob creates and runs a Kubernetes job, using Go install to dynamically install and run the program.
func CreateSelfExecutingJob(jobName, namespace string) error {
	clientset, err := ConnectToCluster()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %v", err)
	}

	ttl := int32(15 * 60) // TTL of 15 minutes after job completion

	// Define the job spec using the official Go image
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "go-job",
							Image: "golang:1.20", // Use the Go Docker image
							Command: []string{
								"sh", "-c",
								`
                                echo "Installing and running Go program...";
                                go install github.com/ngjoh/kokskat-k8@latest;
                                echo "Running installed Go program...";
                                /go/bin/kokskat-k8 showjob;
                                echo "Job finished.";
                                `,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			TTLSecondsAfterFinished: &ttl,
		},
	}

	// Create the job in the specified namespace
	_, err = clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job: %v", err)
	}

	fmt.Printf("Created job %s in namespace %s\n", jobName, namespace)
	return nil
}

// func xmain() {
// 	// Define parameters for the job
// 	jobName := "self-executing-go-job"
// 	namespace := "default"

// 	// Create and execute the self-deleting job
// 	err := CreateSelfExecutingJob(jobName, namespace)
// 	if err != nil {
// 		fmt.Printf("Error creating job: %v\n", err)
// 		return
// 	}

// 	fmt.Println("Job created and executing.")
// }
