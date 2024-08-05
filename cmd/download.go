package cmd

import (
	"fmt"
	"time"

	sf "github.com/catalystcommunity/salesforce-bulk-exporter/internal/salesforce"
	"github.com/urfave/cli/v2"
)

var DownloadCommand = &cli.Command{
	Name:      "download",
	Usage:     "Downloads the results of a bulk job",
	Args:      true,
	ArgsUsage: " job_id",
	Flags:     downloadFlags,
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			_ = cli.ShowCommandHelp(ctx, "download")
			return fmt.Errorf("expected exactly one argument, got %d", ctx.NArg())
		}
		jobID := ctx.Args().First()

		err := sf.InitSFClient()
		if err != nil {
			return err
		}

		// check if or wait until job is done
		if downloadCmdWait {
			err = sf.WaitUntilJobComplete(jobID, downloadCmdWaitInterval)
			if err != nil {
				return err
			}
		} else {
			complete, state, err := sf.CheckIfJobComplete(jobID)
			if err != nil {
				return err
			}
			if !complete {
				return fmt.Errorf("Job not complete, current state is: %s\n", state)
			}
		}

		// download files
		filenames, err := sf.SaveAllResults(jobID, downloadCmdFilePrefix, downloadCmdFileExt)
		if err != nil {
			return err
		}
		fmt.Printf("Completed download, saved to %d files\n", len(filenames))
		return nil
	},
}

var (
	downloadCmdWait         bool
	downloadCmdWaitInterval time.Duration
	downloadCmdFilePrefix   string
	downloadCmdFileExt      string
)

var downloadFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "wait",
		Aliases:     []string{"w"},
		Usage:       "Wait for the job to complete",
		Destination: &downloadCmdWait,
	},
	&cli.DurationFlag{
		Name:        "wait-interval",
		Aliases:     []string{"i"},
		Usage:       "Time to wait in between polls of job status",
		Value:       10 * time.Second,
		Destination: &downloadCmdWaitInterval,
	},
	&cli.StringFlag{
		Name:        "filename-prefix",
		Aliases:     []string{"f"},
		Usage:       "Filename prefix for the downloaded files from Salesforce",
		Value:       "export",
		Destination: &downloadCmdFilePrefix,
	},
	&cli.StringFlag{
		Name:        "file-extension",
		Aliases:     []string{"e"},
		Usage:       "Filename extension for the downloaded files from Salesforce",
		Value:       "csv",
		Destination: &downloadCmdFileExt,
	},
}
