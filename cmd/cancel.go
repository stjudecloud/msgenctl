package cmd

import (
	"log/slog"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/stjudecloud/msgenctl/internal"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel <workflow-id>",
	Short: "cancels a running workflow",
	Args:  cobra.ExactArgs(1),
	RunE:  cancel,
}

func init() {
	rootCmd.AddCommand(cancelCmd)
}

func cancel(cmd *cobra.Command, args []string) error {
	config, err := internal.ServiceConfigFromFlags(cmd.Flags())

	if err != nil {
		return err
	}

	client := internal.NewClient(config.BaseURL, config.AccessKey)

	rawWorkflowID, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	workflowID := internal.WorkflowID(rawWorkflowID)

	slog.Info("cancel", "workflowID", workflowID)

	workflow, err := internal.CancelWorkflow(client, workflowID)

	if err != nil {
		return err
	}

	printWorkflow(workflow)

	return nil
}
