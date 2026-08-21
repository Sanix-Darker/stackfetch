package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type qaEntry struct {
	ID           string
	Area         string
	Feature      string
	UserStory    string
	Expected     string
	Status       string
	TestCoverage string
	OpenIssues   string
	Resolution   string
}

type qaSummary struct {
	total        int
	byStatus     map[string]int
	byArea       map[string]int
	blockedItems int
}

const featureTrackerPath = "FEATURE_TRACKER.csv"

func newQACmd(markdown bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qa",
		Short: "Display post-release feature stories and QA matrix",
		RunE: func(cmd *cobra.Command, args []string) error {
			area, _ := cmd.Flags().GetString("area")
			tracker, err := loadFeatureTrackerCSV()
			if err != nil {
				return err
			}
			entries, err := parseFeatureTracker(tracker)
			if err != nil {
				return err
			}
			filtered := filterQAEntries(entries, area)
			renderQAMatrix(filtered, markdown)
			return nil
		},
	}
	cmd.Flags().String("area", "", "Filter matrix by area (case-insensitive)")
	return cmd
}

func parseFeatureTracker(raw []byte) ([]qaEntry, error) {
	r := csv.NewReader(bytes.NewReader(raw))
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) <= 1 {
		return nil, nil
	}

	entries := make([]qaEntry, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 9 {
			continue
		}
		entries = append(entries, qaEntry{
			ID:           strings.TrimSpace(row[0]),
			Area:         strings.TrimSpace(row[1]),
			Feature:      strings.TrimSpace(row[2]),
			UserStory:    strings.TrimSpace(row[3]),
			Expected:     strings.TrimSpace(row[4]),
			Status:       strings.TrimSpace(row[5]),
			TestCoverage: strings.TrimSpace(row[6]),
			OpenIssues:   strings.TrimSpace(row[7]),
			Resolution:   strings.TrimSpace(row[8]),
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Area == entries[j].Area {
			return entries[i].ID < entries[j].ID
		}
		return entries[i].Area < entries[j].Area
	})
	return entries, nil
}

func summarizeQA(entries []qaEntry) qaSummary {
	s := qaSummary{
		byStatus: make(map[string]int),
		byArea:   make(map[string]int),
	}
	for _, e := range entries {
		s.total++
		s.byStatus[e.Status]++
		s.byArea[e.Area]++
		if strings.EqualFold(e.Status, "fail") || strings.EqualFold(e.Status, "blocked") || strings.EqualFold(e.Status, "todo") {
			s.blockedItems++
		}
	}
	return s
}

func filterQAEntries(entries []qaEntry, area string) []qaEntry {
	if area == "" {
		return entries
	}
	out := make([]qaEntry, 0, len(entries))
	for _, e := range entries {
		if strings.EqualFold(e.Area, area) {
			out = append(out, e)
		}
	}
	return out
}

func renderQAMatrix(entries []qaEntry, markdown bool) {
	s := summarizeQA(entries)
	if markdown {
		fmt.Printf("### QA matrix\n\n")
		fmt.Printf("| ID | Area | Feature | Status | Coverage | Notes |\n")
		fmt.Printf("| --- | --- | --- | --- | --- | --- |\n")
		for _, e := range entries {
			notes := e.OpenIssues
			if e.Resolution != "" {
				notes = e.Resolution
			}
			fmt.Printf("| %s | %s | %s | %s | %s | %s |\n", e.ID, e.Area, e.Feature, e.Status, e.TestCoverage, notes)
		}
		fmt.Printf("\n")
		fmt.Printf("- Total stories: %d\n", s.total)
		for status, n := range s.byStatus {
			fmt.Printf("- %s: %d\n", status, n)
		}
		return
	}

	rows := [][]string{}
	for _, e := range entries {
		notes := e.OpenIssues
		if e.Resolution != "" {
			notes = e.Resolution
		}
		rows = append(rows, []string{e.ID, e.Area, e.Feature, e.Status, notes})
	}
	fmt.Printf("Stories: %d\n", s.total)
	fmt.Printf("Blocked items: %d\n", s.blockedItems)
	for area, n := range s.byArea {
		fmt.Printf("  %s: %d\n", area, n)
	}
	fmt.Printf("\n")
	fmt.Println("ID\tArea\tFeature\tStatus\tCoverage/Notes")
	for _, row := range rows {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", row[0], row[1], row[2], row[3], row[4])
	}
}

func loadFeatureTrackerCSV() ([]byte, error) {
	if tracker, err := os.ReadFile(filepath.Clean(featureTrackerPath)); err == nil {
		return tracker, nil
	}
	if tracker, err := os.ReadFile(filepath.Join("..", featureTrackerPath)); err == nil {
		return tracker, nil
	}
	return nil, fmt.Errorf("cannot locate feature tracker file")
}
