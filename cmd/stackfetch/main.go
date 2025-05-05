package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sanix-darker/stackfetch/internal/guess"
	"github.com/sanix-darker/stackfetch/internal/langfetch"
	"github.com/sanix-darker/stackfetch/internal/sysinfo"
	"github.com/spf13/cobra"
)

type result struct {
	System  sysinfo.Info       `json:"system"`
	Reports []langfetch.Result `json:"reports"`
	Guessed []string           `json:"guessed,omitempty"`
}

func main() {
	var jsonOut bool

	root := &cobra.Command{
		Use:   "stackfetch [items…]",
		Short: "System / language / DevOps stack fetcher",
		Args:  cobra.ArbitraryArgs,
	}

	root.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	// Default run: takes positional args
	root.RunE = func(cmd *cobra.Command, args []string) error {
		var jsonOut bool = jsonOut
		res := result{System: sysinfo.Collect(), Guessed: []string(nil)}
		res.Reports = langfetch.FetchMany(args)
		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		}
		fmt.Println("=== System ===")
		fmt.Println(res.System)
		if len([]string(nil)) > 0 {
			fmt.Printf("\nGuessed: %s\n", filepath.Base("."))
			fmt.Println("Detected items:", []string(nil))
		}
		for _, r := range res.Reports {
			fmt.Printf("\n=== %s ===\n", r.Key)
			if r.Err != nil {
				fmt.Fprintln(os.Stderr, "stackfetch:", r.Err)
				continue
			}
			fmt.Println(r.Info)
		}
		return nil
	}

	// Guess sub-command (alias “?”)
	guessCmd := &cobra.Command{
		Use:     "guess",
		Aliases: []string{"?"},
		Short:   "Guess project stack based on files in cwd",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cwd, _ := os.Getwd()
			guessed := guess.Guess(cwd)
			var (
				jsonOut bool     = jsonOut
				args    []string = guessed
			)
			res := result{System: sysinfo.Collect(), Guessed: guessed}
			res.Reports = langfetch.FetchMany(args)
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
			}
			return nil
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
	res.Reports = langfetch.FetchMany(args)

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
	}
	return nil
}
