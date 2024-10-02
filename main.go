package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var namespace string
	var jobName string
	var natsSubject string
	rootCmd := &cobra.Command{
		Use:   "kokskat-k8",
		Short: "Kubernetes Job Management CLI",
	}

	rootCmd.AddCommand(newShowJobCmd(&namespace))
	rootCmd.AddCommand(newSelfExecuteCmd(&jobName, &namespace))
	rootCmd.AddCommand(newStartCmd(&namespace, &natsSubject))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
