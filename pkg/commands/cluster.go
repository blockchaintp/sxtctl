package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func listClusters(out, errOut io.Writer) error {

	err := checkFormat()

	if err != nil {
		return err
	}

	api, err := setupApi()
	if err != nil {
		return err
	}

	clusters, err := loadClusters(api)
	if err != nil {
		return err
	}

	if OutputFormat == "text" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Created", "Status", "API Server", "Deployments"})

		for _, cluster := range clusters {
			table.Append([]string{
				fmt.Sprintf("%d", cluster.Id),
				cluster.Name,
				cluster.CreatedAt,
				cluster.Status,
				cluster.AppliedState.ApiServer,
				cluster.ActiveDeployments,
			})
		}

		table.Render()
	} else {
		jsonString, err := getJSONString(clusters)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", jsonString)
	}

	return nil
}

func NewCmdClusterList(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List your current clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listClusters(out, errOut)
		},
	}

	return cmd
}

func NewCmdCluster(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage k8s clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listClusters(out, errOut)
		},
	}

	cmd.AddCommand(NewCmdClusterList(out, errOut))

	return cmd
}
