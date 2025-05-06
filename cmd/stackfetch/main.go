package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/sanix-darker/stackfetch/internal/guess"
	"github.com/sanix-darker/stackfetch/internal/langfetch"
	"github.com/sanix-darker/stackfetch/internal/services"
	"github.com/sanix-darker/stackfetch/internal/sysinfo"
	"github.com/spf13/cobra"
)

type result struct {
    System   sysinfo.Info           `json:"system"`
    Reports  []langfetch.Result     `json:"reports"`
    Services []services.Status      `json:"services,omitempty"`
    Guessed  []string               `json:"guessed,omitempty"`
    Ports    []services.PortStatus  `json:"ports,omitempty"`
}

var version = "dev"

func main() {
    var jsonOut bool

    root := &cobra.Command{
        Use:   "stackfetch [items…]",
        Short: "System / language / DevOps stack fetcher",
        Args:  cobra.ArbitraryArgs,
        Version: version,
    }
    // for the --version flag
    root.InitDefaultVersionFlag()
    root.SetVersionTemplate("stackfetch version {{.Version}}\n")

    root.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
    root.RunE = func(cmd *cobra.Command, args []string) error {
        return runStackfetch(jsonOut, nil, args)
    }

    guessCmd := &cobra.Command{
        Use:     "guess",
        Aliases: []string{"?"},
        Short:   "Guess project stack based on files in cwd",
        RunE: func(cmd *cobra.Command, _ []string) error {
            cwd, _ := os.Getwd()
            guessed := guess.Guess(cwd)
            return runStackfetch(jsonOut, guessed, guessed)
        },
    }
    root.AddCommand(guessCmd)

    if err := root.Execute(); err != nil {
        log.Fatalf("%v", err)
    }
	// On Windows, pause so the console doesn’t vanish immediately
    // YES DUDE, THAT'S A THING ON WINDOWS !
    if runtime.GOOS == "windows" {
        fmt.Println("Press Enter to exit…")
        fmt.Scanln()
    }
}

// runStackfetch centralizes both plain-text and JSON output
func runStackfetch(jsonOut bool, guessed, args []string) error {
    res := result{
        System:  sysinfo.Collect(),
        Guessed: guessed,
    }
    // parallel fetch of all requested items
    res.Reports = langfetch.FetchMany(args)

    // build dependency list for all keys we care about
    var deps []string
    for _, key := range append(guessed, args...) {
        deps = append(deps, langfetch.Dependencies(key)...)
    }
    // now check only those services/tools
    res.Services = services.Check(deps)

    if jsonOut {
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        return enc.Encode(res)
    }

    fmt.Println("=== System ===")
    fmt.Println(res.System)

    if len(guessed) > 0 {
        fmt.Printf("\nGuessed: %s\n", filepath.Base("."))
        fmt.Println("Detected items:", guessed)
    }

    for _, r := range res.Reports {
        fmt.Printf("\n=== %s ===\n", r.Key)
        if r.Err != nil {
            fmt.Fprintln(os.Stderr, "stackfetch:", r.Err)
            continue
        }
        fmt.Println(r.Info)
        if depList := langfetch.Dependencies(r.Key); len(depList) > 0 {
            fmt.Println("  └─ depends on:")
            for _, d := range depList {
                st := services.StatusByName(d)
                printServiceLine(d, st)
            }
        }
    }

    // check port status for those same services
    //    500ms is a reasonable per-port timeout
    res.Ports = services.CheckPorts(deps, 500*time.Millisecond)

    // then in your plain-text output:
    fmt.Println("\n=== Ports ===")
    for _, ps := range res.Ports {
        status := "closed"
        if ps.Open {
            status = "open"
        }
        fmt.Printf("  %s:%d → %s\n", ps.Service, ps.Port, status)
    }

    return nil
}

func printServiceLine(name string, s services.Status) {
	var line string
	switch {
	case !s.Installed:
		line = color.RedString("✗ %s (not installed)", name)
	case s.Running:
		line = color.GreenString("✔ %s (running)", name)
	default:
		line = color.YellowString("! %s (installed, not running)", name)
	}
	fmt.Println(line)
}
