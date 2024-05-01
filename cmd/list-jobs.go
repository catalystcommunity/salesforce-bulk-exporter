package cmd

import (
	"os"

	sf "github.com/catalystsquad/salesforce-bulk-exporter/internal/salesforce"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

var ListJobsCommand = &cli.Command{
	Name:    "list-jobs",
	Aliases: []string{"list"},
	Usage:   "List current bulk jobs",
	Action: func(ctx *cli.Context) error {
		err := sf.InitSFClient()
		if err != nil {
			return err
		}

		jobs, err := sf.GetAllBulkJobs()
		if err != nil {
			return err
		}

		// list all jobs in a nice table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "Object", "Operation", "State", "SystemModstamp", "CreatedByID"})
		for _, job := range jobs {
			t.AppendRow(table.Row{
				job.ID, job.Object, job.Operation, job.State, job.SystemModstamp, job.CreatedById,
			})
		}
		t.SortBy([]table.SortBy{{Name: "SystemModstamp"}})
		t.Render()

		return nil
	},
}
