package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type qaJSONResponse struct {
	GeneratedAt string        `json:"generatedAt"`
	Summary     qaJSONSummary `json:"summary"`
	Stories     []qaEntry     `json:"stories"`
	Issues      []string      `json:"issues"`
}

func TestParseFeatureTracker(t *testing.T) {
	raw := []byte(`id,area,feature,user_story,expected_behavior,status_after_changes,test_coverage,open_issues,resolution_notes
F-01,Core,Output,Show version,Display version string,PASS,manual,test,all good
F-02,Core,QA,Display matrix,Output stories table,TODO,manual,none,tracked later`)

	entries, issues, err := parseFeatureTracker(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("issues: %v", issues)
	}
	if len(entries) != 2 {
		t.Fatalf("len=%d, want 2", len(entries))
	}
	if entries[0].ID != "F-01" || entries[1].ID != "F-02" {
		t.Fatalf("unexpected order: %+v", entries)
	}
}

func TestParseFeatureTrackerWithMalformedRows(t *testing.T) {
	raw := []byte(`id,area,feature,user_story,expected_behavior,status_after_changes,test_coverage,open_issues,resolution_notes
F-01,Core,Output,Show version,Display version string,PASS,manual,test,all good
F-02,Core,BadRow,Missing fields,manual,test
F-03,Core,Tool,Status typo,Check status,UNKNOWN,manual,test,
F-04,Core,TooMany,One,Two,Three,Four,Five,Six,Seven,Eight`)

	entries, issues, err := parseFeatureTracker(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len=%d, want 2", len(entries))
	}
	if len(issues) < 2 {
		t.Fatalf("issues=%d, want at least 2", len(issues))
	}
	if entries[0].ID != "F-01" || entries[1].ID != "F-03" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestFilterQAEntries(t *testing.T) {
	entries := []qaEntry{
		{Area: "Core", Status: "PASS"},
		{Area: "Services", Status: "todo"},
		{Area: "core", Status: "FAIL"},
	}

	filtered := filterQAEntries(entries, parseAreaFilter("core"), parseStatusFilter("pass,fail"))
	if len(filtered) != 2 {
		t.Fatalf("len=%d, want 2", len(filtered))
	}

	filtered = filterQAEntries(entries, nil, parseStatusFilter("fail"))
	if len(filtered) != 1 {
		t.Fatalf("len=%d, want 1", len(filtered))
	}

	filtered = filterQAEntries(entries, nil, nil)
	if len(filtered) != 3 {
		t.Fatalf("len=%d, want 3", len(filtered))
	}
}

func TestFilterQAEntriesWildcardArea(t *testing.T) {
	entries := []qaEntry{
		{Area: "Core Infra", Status: "PASS"},
		{Area: "core runtime", Status: "TODO"},
		{Area: "QA", Status: "PASS"},
	}

	filtered := filterQAEntries(entries, parseAreaFilter("core*"), parseStatusFilter("pass"))
	if len(filtered) != 1 {
		t.Fatalf("len=%d, want 1", len(filtered))
	}
}

func TestCountMatchingStatusItems(t *testing.T) {
	entries := []qaEntry{
		{Status: "PASS"},
		{Status: "TODO"},
		{Status: "FAIL"},
		{Status: "todo"},
	}
	count := countMatchingStatusItems(entries, parseStatusFilter("TODO,FAIL"))
	if count != 3 {
		t.Fatalf("count=%d, want 3", count)
	}
}

func TestSummarizeQA(t *testing.T) {
	summary := summarizeQA([]qaEntry{
		{Area: "Core", Status: "PASS"},
		{Area: "Core", Status: "TODO"},
		{Area: "QA", Status: "FAIL"},
	})
	if summary.total != 3 {
		t.Fatalf("total=%d", summary.total)
	}
	if summary.byStatus["PASS"] != 1 || summary.byStatus["TODO"] != 1 || summary.byStatus["FAIL"] != 1 {
		t.Fatalf("unexpected status counts: %+v", summary.byStatus)
	}
	if summary.blockedItems != 2 {
		t.Fatalf("blockedItems=%d", summary.blockedItems)
	}

	status := summarizeQAStatus(summary.total, summary.byStatus)
	if status.PassRate != 33 {
		t.Fatalf("pass rate=%d", status.PassRate)
	}
}

func TestValidateFeatureTracker(t *testing.T) {
	entries := []qaEntry{
		{line: 2, ID: "F-01", Area: "Core", Feature: "Output", UserStory: "Show", Expected: "Show", Status: "PASS", TestCoverage: "manual"},
		{line: 3, ID: "F-01", Area: "Core", Feature: "Dup", UserStory: "Show", Expected: "Show", Status: "UNKNOWN", TestCoverage: ""},
		{line: 4, ID: "", Area: "Core", Feature: "Missing", UserStory: "", Expected: "", Status: "", TestCoverage: "manual"},
	}

	issues := validateFeatureTracker(entries)
	if len(issues) < 4 {
		t.Fatalf("issues=%d, want at least 4", len(issues))
	}
	joined := strings.Join(issues, "\n")
	for _, needle := range []string{"duplicate id", "unknown status", "missing id", "missing user story"} {
		if !strings.Contains(strings.ToLower(joined), needle) {
			t.Fatalf("expected issue %q in %q", needle, joined)
		}
	}
}

func TestRenderQAMatrixJSON(t *testing.T) {
	entries := []qaEntry{
		{ID: "F-01", Area: "Core", Feature: "Output", Status: "PASS", TestCoverage: "manual"},
		{ID: "F-02", Area: "Core", Feature: "QA", Status: "TODO", TestCoverage: "manual", Resolution: "later"},
	}
	summary := summarizeQA(entries)

	var out strings.Builder
	if err := renderQAMatrix(entries, summary, []string{"issue: x"}, false, true, &out); err != nil {
		t.Fatalf("render json: %v", err)
	}

	decoded := qaJSONResponse{}

	if err := json.Unmarshal([]byte(out.String()), &decoded); err != nil {
		t.Fatalf("decode json: %v\nraw=%q", err, out.String())
	}
	if decoded.Summary.Total != 2 {
		t.Fatalf("summary total=%d", decoded.Summary.Total)
	}
	if len(decoded.Stories) != 2 {
		t.Fatalf("stories=%d", len(decoded.Stories))
	}
	if decoded.Issues[0] != "issue: x" {
		t.Fatalf("issues=%v", decoded.Issues)
	}
}

func TestParseStatusFilter(t *testing.T) {
	filters := parseStatusFilter("pass, todo,FAIL")
	if _, ok := filters["PASS"]; !ok {
		t.Fatalf("missing PASS")
	}
	if _, ok := filters["TODO"]; !ok {
		t.Fatalf("missing TODO")
	}
	if _, ok := filters["FAIL"]; !ok {
		t.Fatalf("missing FAIL")
	}
	if _, ok := filters["MISS"]; ok {
		t.Fatalf("unexpected status")
	}
}

func TestParseAreaFilter(t *testing.T) {
	filters := parseAreaFilter("core, QA, core*")
	if _, ok := filters["core"]; !ok {
		t.Fatalf("missing core")
	}
	if _, ok := filters["qa"]; !ok {
		t.Fatalf("missing qa")
	}
	if _, ok := filters["core*"]; !ok {
		t.Fatalf("missing wildcard core*")
	}
}

func TestLoadFeatureTrackerCSV(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(workingDir)

	tmpDir := t.TempDir()
	tracker := filepath.Join(tmpDir, "FEATURE_TRACKER.csv")
	payload := "id,area,feature,user_story,expected_behavior,status_after_changes,test_coverage,open_issues,resolution_notes\nF-LOCAL,Core,Output,Local story,Local check,PASS,manual,,\n"
	if err := os.WriteFile(tracker, []byte(payload), 0o644); err != nil {
		t.Fatalf("write temp tracker: %v", err)
	}

	got, err := loadFeatureTrackerCSV(tracker)
	if err != nil {
		t.Fatalf("load tracker: %v", err)
	}
	if string(got) != payload {
		t.Fatalf("unexpected tracker content")
	}

	if _, err := loadFeatureTrackerCSV(filepath.Join(tmpDir, "missing.csv")); err == nil {
		t.Fatalf("expected error when tracker path does not exist")
	}
}

func TestQACmdFailOnStatus(t *testing.T) {
	payload := "id,area,feature,user_story,expected_behavior,status_after_changes,test_coverage,open_issues,resolution_notes\n" +
		"F-01,Core,Output,Local story,Local check,PASS,manual,,\n" +
		"F-02,Core,Output,Local story,Local check,TODO,manual,,\n"
	tmpDir := t.TempDir()
	tracker := filepath.Join(tmpDir, "FEATURE_TRACKER.csv")
	if err := os.WriteFile(tracker, []byte(payload), 0o644); err != nil {
		t.Fatalf("write temp tracker: %v", err)
	}

	cmd := newQACmd(false)
	cmd.SetArgs([]string{"--tracker", tracker, "--fail-on", "todo"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected qa gate failure")
	}
}

func TestQACmdStrictUsesDefaultGate(t *testing.T) {
	payload := "id,area,feature,user_story,expected_behavior,status_after_changes,test_coverage,open_issues,resolution_notes\n" +
		"F-01,Core,Output,Local story,Local check,PASS,manual,,\n" +
		"F-02,Core,Output,Local story,Local check,FAIL,manual,,\n"
	tmpDir := t.TempDir()
	tracker := filepath.Join(tmpDir, "FEATURE_TRACKER.csv")
	if err := os.WriteFile(tracker, []byte(payload), 0o644); err != nil {
		t.Fatalf("write temp tracker: %v", err)
	}

	cmd := newQACmd(false)
	cmd.SetArgs([]string{"--tracker", tracker, "--strict", "true"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected strict fail on FAIL status")
	}
}
