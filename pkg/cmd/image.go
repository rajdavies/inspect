package cmd

import (
	"errors"
	"net/http"
	"os"

	"github.com/garethjevans/inspect/pkg/inspect"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ImageCmd struct for the image command.
type ImageCmd struct {
	Cmd    *cobra.Command
	Args   []string
	Client inspect.Client

	Repository string
	Tag        string
}

// NewImageCmd creates a new ImageCmd.
func NewImageCmd() *cobra.Command {
	c := &ImageCmd{
		Client: inspect.Client{
			Client: &http.Client{},
		},
	}
	cmd := &cobra.Command{
		Use:     "image <name>",
		Short:   "Inspect the docker container",
		Long:    "",
		Example: "",
		Aliases: []string{"i", "in", "ins"},
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				logrus.Fatalf("unable to run command: %s", err)
			}
		},
		Args: cobra.MaximumNArgs(1),
	}

	cmd.Flags().StringVarP(&c.Repository, "repository", "r", "",
		"Repository to query")
	cmd.Flags().StringVarP(&c.Tag, "tag", "t", "",
		"Tag to query")

	return cmd
}

// Run runs the command.
func (c *ImageCmd) Run() error {
	if len(c.Args) == 1 {
		c.Repository, c.Tag = ParseRepo(c.Args[0])
	}

	if c.Repository == "" {
		return errors.New("no repository has been configured")
	}

	if c.Tag == "" {
		return errors.New("no tag has been configured")
	}

	labels, err := c.Client.Labels(c.Repository, c.Tag)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(TableStyle)

	if Headers {
		t.AppendHeader(table.Row{"Label", "Value"})
	}

	if WriteSeparators {
		t.AppendSeparator()
	}

	for k, v := range labels {
		t.AppendRow(table.Row{k, v})
	}

	if WriteSeparators {
		t.AppendSeparator()
	}
	t.AppendRow(table.Row{"GitHub URL", inspect.GitHubURL(labels)})

	t.Render()

	return nil
}
