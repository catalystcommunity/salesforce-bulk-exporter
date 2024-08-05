package cmd

import (
	"fmt"
	"strings"
	"time"

	sf "github.com/catalystcommunity/salesforce-bulk-exporter/internal/salesforce"
	"github.com/urfave/cli/v2"
)

var ExportCommand = &cli.Command{
	Name:      "export",
	Usage:     "Exports all object records from Salesforce",
	Args:      true,
	ArgsUsage: " object",
	Flags:     exportFlags,
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			_ = cli.ShowCommandHelp(ctx, "export")
			return fmt.Errorf("expected exactly one argument, got %d", ctx.NArg())
		}
		object := ctx.Args().First()

		err := sf.InitSFClient()
		if err != nil {
			return err
		}

		// generate the query to use in the job
		var query string
		var exportFields []string
		if exportCmdFieldsCli != nil {
			exportFields = exportCmdFieldsCli.Value()
		}
		if len(exportFields) == 0 {
			query, err = sf.GenerateQueryWithAllFields(object)
			if err != nil {
				return err
			}
		} else {
			query = fmt.Sprintf("SELECT %s FROM %s", strings.Join(exportFields[:], ", "), object)
		}

		// submit the query
		jobID, err := sf.SubmitBulkQueryJob(query, exportCmdIncludeArchived)
		if err != nil {
			return err
		}
		fmt.Printf("Submitted bulk query job with ID: %s\n", jobID)

		// wait for completion
		if exportCmdDownload {
			err = sf.WaitUntilJobComplete(jobID, exportCmdWaitInterval)
			if err != nil {
				return err
			}
			fmt.Println("Job completed, beginning download...")
			filenames, err := sf.SaveAllResults(jobID, exportCmdFilePrefix, exportCmdFileExt)
			if err != nil {
				return err
			}
			fmt.Printf("Completed download, saved to %d files\n", len(filenames))
		}
		return nil
	},
}

var (
	exportCmdDownload        bool
	exportCmdWaitInterval    time.Duration
	exportCmdFieldsCli       *cli.StringSlice
	exportCmdFilePrefix      string
	exportCmdFileExt         string
	exportCmdIncludeArchived bool
)

var exportFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "wait",
		Aliases:     []string{"w"},
		Usage:       "Wait for the job to complete and download",
		EnvVars:     []string{"EXPORT_DOWNLOAD"},
		Value:       false,
		Destination: &exportCmdDownload,
	},
	&cli.DurationFlag{
		Name:        "wait-interval",
		Aliases:     []string{"i"},
		Usage:       "Time to wait in between polls of job status",
		EnvVars:     []string{"EXPORT_WAIT_INTERVAL"},
		Value:       10 * time.Second,
		Destination: &exportCmdWaitInterval,
	},
	&cli.StringSliceFlag{
		Name:        "fields",
		Usage:       "Which fields to export, by default all fields are discovered",
		EnvVars:     []string{"EXPORT_FIELDS"},
		Value:       nil,
		Destination: exportCmdFieldsCli,
	},
	&cli.StringFlag{
		Name:        "filename-prefix",
		Aliases:     []string{"f"},
		Usage:       "Filename prefix for the downloaded files from Salesforce",
		EnvVars:     []string{"EXPORT_FILENAME_PREFIX"},
		Value:       "export",
		Destination: &exportCmdFilePrefix,
	},
	&cli.StringFlag{
		Name:        "file-extension",
		Aliases:     []string{"e"},
		Usage:       "Filename extension for the downloaded files from Salesforce",
		EnvVars:     []string{"EXPORT_FILE_EXTENSION"},
		Value:       "csv",
		Destination: &exportCmdFileExt,
	},
	&cli.BoolFlag{
		Name:        "include-archived",
		Aliases:     []string{"a"},
		Usage:       "Include archived records in the export",
		EnvVars:     []string{"EXPORT_INCLUDE_ARCHIVED"},
		Value:       false,
		Destination: &exportCmdIncludeArchived,
	},
}
