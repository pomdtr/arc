package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func runApplescript(code string) ([]byte, error) {
	output, err := exec.Command("osascript", "-e", code).Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s", exitError.Stderr)
		}

		return nil, err
	}

	return output, nil
}

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of Arc",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runApplescript(`tell application "Arc" to return version`)
			if err != nil {
				return err
			}

			cmd.Print(string(output))
			return nil
		},
	}

	return cmd
}

func buildDoc(command *cobra.Command) (string, error) {
	out := strings.Builder{}

	var page strings.Builder
	err := doc.GenMarkdown(command, &page)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(page.String(), "\n") {
		if strings.Contains(line, "SEE ALSO") {
			break
		}

		out.WriteString(line + "\n")
	}

	for _, child := range command.Commands() {
		if child.Hidden {
			continue
		}

		childPage, err := buildDoc(child)
		if err != nil {
			return "", err
		}
		out.WriteString(childPage)
	}

	return out.String(), nil
}

func NewDocCmd() *cobra.Command {
	docCmd := &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation for sunbeam",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			doc, err := buildDoc(cmd.Root())
			if err != nil {
				return err
			}

			fmt.Println("# Reference")
			fmt.Println()
			fmt.Println(doc)

			return nil
		},
	}
	return docCmd
}

func main() {
	cmd := cobra.Command{
		Use:          "arc",
		Short:        "Arc Companion CLI",
		SilenceUsage: true,
	}

	cmd.AddCommand(NewCmdTab())
	cmd.AddCommand(NewCmdSpace())
	cmd.AddCommand(NewCmdWindow())
	cmd.AddCommand(NewCmdHistory())
	cmd.AddCommand(NewCmdVersion())
	cmd.AddCommand(NewDocCmd())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
