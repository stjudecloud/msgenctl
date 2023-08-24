package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/stjudecloud/msgenctl/internal"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "submits a new workflow",
	RunE:  submit,
}

func init() {
	flags := submitCmd.Flags()

	// process
	flags.String("process-name", "", "process name")
	flags.String("process-args", "", "process arguments")

	// input
	flags.String("input-storage-connection-string", "", "input Azure Storage connection string")
	flags.String("input-storage-container-name", "", "input Azure Storage container name")
	flags.String("input-blob-name", "", "input blob name")

	flags.String("description", "", "workflow description")

	// output
	flags.String("output-storage-connection-string", "", "output Azure Storage connection string")
	flags.String("output-storage-container-name", "", "output Azure Storage container name")
	flags.String("output-basename", "", "output basename")
	flags.Bool("output-overwrite", false, "overwrite outputs")
	flags.Bool("output-include-log", true, "upload logs")

	// optional
	flags.String(
		"emit-ref-confidence",
		internal.ReferenceConfidenceModeNone,
		"mode for emitting reference confidence scores",
	)
	flags.Bool("bgzip-output", false, "compress VCF/GVCF files with bgzip")

	flags.Bool("ignore-azure-region", false, "allow data and service to be in different regions")

	rootCmd.AddCommand(submitCmd)
}

func submit(cmd *cobra.Command, args []string) error {
	config, err := internal.SubmitConfigFromFlags(cmd.Flags())

	if err != nil {
		return err
	}

	slog.Info("submit", "description", config.Description)

	client := internal.NewClient(config.Service.BaseURL, config.Service.AccessKey)
	workflow, err := internal.SubmitWorkflow(client, config)

	if err != nil {
		return err
	}

	fmt.Println(workflow.ID)

	return nil
}
