package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/spf13/cobra"
)

// JobRequest represents a request structure for job operations
type JobRequest struct {
	Action    string `json:"action"`
	JobName   string `json:"job_name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// JobResponse represents the response structure from the microservice
type JobResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// CreateJobHandler handles the "create" job operation
func CreateJobHandler(ctx context.Context, req micro.Request) {
	var jobReq JobRequest
	err := json.Unmarshal(req.Data(), &jobReq)
	if err != nil {
		log.Printf("Error unmarshalling request: %v", err)
		req.Respond([]byte(fmt.Sprintf(`{"status": "error", "message": "Invalid request format: %v"}`, err)))
		return
	}

	// Simulate job creation logic (e.g., interacting with Kubernetes)
	log.Printf("Creating job '%s' in namespace '%s'", jobReq.JobName, jobReq.Namespace)

	// Build response
	response := JobResponse{
		Status:  "success",
		Message: fmt.Sprintf("Job '%s' created in namespace '%s'", jobReq.JobName, jobReq.Namespace),
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		req.Respond([]byte(fmt.Sprintf(`{"status": "error", "message": "Failed to create job"}`)))
		return
	}

	// Send response
	req.Respond(responseData)
}

// DeleteJobHandler handles the "delete" job operation
func DeleteJobHandler(ctx context.Context, req micro.Request) {
	var jobReq JobRequest
	err := json.Unmarshal(req.Data(), &jobReq)
	if err != nil {
		log.Printf("Error unmarshalling request: %v", err)
		req.Respond([]byte(fmt.Sprintf(`{"status": "error", "message": "Invalid request format: %v"}`, err)))
		return
	}

	// Simulate job deletion logic (e.g., interacting with Kubernetes)
	log.Printf("Deleting job '%s' in namespace '%s'", jobReq.JobName, jobReq.Namespace)

	// Build response
	response := JobResponse{
		Status:  "success",
		Message: fmt.Sprintf("Job '%s' deleted from namespace '%s'", jobReq.JobName, jobReq.Namespace),
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		req.Respond([]byte(fmt.Sprintf(`{"status": "error", "message": "Failed to delete job"}`)))
		return
	}

	// Send response
	req.Respond(responseData)
}

// startKoksmatJobService starts the NATS microservice and registers it in Cobra as "service" command
func startKoksmatJobService(nc *nats.Conn, rootCmd *cobra.Command) {
	// Define the "service" Cobra command
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Starts the koksmat.jobs microservice",
		Run: func(cmd *cobra.Command, args []string) {
			// Create a new NATS Microservice
			svc, err := micro.AddService(nc, micro.Config{
				Name:        "koksmat.jobs",
				Version:     "1.0.0",
				Description: "A microservice to manage Kubernetes jobs",
			})
			if err != nil {
				log.Fatalf("Error creating microservice: %v", err)
			}

			// Add "create" endpoint
			err = svc.AddEndpoint("create", micro.HandlerFunc(func(req micro.Request) {
				CreateJobHandler(context.Background(), req)
			}))
			if err != nil {
				log.Fatalf("Error adding 'create' endpoint: %v", err)
			}

			// Add "delete" endpoint

			err = svc.AddEndpoint("delete", micro.HandlerFunc(func(req micro.Request) {
				DeleteJobHandler(context.Background(), req)
			}))
			fmt.Println("Microservice 'koksmat.jobs' is running...")

			// Keep the service running
			select {}
		},
	}

	// Register the "service" command under the root command
	rootCmd.AddCommand(serviceCmd)
}
