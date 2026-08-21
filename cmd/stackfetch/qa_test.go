package main

import (
	"testing"
)

func TestParseFeatureTracker(t *testing.T) {
	raw := []byte(`id,area,feature,user_story,expected_behavior,status_after_changes,test_coverage,open_issues,resolution_notes
F-01,Core,Output,Show version,Display version string,PASS,manual,test,all good
F-02,Core,QA,Display matrix,Output stories table,TODO,manual,none,tracked later`)

	entries, err := parseFeatureTracker(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len=%d, want 2", len(entries))
	}
	if entries[0].ID != "F-01" || entries[1].Status != "TODO" {
		t.Fatalf("unexpected parse order: %+v", entries)
	}
}

func TestFilterQAEntries(t *testing.T) {
	entries := []qaEntry{
		{Area: "Core"},
		{Area: "Services"},
		{Area: "core"},
	}

	filtered := filterQAEntries(entries, "core")
	if len(filtered) != 2 {
		t.Fatalf("len=%d, want 2", len(filtered))
	}

	filtered = filterQAEntries(entries, "")
	if len(filtered) != 3 {
		t.Fatalf("len=%d, want 3", len(filtered))
	}
}

func TestSummarizeQA(t *testing.T) {
	entries := []qaEntry{
		{Area: "Core", Status: "PASS"},
		{Area: "Core", Status: "TODO"},
		{Area: "QA", Status: "FAIL"},
	}

	s := summarizeQA(entries)
	if s.total != 3 {
		t.Fatalf("total=%d", s.total)
	}
	if s.byStatus["PASS"] != 1 || s.byStatus["TODO"] != 1 || s.byStatus["FAIL"] != 1 {
		t.Fatalf("unexpected status counts: %+v", s.byStatus)
	}
	if s.blockedItems != 2 {
		t.Fatalf("blockedItems=%d", s.blockedItems)
	}
}
