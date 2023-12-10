package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cli/go-gh/pkg/tableprinter"
	sb "github.com/huandu/go-sqlbuilder"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	_ "modernc.org/sqlite"
)

var historyPath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Arc", "User Data", "Default", "History")

type HistoryEntry struct {
	ID            int    `db:"id" json:"id"`
	URL           string `db:"url" json:"url"`
	Title         string `db:"title" json:"title"`
	LastVisitedAt string `db:"lastVisitedAt" json:"lastVisitedAt"`
}

func NewCmdHistory() *cobra.Command {
	var flags struct {
		query string
		limit int
		json  bool
	}

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Search history",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, _ []string) error {
			dbFile, err := os.Open(historyPath)
			if err != nil {
				return fmt.Errorf("failed to open db file: %w", err)
			}
			defer dbFile.Close()

			tempfile, err := os.CreateTemp("", "arc-history-*.sqlite")
			if err != nil {
				return fmt.Errorf("failed to create tempfile: %w", err)
			}
			defer os.Remove(tempfile.Name())

			if _, err := io.Copy(tempfile, dbFile); err != nil {
				return fmt.Errorf("failed to copy db file: %w", err)
			}

			db, err := sql.Open("sqlite", tempfile.Name())
			if err != nil {
				return fmt.Errorf("failed to open db: %w", err)
			}

			sb := sb.NewSelectBuilder()
			sb.Select("id", "url", "title", sb.As("datetime(last_visit_time / 1000000 + (strftime('%s', '1601-01-01')), 'unixepoch', 'localtime')", "lastVisitedAt"))
			sb.From("urls")
			sb.GroupBy("url")
			sb.OrderBy("last_visit_time DESC")

			if flags.limit > 0 {
				sb.Limit(flags.limit)
			}

			if len(flags.query) > 0 {
				sb.Where(sb.Or(
					sb.Like("url", fmt.Sprintf("%%%s%%", flags.query)),
					sb.Like("title", fmt.Sprintf("%%%s%%", flags.query)),
				))
			}

			sql, sqlArgs := sb.Build()
			rows, err := db.Query(sql, sqlArgs...)
			if err != nil {
				return fmt.Errorf("failed to query: %w", err)
			}
			defer rows.Close()

			var entries []HistoryEntry
			for rows.Next() {
				var entry HistoryEntry
				if err := rows.Scan(&entry.ID, &entry.URL, &entry.Title, &entry.LastVisitedAt); err != nil {
					return fmt.Errorf("failed to scan: %w", err)
				}
				entries = append(entries, entry)
			}

			if flags.json {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				encoder.SetEscapeHTML(false)
				return encoder.Encode(entries)
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

			for _, entry := range entries {
				printer.AddField(entry.URL)
				printer.AddField(entry.Title)
				printer.AddField(entry.LastVisitedAt)
				printer.EndRow()
			}

			return printer.Render()
		},
	}

	cmd.Flags().IntVarP(&flags.limit, "limit", "l", 100, "limit")
	cmd.Flags().StringVarP(&flags.query, "query", "q", "", "query")

	return cmd
}
