package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func NewCmdWindow() *cobra.Command {
	cmd := &cobra.Command{
		Use: "window",
	}

	cmd.AddCommand(NewCmdWindowCreate())

	return cmd
}

func NewCmdWindowCreate() *cobra.Command {
	var flags struct {
		Incognito bool
	}

	cmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			var osascript string
			if flags.Incognito {
				osascript = `tell application "Arc"
					make new window with properties {incognito:true}
					activate
				end tell`
			} else {
				osascript = `tell application "Arc"
					make new window
					activate
				end tell`
			}

			if _, err := runApplescript(osascript); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&flags.Incognito, "incognito", false, "open in incognito mode")

	return cmd
}

func NewCmdWindowList() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO
			return nil
		},
	}

	return cmd
}

func NewCmdWindowClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "close",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var osascript string
			if len(args) > 0 {
				index, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("invalid window index: %s", args[0])
				}

				osascript = fmt.Sprintf(`tell application "Arc"
					close window %d
				end tell`, index)
			} else {
				osascript = `tell application "Arc"
					close front window
				end tell`
			}

			if _, err := runApplescript(osascript); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
