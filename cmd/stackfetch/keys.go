package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/sanix-darker/stackfetch/internal/langfetch"
	"github.com/spf13/cobra"
)

func newKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "List known language and stack keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern, _ := cmd.Flags().GetString("contains")
			sortBy, _ := cmd.Flags().GetString("sort")
			sortBy = strings.ToLower(strings.TrimSpace(sortBy))
			switch sortBy {
			case "", "asc", "desc":
			default:
				return fmt.Errorf("invalid --sort value %q: use asc or desc", sortBy)
			}
			keys := langfetch.Keys()
			filtered := filterStringSlice(keys, pattern, sortBy)
			if isJSONRequested(cmd) {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(struct {
					Count int      `json:"count"`
					Keys  []string `json:"keys"`
				}{
					Count: len(filtered),
					Keys:  filtered,
				})
			}

			for _, key := range filtered {
				fmt.Fprintln(cmd.OutOrStdout(), key)
			}
			return nil
		},
	}

	cmd.Flags().String("contains", "", "Filter keys by substring (case-insensitive)")
	cmd.Flags().String("sort", "asc", "Sort order: asc or desc")
	return cmd
}

func filterStringSlice(values []string, contains, sortOrder string) []string {
	sortOrder = strings.ToLower(strings.TrimSpace(sortOrder))
	desc := sortOrder == "desc"

	if contains == "" {
		if !desc {
			return values
		}
		out := append([]string{}, values...)
		sort.SliceStable(out, func(i, j int) bool {
			return out[i] > out[j]
		})
		return out
	}

	needle := strings.ToLower(strings.TrimSpace(contains))
	filtered := make([]string, 0, len(values))
	for _, v := range values {
		if strings.Contains(strings.ToLower(v), needle) {
			filtered = append(filtered, v)
		}
	}
	if desc {
		sort.SliceStable(filtered, func(i, j int) bool {
			return filtered[i] > filtered[j]
		})
	}
	return filtered
}
