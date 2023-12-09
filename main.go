package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
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
		Use: "version",
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

func main() {
	cmd := cobra.Command{
		Use:          "arc",
		SilenceUsage: true,
	}

	cmd.AddCommand(NewCmdTab())
	cmd.AddCommand(NewCmdSpace())
	cmd.AddCommand(NewCmdWindow())
	cmd.AddCommand(NewCmdHistory())
	cmd.AddCommand(NewCmdVersion())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
