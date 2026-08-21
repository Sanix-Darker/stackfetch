package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sanix-darker/stackfetch/internal/langfetch"
	"github.com/spf13/cobra"
)

const keyInfoSuggestionLimit = 5

type keyInfoResponse struct {
	Key          string              `json:"key"`
	Supported    bool                `json:"supported"`
	Dependencies []string            `json:"dependencies,omitempty"`
	Suggestions  []string            `json:"suggestions,omitempty"`
	Fetch        bool                `json:"fetch"`
	Error        string              `json:"error,omitempty"`
	Info         *langfetch.LangInfo `json:"info,omitempty"`
}

func newKeyInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keyinfo <key>",
		Short: "Inspect registry metadata and dependencies for a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fetch, _ := cmd.Flags().GetBool("fetch")
			key := strings.ToLower(strings.TrimSpace(args[0]))
			resp := buildKeyInfoResponse(key, fetch)

			if isJSONRequested(cmd) {
				enc := json.NewEncoder(cmd.OutOrStdout())
				return enc.Encode(resp)
			}
			return writeKeyInfoOutput(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().Bool("fetch", false, "Also collect live report for this key")
	return cmd
}

func buildKeyInfoResponse(key string, fetch bool) keyInfoResponse {
	key = strings.ToLower(strings.TrimSpace(key))
	resp := keyInfoResponse{
		Key:         key,
		Fetch:       fetch,
		Supported:   langfetch.IsSupported(key),
		Suggestions: langfetch.SuggestKeys(key, keyInfoSuggestionLimit),
	}

	if resp.Supported {
		resp.Dependencies = langfetch.Dependencies(key)
	}

	if fetch {
		info, err := langfetch.Fetch(key)
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Info = &info
		}
	}

	return resp
}

func writeKeyInfoOutput(out io.Writer, resp keyInfoResponse) error {
	if resp.Supported {
		fmt.Fprintf(out, "key: %s\n", resp.Key)
		fmt.Fprintln(out, "supported: true")
		if len(resp.Dependencies) > 0 {
			fmt.Fprintf(out, "dependencies: %s\n", strings.Join(resp.Dependencies, ", "))
		} else {
			fmt.Fprintln(out, "dependencies: none")
		}
	} else {
		fmt.Fprintf(out, "key: %s\n", resp.Key)
		fmt.Fprintln(out, "supported: false")
	}

	if len(resp.Error) > 0 {
		fmt.Fprintf(out, "error: %s\n", resp.Error)
	}
	if resp.Supported && resp.Info != nil && resp.Fetch {
		fmt.Fprintln(out, "info:")
		fmt.Fprintln(out, resp.Info.String())
	}
	if len(resp.Suggestions) > 0 && !resp.Supported {
		fmt.Fprintf(out, "suggestions: %s\n", strings.Join(resp.Suggestions, ", "))
	}
	return nil
}
