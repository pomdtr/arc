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

func NewCmdSpace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Manage spaces",
	}

	cmd.AddCommand(NewCmdSpaceFocus())
	cmd.AddCommand(NewCmdSpaceList())
	return cmd
}

func NewCmdSpaceFocus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "focus",
		Short: "Focus a space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spaceID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			if _, err := runApplescript(fmt.Sprintf(`tell application "Arc"
				tell front window
			  		tell space %d to focus
				end tell
		  	end tell`, spaceID)); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

//go:embed applescript/list-spaces.applescript
var listSpacesScript string

type Space struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func NewCmdSpaceList() *cobra.Command {
	var flags struct {
		Json bool
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List spaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runApplescript(listSpacesScript)
			if err != nil {
				return err
			}

			var spaces []Space
			if err := json.Unmarshal(output, &spaces); err != nil {
				return err
			}

			if flags.Json {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				encoder.SetEscapeHTML(false)
				return encoder.Encode(spaces)
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

			for _, space := range spaces {
				printer.AddField(strconv.Itoa(space.ID))
				printer.AddField(space.Title)
				printer.EndRow()
			}

			return printer.Render()
		},
	}

	cmd.Flags().BoolVar(&flags.Json, "json", false, "output as json")
	return cmd
}
