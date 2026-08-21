package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type qaEntry struct {
	line         int
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

type qaStatusValues struct {
	Pass     int
	Blocked  int
	Todo     int
	Fail     int
	Other    int
	Total    int
	PassRate int
}

const featureTrackerPath = "FEATURE_TRACKER.csv"

func newQACmd(markdown bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qa",
		Short: "Display post-release feature stories and QA matrix",
		RunE: func(cmd *cobra.Command, args []string) error {
			area, _ := cmd.Flags().GetString("area")
			statusFilterRaw, _ := cmd.Flags().GetString("status")
			strict, _ := cmd.Flags().GetBool("strict")
			validateOnly, _ := cmd.Flags().GetBool("validate")

			jsonOut := isJSONRequested(cmd)
			trailer, err := loadFeatureTrackerCSV()
			if err != nil {
				return err
			}
			entries, parseIssues, err := parseFeatureTracker(trailer)
			if err != nil {
				return err
			}

			issues := append([]string{}, parseIssues...)
			issues = append(issues, validateFeatureTracker(entries)...)
			report := summarizeQA(entries)

			if validateOnly {
				return renderQAValidation(entries, report, issues, jsonOut, cmd.OutOrStdout())
			}

			filtered := filterQAEntries(entries, area, parseStatusFilter(statusFilterRaw))
			summary := summarizeQA(filtered)
			if err := renderQAMatrix(filtered, summary, issues, markdown, jsonOut, cmd.OutOrStdout()); err != nil {
				return err
			}
			if strict && summary.blockedItems > 0 {
				return fmt.Errorf("qa strict check failed: %d blocked stories", summary.blockedItems)
			}
			if strict && len(issues) > 0 {
				return fmt.Errorf("qa validation issues: %d", len(issues))
			}
			return nil
		},
	}
	cmd.Flags().String("area", "", "Filter matrix by area (case-insensitive)")
	cmd.Flags().String("status", "", "Filter matrix by status (comma-separated, case-insensitive)")
	cmd.Flags().Bool("strict", false, "Exit non-zero when blocked stories are present in the filtered set")
	cmd.Flags().Bool("validate", false, "Validate tracker integrity and exit with non-zero on issues")
	return cmd
}

func parseFeatureTracker(raw []byte) ([]qaEntry, []string, error) {
	r := csv.NewReader(bytes.NewReader(raw))
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	rows := []qaEntry{}
	issues := []string{}
	line := 0
	seenHeader := false

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, issues, err
		}
		line++
		if len(row) == 0 {
			continue
		}
		allBlank := true
		for _, field := range row {
			if strings.TrimSpace(field) != "" {
				allBlank = false
				break
			}
		}
		if allBlank {
			continue
		}

		if !seenHeader {
			seenHeader = true
			continue
		}

		if len(row) < 9 {
			issues = append(issues, fmt.Sprintf("line %d: expected 9 fields, got %d", line, len(row)))
			continue
		}
		if len(row) > 9 {
			issues = append(issues, fmt.Sprintf("line %d: expected 9 fields, got %d", line, len(row)))
			continue
		}

		rows = append(rows, qaEntry{
			line:         line,
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

	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Area == rows[j].Area {
			return rows[i].ID < rows[j].ID
		}
		return rows[i].Area < rows[j].Area
	})
	return rows, issues, nil
}

func parseStatusFilter(raw string) map[string]struct{} {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	statuses := make(map[string]struct{}, len(parts))
	for _, status := range parts {
		normalized := normalizeStatus(status)
		if normalized == "" {
			continue
		}
		statuses[normalized] = struct{}{}
	}
	if len(statuses) == 0 {
		return nil
	}
	return statuses
}

func normalizeStatus(raw string) string {
	return strings.ToUpper(strings.TrimSpace(raw))
}

func validateFeatureTracker(entries []qaEntry) []string {
	issues := []string{}
	ids := map[string]int{}

	for _, entry := range entries {
		if entry.ID == "" {
			issues = append(issues, fmt.Sprintf("line %d: missing id", entry.line))
		}
		if entry.Area == "" {
			issues = append(issues, fmt.Sprintf("line %d: %s missing area", entry.line, entry.ID))
		}
		if entry.Feature == "" {
			issues = append(issues, fmt.Sprintf("line %d: %s missing feature", entry.line, entry.ID))
		}
		if entry.UserStory == "" {
			issues = append(issues, fmt.Sprintf("line %d: %s missing user story", entry.line, entry.ID))
		}
		if entry.Expected == "" {
			issues = append(issues, fmt.Sprintf("line %d: %s missing expected behavior", entry.line, entry.ID))
		}
		if entry.Status == "" {
			issues = append(issues, fmt.Sprintf("line %d: %s missing status", entry.line, entry.ID))
		}
		status := normalizeStatus(entry.Status)
		if status != "" && !isKnownStatus(status) {
			issues = append(issues, fmt.Sprintf("line %d: %s unknown status %q", entry.line, entry.ID, entry.Status))
		}
		if entry.ID != "" {
			if firstLine, exists := ids[entry.ID]; exists {
				issues = append(issues, fmt.Sprintf("line %d: duplicate id %s (first seen on line %d)", entry.line, entry.ID, firstLine))
			} else {
				ids[entry.ID] = entry.line
			}
		}
		if entry.TestCoverage == "" {
			issues = append(issues, fmt.Sprintf("line %d: %s missing test coverage", entry.line, entry.ID))
		}
	}

	return issues
}

func isKnownStatus(status string) bool {
	switch status {
	case "PASS", "TODO", "FAIL", "BLOCKED", "IN_PROGRESS", "WIP":
		return true
	default:
		return false
	}
}

func summarizeQA(entries []qaEntry) qaSummary {
	s := qaSummary{
		byStatus: make(map[string]int),
		byArea:   make(map[string]int),
	}
	for _, e := range entries {
		s.total++
		status := normalizeStatus(e.Status)
		s.byStatus[status]++
		s.byArea[e.Area]++
		if status == "FAIL" || status == "BLOCKED" || status == "TODO" {
			s.blockedItems++
		}
	}
	return s
}

func summarizeQAStatus(entryCount int, byStatus map[string]int) qaStatusValues {
	sv := qaStatusValues{Total: entryCount}
	for status, n := range byStatus {
		sv.Total += 0
		switch status {
		case "PASS":
			sv.Pass += n
		case "BLOCKED":
			sv.Blocked += n
		case "TODO":
			sv.Todo += n
		case "FAIL":
			sv.Fail += n
		default:
			sv.Other += n
		}
	}
	if entryCount > 0 {
		sv.PassRate = (sv.Pass * 100) / entryCount
	}
	return sv
}

func filterQAEntries(entries []qaEntry, area string, statusSet map[string]struct{}) []qaEntry {
	if area == "" && len(statusSet) == 0 {
		return entries
	}

	normalizedArea := strings.ToLower(strings.TrimSpace(area))
	out := make([]qaEntry, 0, len(entries))
	for _, e := range entries {
		if normalizedArea != "" && strings.ToLower(e.Area) != normalizedArea {
			continue
		}
		if len(statusSet) != 0 {
			if _, ok := statusSet[normalizeStatus(e.Status)]; !ok {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}

func renderQAMatrix(entries []qaEntry, summary qaSummary, issues []string, markdown bool, jsonOut bool, out io.Writer) error {
	if jsonOut {
		return renderQAJSON(entries, summary, issues, out)
	}
	if markdown {
		return renderQAMarkdown(entries, summary, issues, out)
	}
	return renderQAPlain(entries, summary, issues, out)
}

func renderQAMarkdown(entries []qaEntry, summary qaSummary, issues []string, out io.Writer) error {
	fmt.Fprintf(out, "### QA matrix\n\n")
	fmt.Fprintf(out, "| ID | Area | Feature | Status | Coverage | Notes |\n")
	fmt.Fprintf(out, "| --- | --- | --- | --- | --- | --- |\n")
	for _, e := range entries {
		notes := e.OpenIssues
		if e.Resolution != "" {
			notes = e.Resolution
		}
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s | %s |\n", e.ID, e.Area, e.Feature, e.Status, e.TestCoverage, notes)
	}
	fmt.Fprintf(out, "\n")
	writeSummaryBlocks(out, summary, issues)
	return nil
}

func renderQAPlain(entries []qaEntry, summary qaSummary, issues []string, out io.Writer) error {
	fmt.Fprintf(out, "Stories: %d\n", summary.total)
	fmt.Fprintf(out, "Blocked items: %d\n", summary.blockedItems)
	fmt.Fprintf(out, "Pass rate: %d%%\n", summarizeQAStatus(summary.total, summary.byStatus).PassRate)

	for _, area := range sortedKeys(summary.byArea) {
		fmt.Fprintf(out, "  %s: %d\n", area, summary.byArea[area])
	}

	writeSummaryBlocks(out, summary, issues)

	fmt.Fprint(out, "\n")
	fmt.Fprintln(out, "ID\tArea\tFeature\tStatus\tCoverage/Notes")
	for _, e := range entries {
		notes := e.OpenIssues
		if e.Resolution != "" {
			notes = e.Resolution
		}
		fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%s\n", e.ID, e.Area, e.Feature, e.Status, notes)
	}
	return nil
}

func writeSummaryBlocks(out io.Writer, summary qaSummary, issues []string) {
	if len(issues) > 0 {
		fmt.Fprintf(out, "Tracker issues: %d\n", len(issues))
		for _, issue := range issues {
			fmt.Fprintf(out, "- %s\n", issue)
		}
		fmt.Fprintf(out, "\n")
	}

	stats := summarizeQAStatus(summary.total, summary.byStatus)
	fmt.Fprintf(out, "- Total stories: %d\n", summary.total)
	fmt.Fprintf(out, "- Pass rate: %d%%\n", stats.PassRate)
	for _, status := range sortedKeys(summary.byStatus) {
		fmt.Fprintf(out, "- %s: %d\n", status, summary.byStatus[status])
	}
}

func renderQAJSON(entries []qaEntry, summary qaSummary, issues []string, out io.Writer) error {
	payload := struct {
		GeneratedAt string        `json:"generatedAt"`
		Stories     []qaEntry     `json:"stories"`
		Summary     qaJSONSummary `json:"summary"`
		Issues      []string      `json:"issues,omitempty"`
	}{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Stories:     entries,
		Issues:      issues,
		Summary:     toJSONSummary(summary),
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func renderQAValidation(entries []qaEntry, summary qaSummary, issues []string, jsonOut bool, out io.Writer) error {
	if jsonOut {
		if len(entries) == 0 {
			issues = append(issues, "tracker is empty")
		}
		report := struct {
			GeneratedAt string        `json:"generatedAt"`
			Validation  qaJSONSummary `json:"validation"`
			Issues      []string      `json:"issues"`
		}{
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			Validation:  toJSONSummary(summary),
			Issues:      issues,
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			return err
		}
		if len(issues) > 0 {
			return fmt.Errorf("qa tracker validation failed: %d issues", len(issues))
		}
		return nil
	}

	fmt.Fprintln(out, "Tracker validation report")
	fmt.Fprintf(out, "Stories: %d\n", summary.total)
	if len(issues) == 0 {
		fmt.Fprintln(out, "Status: PASS")
		return nil
	}

	fmt.Fprintln(out, "Status: FAIL")
	for _, issue := range issues {
		fmt.Fprintf(out, "- %s\n", issue)
	}
	return fmt.Errorf("qa tracker validation failed: %d issues", len(issues))
}

func toJSONSummary(summary qaSummary) qaJSONSummary {
	sv := summarizeQAStatus(summary.total, summary.byStatus)
	return qaJSONSummary{
		Total:        summary.total,
		BlockedItems: summary.blockedItems,
		PassRate:     sv.PassRate,
		ByStatus:     summary.byStatus,
		ByArea:       summary.byArea,
	}
}

type qaJSONSummary struct {
	Total        int            `json:"total"`
	PassRate     int            `json:"passRate"`
	BlockedItems int            `json:"blockedItems"`
	ByStatus     map[string]int `json:"byStatus"`
	ByArea       map[string]int `json:"byArea"`
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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

func isJSONRequested(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if cmd.Root() != nil {
		if b, err := cmd.Root().PersistentFlags().GetBool("json"); err == nil {
			return b
		}
	}
	if cmd.InheritedFlags().Lookup("json") != nil {
		if b, err := cmd.InheritedFlags().GetBool("json"); err == nil {
			return b
		}
	}
	return false
}
