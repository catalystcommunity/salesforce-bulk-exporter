package cmd

import (
	"os"

	sf "github.com/catalystsquad/salesforce-bulk-exporter/internal/salesforce"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

var DescribeJobCommand = &cli.Command{
	Name:      "describe-job",
	Aliases:   []string{"describe"},
	Usage:     "Describes a bulk job",
	Args:      true,
	ArgsUsage: " job_id",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			_ = cli.ShowCommandHelp(ctx, "describe-job")
			return cli.Exit("expected exactly one argument, got %d", 1)
		}
		jobID := ctx.Args().First()

		err := sf.InitSFClient()
		if err != nil {
			return err
		}

		resp, err := sf.GetBulkJob(jobID)
		if err != nil {
			return err
		}

		// display all fields in a table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"Field Name", "Value"})

		t.AppendRow(table.Row{"Operation", resp.Operation})
		t.AppendRow(table.Row{"Object", resp.Object})
		t.AppendRow(table.Row{"CreatedById", resp.CreatedById})
		t.AppendRow(table.Row{"CreatedDate", resp.CreatedDate})
		t.AppendRow(table.Row{"SystemModstamp", resp.SystemModstamp})
		t.AppendRow(table.Row{"State", resp.State})
		t.AppendRow(table.Row{"ConcurrencyMode", resp.ConcurrencyMode})
		t.AppendRow(table.Row{"ContentType", resp.ContentType})
		t.AppendRow(table.Row{"ApiVersion", resp.ApiVersion})
		t.AppendRow(table.Row{"LineEnding", resp.LineEnding})
		t.AppendRow(table.Row{"ColumnDelimiter", resp.ColumnDelimiter})
		t.AppendRow(table.Row{"NumberRecordsProcessed", resp.NumberRecordsProcessed})
		t.AppendRow(table.Row{"Retries", resp.Retries})
		t.AppendRow(table.Row{"TotalProcessingTimeMilliseconds", resp.TotalProcessingTimeMilliseconds})

		t.Render()

		return nil
	},
}
