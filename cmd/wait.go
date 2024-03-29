package cmd

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stjudecloud/msgenctl/internal"
)

var waitCmd = &cobra.Command{
	Use:   "wait <workflow-id>",
	Short: "polls until the completion of a workflow",
	Args:  cobra.ExactArgs(1),
	RunE:  wait,
}

func init() {
	flags := waitCmd.Flags()

	flags.Int("interval", 60, "poll interval in seconds")

	rootCmd.AddCommand(waitCmd)
}

func wait(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	config, err := internal.ServiceConfigFromFlags(flags)

	if err != nil {
		return err
	}

	interval, err := intervalFromFlags(flags)

	if err != nil {
		return err
	}

	rawWorkflowID, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	workflowID := internal.WorkflowID(rawWorkflowID)

	slog.Info("wait", "workflowID", workflowID)

	client := internal.NewClient(config.BaseURL, config.AccessKey)

	for {
		workflow, err := internal.FetchWorkflow(client, workflowID)

		if err != nil {
			return err
		}

		slog.Info("wait", "workflowID", workflowID, "status", workflow.Status, "message", workflow.Message)

		switch workflow.Status {
		case internal.StatusSuccess:
			return nil
		case internal.StatusFailed, internal.StatusCancelled:
			return fmt.Errorf("workflow unsuccessful: %d: %s", workflow.Status, workflow.Message)
		default:
			time.Sleep(interval)
		}
	}
}

func intervalFromFlags(flags *pflag.FlagSet) (time.Duration, error) {
	rawInterval, err := flags.GetInt("interval")

	if err != nil {
		return 0, err
	}

	interval := time.Duration(rawInterval) * time.Second

	return interval, nil
}
