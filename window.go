package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	_ "embed"

	"github.com/cli/go-gh/pkg/tableprinter"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type Window struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func NewCmdWindow() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Manage windows",
		Use:   "window",
	}

	cmd.AddCommand(NewCmdWindowCreate())
	cmd.AddCommand(NewCmdWindowClose())
	cmd.AddCommand(NewCmdWindowList())

	return cmd
}

func NewCmdWindowCreate() *cobra.Command {
	var flags struct {
		Incognito bool
		Little    bool
	}

	cmd := &cobra.Command{
		Use:   "create [url]",
		Short: "Create a new window",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var applescript string
			if flags.Incognito {
				applescript = `tell application "Arc"
					make new window with properties {incognito:true}
					activate
				end tell`
			} else {
				applescript = `tell application "Arc"
					make new window
				end tell`
			}

			if _, err := runApplescript(applescript); err != nil {
				return err
			}

			if len(args) > 0 {
				if _, err := runApplescript(fmt.Sprintf(`tell application "Arc"
					tell front window
						make new tab with properties {URL:"%s"}
					end tell
				end tell`, args[0])); err != nil {
					return err
				}
			}

			if _, err := runApplescript(`tell application "Arc" to activate`); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&flags.Incognito, "incognito", false, "open in incognito mode")

	return cmd
}

//go:embed applescript/list-windows.applescript
var listWindowsScript string

func NewCmdWindowList() *cobra.Command {
	flags := struct {
		Json bool
	}{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List windows",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runApplescript(listWindowsScript)
			if err != nil {
				return err
			}

			var windows []Window
			if err := json.Unmarshal(output, &windows); err != nil {
				return err
			}

			if flags.Json {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				encoder.SetEscapeHTML(false)
				return encoder.Encode(windows)
			}

			var printer tableprinter.TablePrinter
			if !isatty.IsTerminal(os.Stdout.Fd()) {
				printer = tableprinter.New(os.Stdout, false, 0)
			} else {
				w, _, err := term.GetSize(int(os.Stdout.Fd()))
				if err != nil {
					return err
				}

				printer = tableprinter.New(os.Stdout, true, w)
			}

			for _, window := range windows {
				printer.AddField(fmt.Sprintf("%d", window.ID))
				printer.AddField(window.Title)
				printer.EndRow()
			}

			return printer.Render()
		},
	}

	cmd.Flags().BoolVar(&flags.Json, "json", false, "output as json")
	return cmd
}

func NewCmdWindowClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close a window",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				if _, err := runApplescript(`tell application "Arc" to tell front window to close`); err != nil {
					return err
				}
				return nil
			}

			for _, id := range args {
				windowID, err := strconv.Atoi(id)
				if err != nil {
					return err
				}

				if _, err := runApplescript(fmt.Sprintf(`tell application "Arc" to tell window %d to close`, windowID)); err != nil {
					return err
				}

			}
			return nil
		},
	}

	return cmd
}
