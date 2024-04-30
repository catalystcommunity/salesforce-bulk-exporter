package cmd

import (
	"fmt"
	"strings"
	"time"

	sf "github.com/catalystsquad/salesforce-bulk-exporter/internal/salesforce"
	"github.com/urfave/cli/v2"
)

var ExportCommand = &cli.Command{
	Name:      "export",
	Usage:     "Exports all object records from Salesforce",
	Args:      true,
	ArgsUsage: "object",
	Flags:     exportFlags,
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			return fmt.Errorf("expected exactly one argument, got %d", ctx.NArg())
		}
		object := ctx.Args().First()

		err := sf.InitSFClient()
		if err != nil {
			return err
		}

		// generate the query to use in the job
		var query string
		if len(exportCmdFields) == 0 {
			query, err = sf.GenerateQueryWithAllFields(object)
			if err != nil {
				return err
			}
		} else {
			query = fmt.Sprintf("SELECT %s FROM %s", strings.Join(exportCmdFields[:], ", "), object)
		}

		// submit the query
		jobID, err := sf.SubmitBulkQueryJob(query)
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
			filenames, err := sf.SaveAllResults(jobID, exportCmdFilePrefix, exportCmdFileExt)
			if err != nil {
				return err
			}
			fmt.Printf("Saved export to files: %s\n", strings.Join(filenames[:], ","))
		}
		return nil
	},
}

var (
	exportCmdDownload     bool
	exportCmdWaitInterval time.Duration
	exportCmdFieldsCli    *cli.StringSlice
	exportCmdFields       []string
	exportCmdFilePrefix   string
	exportCmdFileExt      string
)

var exportFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "download",
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
}
