package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ClusterNameOrId string
var DeploymentNameOrId string

func listDeployments(out, errOut io.Writer) error {

	err := checkFormat()

	if err != nil {
		return err
	}

	api, cluster, err := initialiseCluster()
	if err != nil {
		return err
	}

	deployments, err := loadDeployments(api, cluster.Id)

	if err != nil {
		return err
	}

	if OutputFormat == "text" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Created", "Status", "DeploymentType", "DeploymentVersion"})

		for _, deployment := range deployments {
			table.Append([]string{
				fmt.Sprintf("%d", deployment.Id),
				deployment.Name,
				deployment.CreatedAt,
				deployment.Status,
				deployment.DeploymentType,
				deployment.DeploymentVersion,
			})
		}

		table.Render()
	} else {
		jsonString, err := getJSONString(deployments)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", jsonString)
	}

	return nil
}

func NewCmdDeploymentList(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List deployments on a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listDeployments(out, errOut)
		},
	}

	return cmd
}

func NewCmdDeploymentUndeploy(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undeploy",
		Short: "Pause a deployment on a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, _, deployment, err := initialiseDeployment()
			if err != nil {
				return err
			}

			if deployment.Status != "provisioned" {
				return fmt.Errorf("Deployment must be in provisioned state to un-deploy - currently: %s", deployment.Status)
			}

			err = undeployDeployment(api, deployment)
			if err != nil {
				return err
			}

			err = waitForDeploymentState(api, deployment, "deleted")
			if err != nil {
				return err
			}

			fmt.Printf("Deployment is now paused\n")
			return nil
		},
	}

	return cmd
}

func NewCmdDeploymentRedeploy(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeploy",
		Short: "Reactivate a deployment on a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, _, deployment, err := initialiseDeployment()
			if err != nil {
				return err
			}

			if deployment.Status != "deleted" {
				return fmt.Errorf("Deployment must be in paused state to re-deploy - currently: %s", deployment.Status)
			}

			err = redeployDeployment(api, deployment)
			if err != nil {
				return err
			}

			err = waitForDeploymentState(api, deployment, "provisioned")
			if err != nil {
				return err
			}

			fmt.Printf("Deployment is now reactivated\n")
			return nil
		},
	}

	return cmd
}

func NewCmdDeployment(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployment",
		Short: "Manage sextant deployments",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listDeployments(out, errOut)
		},
	}

	cmd.AddCommand(NewCmdDeploymentList(out, errOut))
	cmd.AddCommand(NewCmdDeploymentUndeploy(out, errOut))
	cmd.AddCommand(NewCmdDeploymentRedeploy(out, errOut))
	cmd.PersistentFlags().StringVarP(&ClusterNameOrId, "cluster", "c", os.Getenv("SEXTANT_CLUSTER"), "The name or id of the cluster")
	cmd.PersistentFlags().StringVarP(&DeploymentNameOrId, "deployment", "d", os.Getenv("SEXTANT_DEPLOYMENT"), "The name or id of the deployment")

	return cmd
}
