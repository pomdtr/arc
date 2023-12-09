package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	_ "embed"

	"github.com/cli/go-gh/pkg/tableprinter"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type Tab struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	WindowID int    `json:"windowId"`
	TabID    int    `json:"tabId"`
	Location string `json:"location"`
}

type State string

var (
	TabStateUnknown  State = ""
	TabStatePinned   State = "Pinned"
	TabStateUnpinned State = "Unpinned"
	TabStateFavorite State = "Favorite"
)

func (t Tab) State() State {
	switch t.Location {
	case "pinned":
		return TabStatePinned
	case "unpinned":
		return TabStateUnpinned
	case "topApp":
		return TabStateFavorite
	default:
		return TabStateUnknown
	}
}

func NewCmdTab() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tab",
	}

	cmd.AddCommand(NewCmdTabUrl())
	cmd.AddCommand(NewCmdTabTitle())
	cmd.AddCommand(NewCmdTabList())
	cmd.AddCommand(NewCmdTabFocus())
	cmd.AddCommand(NewCmdTabCreate())
	cmd.AddCommand(NewCmdTabClose())
	cmd.AddCommand(NewCmdTabReload())
	cmd.AddCommand(NewCmdTabUpdate())
	cmd.AddCommand(NewCmdTabExecute())

	return cmd
}

func NewCmdTabUrl() *cobra.Command {
	cmd := &cobra.Command{
		Use: "url",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runApplescript(`tell application "Arc" to set currentURL to URL of active tab of front window`)
			if err != nil {
				return err
			}

			if _, err := os.Stdout.Write(output); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func NewCmdTabTitle() *cobra.Command {
	cmd := &cobra.Command{
		Use: "title",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runApplescript(`tell application "Arc" to set currentTitle to title of active tab of window 1`)
			if err != nil {
				return err
			}

			if _, err := os.Stdout.Write(output); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func NewCmdTabCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func NewCmdTabFocus() *cobra.Command {
	var flags struct {
		ID int
	}
	cmd := &cobra.Command{
		Use: "focus",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := runApplescript(fmt.Sprintf(`tell application "Arc"
			tell front window
			  tell tab %d to focus
			end tell
			activate
		  end tell`, flags.ID))
			return err
		},
	}

	cmd.Flags().IntVar(&flags.ID, "id", 0, "tab id")
	cmd.MarkFlagRequired("id")

	return cmd
}

//go:embed applescript/list-tabs.applescript
var listTabsScript string

func NewCmdTabList() *cobra.Command {
	var flags struct {
		Pinned   bool
		Favorite bool
		Unpinned bool
		Json     bool
	}

	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := runApplescript(listTabsScript)
			if err != nil {
				return err
			}

			var tabs []Tab
			if err := json.Unmarshal(output, &tabs); err != nil {
				return err
			}

			var filteredTabs []Tab
			if !flags.Pinned && !flags.Unpinned && !flags.Favorite {
				filteredTabs = tabs
			} else {
				for _, tab := range tabs {
					if flags.Pinned && tab.Location == "pinned" {
						filteredTabs = append(filteredTabs, tab)
					}

					if flags.Unpinned && tab.Location == "unpinned" {
						filteredTabs = append(filteredTabs, tab)
					}

					if flags.Favorite && tab.Location == "topApp" {
						filteredTabs = append(filteredTabs, tab)
					}
				}
			}

			sort.SliceStable(filteredTabs, func(i, j int) bool {
				if filteredTabs[i].State() == filteredTabs[j].State() {
					return filteredTabs[i].TabID < filteredTabs[j].TabID
				}

				if filteredTabs[i].State() == TabStateFavorite {
					return true
				}

				if filteredTabs[j].State() == TabStateFavorite {
					return false
				}

				if filteredTabs[i].State() == TabStatePinned {
					return true
				}

				return false
			})

			if flags.Json {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				encoder.SetEscapeHTML(false)
				return encoder.Encode(filteredTabs)
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

			for _, tab := range filteredTabs {
				printer.AddField(strconv.Itoa(tab.TabID))
				printer.AddField(string(tab.State()))
				printer.AddField(tab.Title)
				printer.AddField(tab.URL)
				printer.EndRow()
			}

			return printer.Render()
		},
	}

	cmd.Flags().BoolVar(&flags.Json, "json", false, "output as json")
	cmd.Flags().BoolVar(&flags.Pinned, "pinned", false, "only show pinned tabs")
	cmd.Flags().BoolVar(&flags.Unpinned, "unpinned", false, "only show unpinned tabs")
	cmd.Flags().BoolVar(&flags.Favorite, "favorite", false, "only show favorite tabs")
	return cmd
}

func NewCmdTabClose() *cobra.Command {
	var flags struct {
		ID int
	}

	cmd := &cobra.Command{
		Use:  "close",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var osascript string
			if cmd.Flags().Changed("id") {
				osascript = fmt.Sprintf(`tell application "Arc"
					tell front window
					  tell tab %d to close
					end tell
				  end tell`, flags.ID)
			} else {
				osascript = `tell application "Arc"
					tell front window
					  tell active tab to close
					end tell
				end tell`
			}

			if _, err := runApplescript(osascript); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&flags.ID, "id", 0, "tab id")
	return cmd
}

func NewCmdTabReload() *cobra.Command {
	var flags struct {
		ID int
	}

	cmd := &cobra.Command{
		Use: "reload",

		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().IntVar(&flags.ID, "id", 0, "tab id")

	return cmd
}

func NewCmdTabUpdate() *cobra.Command {
	var flags struct {
		ID   int
		Eval string
	}

	cmd := &cobra.Command{
		Use:  "update <url>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().IntVar(&flags.ID, "id", 0, "tab id")
	cmd.Flags().StringVarP(&flags.Eval, "eval", "e", "", "javascript to evaluate")

	return cmd
}

func NewCmdTabExecute() *cobra.Command {
	var flags struct {
		ID int
	}

	cmd := &cobra.Command{
		Use: "execute <script>",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().IntVar(&flags.ID, "id", 0, "tab id")
	return cmd
}
