package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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
		Short: "Manage tabs",
		Use:   "tab",
	}

	cmd.AddCommand(NewCmdTabGet())
	cmd.AddCommand(NewCmdTabList())
	cmd.AddCommand(NewCmdTabFocus())
	cmd.AddCommand(NewCmdTabCreate())
	cmd.AddCommand(NewCmdTabClose())
	cmd.AddCommand(NewCmdTabReload())
	cmd.AddCommand(NewCmdTabExecute())

	return cmd
}

func NewCmdTabGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get information about the active tab",
	}

	cmd.AddCommand(NewCmdTabUrl())
	cmd.AddCommand(NewCmdTabTitle())

	return cmd
}

func NewCmdTabUrl() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "url",
		Short: "Get the url of the active tab",
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
		Use:   "title",
		Short: "Get the title of the active tab",
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
	var flags struct {
		Space     int
		LittleArc bool
	}
	cmd := &cobra.Command{
		Use:   "create <url>",
		Short: `Create a new tab.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var osascript string
			if flags.LittleArc {
				osascript = fmt.Sprintf(`tell application "Arc" to make new tab with properties {URL:"%s"}`, args[0])
			} else if cmd.Flags().Changed("space") {
				osascript = fmt.Sprintf(`tell application "Arc"
				    tell space %d
					    make new tab with properties {URL:"%s"}
					end tell
					activate
			    end tell`, flags.Space, args[0])
			} else {
				osascript = fmt.Sprintf(`tell application "Arc"
					tell front window
					  make new tab with properties {URL:"%s"}
					end tell
					activate
				end tell`, args[0])
			}

			if _, err := runApplescript(osascript); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&flags.LittleArc, "little", false, "open in little arc")
	cmd.Flags().IntVar(&flags.Space, "space", 0, "space to create tab in")
	return cmd
}

func NewCmdTabFocus() *cobra.Command {
	var flags struct {
		ID int
	}

	cmd := &cobra.Command{
		Use:   "select",
		Short: "Select a tab by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := runApplescript(fmt.Sprintf(`tell application "Arc"
				tell front window
			  		tell tab %d to select
				end tell
				activate
			end tell`, flags.ID)); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&flags.ID, "id", "i", 0, "tab id")
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
		Use:   "list",
		Short: `List tabs`,
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
		ID []int
	}

	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close a tab",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("id") {
				osascript := `tell application "Arc"
					tell front window
					  tell active tab to close
					end tell
				end tell`

				if _, err := runApplescript(osascript); err != nil {
					return err
				}
			}

			for _, id := range flags.ID {
				osascript := fmt.Sprintf(`tell application "Arc"
				tell front window
				  tell tab %d to close
				end tell
			  end tell`, id)

				if _, err := runApplescript(osascript); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().IntSliceVarP(&flags.ID, "id", "i", nil, "tab id")
	return cmd
}

func NewCmdTabReload() *cobra.Command {
	var flags struct {
		ID []int
	}

	cmd := &cobra.Command{
		Use:   "reload",
		Short: `Reload a tab"`,

		RunE: func(cmd *cobra.Command, args []string) error {
			var osascript string
			if !cmd.Flags().Changed("id") {
				osascript = `tell application "Arc"
					tell front window
					  tell active tab to reload
					end tell
				end tell`
				if _, err := runApplescript(osascript); err != nil {
					return err
				}

				return nil
			}

			for _, id := range flags.ID {
				osascript = fmt.Sprintf(`tell application "Arc"
					tell front window
					  tell tab %d to reload
					end tell
				  end tell`, id)

				if _, err := runApplescript(osascript); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().IntSliceVarP(&flags.ID, "id", "i", nil, "tab id")

	return cmd
}

func NewCmdTabExecute() *cobra.Command {
	var flags struct {
		ID   int
		Eval string
	}

	cmd := &cobra.Command{
		Use:   "exec <script>",
		Short: "Execute javascript in the active tab",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(flags.Eval) == 0 && len(args) == 0 {
				return fmt.Errorf("either --eval or a script is required")
			}

			if len(flags.Eval) > 0 && len(args) > 0 {
				return fmt.Errorf("only one of --eval or a script is allowed")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var javascript string
			if len(flags.Eval) > 0 {
				javascript = escapeJavascript(flags.Eval)
			} else {
				content, err := os.ReadFile(args[0])
				if err != nil {
					return err
				}

				javascript = escapeJavascript(string(content))
			}

			var osascript string
			if cmd.Flags().Changed("id") {
				osascript = fmt.Sprintf(`tell application "Arc"
				tell front window
				  tell tab %d
				  	execute javascript "%s"
				  end tell
				end tell
			  end tell`, flags.ID, javascript)
			} else {
				osascript = fmt.Sprintf(`tell application "Arc"
				tell front window
				  tell active tab
				  	execute javascript "%s"
				  end tell
				end tell
			  end tell`, javascript)
			}

			output, err := runApplescript(osascript)
			if err != nil {
				return err
			}

			if (len(output)) > 0 {
				cmd.Print(string(output))
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&flags.ID, "id", "i", 0, "tab id")
	cmd.Flags().StringVarP(&flags.Eval, "eval", "e", "", "javascript to evaluate")
	return cmd
}

func escapeJavascript(javascript string) string {
	javascript = strings.ReplaceAll(javascript, `\`, `\\'`)
	javascript = strings.ReplaceAll(javascript, `"`, `\"`)
	return javascript
}
