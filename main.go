package main

import (
	"fmt"
	"os"

	"github.com/catalystsquad/salesforce-bulk-exporter/cmd"
	"github.com/catalystsquad/salesforce-bulk-exporter/internal/salesforce"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "salesforce-bulk-exporter",
		Usage:   "Export data from Salesforce using the Bulk API",
		Suggest: true,
		Commands: []*cli.Command{
			cmd.ExportCommand,
			cmd.DownloadCommand,
			cmd.DescribeJobCommand,
			cmd.ListJobsCommand,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "base-url",
				Usage:       "salesforce base url",
				EnvVars:     []string{"SALESFORCE_BASE_URL"},
				Category:    "salesforce",
				Destination: &salesforce.BaseUrl,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "api-version",
				Value:       "55.0",
				Usage:       "salesforce api version",
				EnvVars:     []string{"SALESFORCE_API_VERSION"},
				Category:    "salesforce",
				Destination: &salesforce.ApiVersion,
			},
			&cli.StringFlag{
				Name:        "client-id",
				Usage:       "salesforce client id",
				EnvVars:     []string{"SALESFORCE_CLIENT_ID"},
				Category:    "salesforce",
				Destination: &salesforce.ClientId,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "client-secret",
				Usage:       "salesforce client secret",
				EnvVars:     []string{"SALESFORCE_CLIENT_SECRET"},
				Category:    "salesforce",
				Destination: &salesforce.ClientSecret,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "username",
				Usage:       "salesforce username",
				EnvVars:     []string{"SALESFORCE_USERNAME"},
				Category:    "salesforce",
				Destination: &salesforce.Username,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "password",
				Usage:       "salesforce password",
				EnvVars:     []string{"SALESFORCE_PASSWORD"},
				Category:    "salesforce",
				Destination: &salesforce.Password,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "grant-type",
				Value:       "password",
				Usage:       "salesforce grant type",
				EnvVars:     []string{"SALESFORCE_GRANT_TYPE"},
				Category:    "salesforce",
				Destination: &salesforce.GrantType,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
