package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/stjudecloud/msgenctl/internal"
	"go.uber.org/zap"
)

var statusCmd = &cobra.Command{
	Use:   "status [workflow-id]",
	Short: "prints the status a workflow or all workflows",
	Args:  cobra.MaximumNArgs(1),
	RunE:  status,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) error {
	config, err := internal.ServiceConfigFromFlags(cmd.Flags())

	if err != nil {
		return err
	}

	client := internal.NewClient(config.BaseURL, config.AccessKey)

	if len(args) > 0 {
		rawWorkflowID, err := strconv.Atoi(args[0])

		if err != nil {
			return err
		}

		workflowID := internal.WorkflowID(rawWorkflowID)

		zap.S().Infow("status", "workflowID", workflowID)

		workflow, err := internal.FetchWorkflow(client, workflowID)

		if err != nil {
			return err
		}

		printWorkflow(workflow)
	} else {
		zap.S().Infow("status", "workflowID", "*")

		workflows, err := internal.FetchWorkflows(client)

		if err != nil {
			return err
		}

		printWorkflows(workflows)
	}

	return nil
}

func printWorkflows(workflows []internal.Workflow) {
	for _, workflow := range workflows {
		printWorkflow(workflow)
		fmt.Println()
	}
}

func printWorkflow(workflow internal.Workflow) {
	fmt.Printf("Workflow ID     : %v\n", workflow.ID)
	fmt.Printf("Status          : %s (%d)\n", workflow.Status, workflow.Status)
	fmt.Printf("Message         : %s\n", workflow.Message)
	fmt.Printf("Process         : %v\n", workflow.Process)
	fmt.Printf("Description     : %v\n", workflow.Description)
	fmt.Printf("Created Date    : %v\n", workflow.CreatedDate)
	fmt.Printf("End Date        : %v\n", workflow.EndDate)
	fmt.Printf("Wall Clock Time : %v\n", workflow.Duration())
	fmt.Printf("Bases Processed : %d\n", workflow.BasesProcessed)
}
