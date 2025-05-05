package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "runtime"

    "github.com/spf13/cobra"
    "github.com/sanix-darker/stackfetch/internal/langfetch"
    "github.com/sanix-darker/stackfetch/internal/sysinfo"
)

type result struct {
    System  sysinfo.Info        `json:"system"`
    Reports []langfetch.Result  `json:"reports"`
}

func main() {
    var jsonOut bool

    root := &cobra.Command{
        Use:   "stackfetch [items…]",
        Short: "System / language / DevOps stack fetcher",
        Args:  cobra.ArbitraryArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            res := result{System: sysinfo.Collect()}
            res.Reports = langfetch.FetchMany(args) // parallel

            if jsonOut {
                enc := json.NewEncoder(os.Stdout)
                enc.SetIndent("", "  ")
                return enc.Encode(res)
            }

            fmt.Println("=== System ===")
            fmt.Println(res.System)
            for _, r := range res.Reports {
                fmt.Printf("\n=== %s ===\n", r.Key)
                if r.Err != nil {
                    fmt.Fprintln(os.Stderr, "stackfetch:", r.Err)
                    continue
                }
                fmt.Println(r.Info)
            }
            return nil
        },
    }

    root.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

    // Enable shell‑completion generation: `stackfetch completion bash` etc.
    root.CompletionOptions.DisableDefaultCmd = false

    if err := root.Execute(); err != nil {
        log.Fatalf("%v", err)
    }

    // Prevent Windows double‑close flash (when built as GUI). On *nix it is NOP.
    if runtime.GOOS == "windows" {
        fmt.Println("Press Enter to exit…")
        fmt.Scanln()
    }
}
