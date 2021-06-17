package commands

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/asaskevich/govalidator"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type AuthResponse struct {
	Ok bool `json:"ok"`
}

func listRemotes(out, errOut io.Writer) error {
	config, err := loadCLIConfig()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Url"})

	for _, remote := range config.Remotes {
		appendName := ""
		if config.ActiveRemote == remote.Name {
			appendName = " (*)"
		}
		table.Append([]string{
			remote.Name + appendName,
			remote.Url,
		})
	}

	table.Render()

	return nil
}

func NewCmdRemoteList(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List your current remotes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRemotes(out, errOut)
		},
	}

	return cmd
}

func NewCmdRemoteAdd(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add a new remote",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if Url == "" {
				return fmt.Errorf("Please provide a --url argument")
			}
			if AccessToken == "" {
				return fmt.Errorf("Please provide a --token argument")
			}

			if !govalidator.IsURL(Url) {
				return fmt.Errorf("'%s' is not a valid URL", Url)
			}

			name := args[0]

			re := regexp.MustCompile(`^[\w-]+$`)
			if !re.Match([]byte(name)) {
				return fmt.Errorf("Remote name must be alphanumeric")
			}

			existingRemote, err := loadRemoteByName(name)
			if err != nil {
				return err
			}
			if existingRemote != nil {
				return fmt.Errorf("Remote with that name already exists")
			}

			_, err = setupApi()
			if err != nil {
				return err
			}

			config, err := loadCLIConfig()
			if err != nil {
				return err
			}

			remote := Remote{
				Name:  name,
				Url:   Url,
				Token: AccessToken,
			}
			config.Remotes = append(config.Remotes, remote)
			config.ActiveRemote = name

			err = saveCLIConfig(config)

			if err != nil {
				return err
			}

			fmt.Fprintf(out, "remote %s added\n", name)

			return nil
		},
	}
	return cmd
}

func NewCmdRemoteRemove(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove a remote",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			existingRemote, err := loadRemoteByName(name)
			if err != nil {
				return err
			}
			if existingRemote == nil {
				return fmt.Errorf("No remote with that name found")
			}
			config, err := loadCLIConfig()
			if err != nil {
				return err
			}

			newRemotes := []Remote{}

			for _, remote := range config.Remotes {
				if remote.Name != name {
					newRemotes = append(newRemotes, remote)
				}
			}

			config.Remotes = newRemotes

			if name == config.ActiveRemote {
				if len(config.Remotes) > 0 {
					config.ActiveRemote = config.Remotes[0].Name
				} else {
					config.ActiveRemote = ""
				}
			}

			err = saveCLIConfig(config)

			if err != nil {
				return err
			}

			fmt.Fprintf(out, "remote %s removed\n", name)

			return nil
		},
	}
	return cmd
}

func NewCmdRemoteUse(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [name]",
		Short: "Default to a remote",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			existingRemote, err := loadRemoteByName(name)
			if err != nil {
				return err
			}
			if existingRemote == nil {
				return fmt.Errorf("No remote with that name found")
			}
			config, err := loadCLIConfig()
			if err != nil {
				return err
			}

			config.ActiveRemote = name

			err = saveCLIConfig(config)

			if err != nil {
				return err
			}

			fmt.Fprintf(out, "remote %s activated\n", name)

			return nil
		},
	}
	return cmd
}

func NewCmdRemote(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage your sextant credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRemotes(out, errOut)
		},
	}

	cmd.AddCommand(NewCmdRemoteList(out, errOut))
	cmd.AddCommand(NewCmdRemoteAdd(out, errOut))
	cmd.AddCommand(NewCmdRemoteRemove(out, errOut))
	cmd.AddCommand(NewCmdRemoteUse(out, errOut))

	return cmd
}
