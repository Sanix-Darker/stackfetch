package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sanix-darker/stackfetch/internal/cloudmeta"
	"github.com/sanix-darker/stackfetch/internal/containerexec"
	"github.com/sanix-darker/stackfetch/internal/guess"
	"github.com/sanix-darker/stackfetch/internal/langfetch"
	"github.com/sanix-darker/stackfetch/internal/security"
	"github.com/sanix-darker/stackfetch/internal/services"
	"github.com/sanix-darker/stackfetch/internal/sysinfo"
	"github.com/sanix-darker/stackfetch/internal/ui"
	"github.com/spf13/cobra"
)

type result struct {
	System   sysinfo.Info          `json:"system"`
	Reports  []langfetch.Result    `json:"reports"`
	Services []services.Status     `json:"services,omitempty"`
	Guessed  []string              `json:"guessed,omitempty"`
	Ports    []services.PortStatus `json:"ports,omitempty"`
	Cloud    cloudmeta.Info        `json:"cloud,omitempty"`
	Security security.Report       `json:"security,omitempty"`
}

var version = "dev"

type runOpts struct {
	json, md  bool
	badge     string
	container string
}

func main() {
	var jsonOut, mdOut bool
	var containerID string

	root := &cobra.Command{
		Use:     "stackfetch [items…]",
		Short:   "System / language / DevOps stack fetcher",
		Args:    cobra.ArbitraryArgs,
		Version: version,
	}
	// for the --version flag
	root.InitDefaultVersionFlag()
	root.SetVersionTemplate("stackfetch version {{.Version}}\n")

	root.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "For JSON output")
	root.PersistentFlags().BoolVar(&mdOut, "markdown", false, "For Markdown output")
	root.PersistentFlags().StringVar(&containerID, "container", "", "To be executed inside container")

	root.RunE = func(cmd *cobra.Command, args []string) error {
		return runStackfetch(runOpts{
			json: jsonOut, md: mdOut,
			container: containerID,
		}, nil, args)
	}

	guessCmd := &cobra.Command{
		Use:     "guess",
		Aliases: []string{"?"},
		Short:   "Guess project stack based on files in cwd",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cwd, _ := os.Getwd()
			guessed := guess.Guess(cwd)
			return runStackfetch(runOpts{
				json: jsonOut, md: mdOut,
				container: containerID,
			}, guessed, guessed)
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
func runStackfetch(opt runOpts, guessed, args []string) error {
	// ─── Container exec override ─────────────────────────────────────────────
	rt := containerexec.DetectRuntime()
	if opt.container != "" && rt != "" {
		cid := opt.container
		langfetch.ExecRunner = func(bin string, args ...string) ([]byte, error) {
			out, err := containerexec.Exec(rt, cid, bin, args...)
			return []byte(out), err
		}
	}

	// ─── System / Cloud / Security context ───────────────────────────────────
	var sysInfo sysinfo.Info
	var cloud cloudmeta.Info
	var sec security.Report

	if opt.container == "" { // host
		sysInfo = sysinfo.Collect()
		cloud = cloudmeta.Collect()
		sec = security.Collect()
	} else { // inside container → lightweight, command-based gather
		sysInfo = containerSystem(rt, opt.container)
	}

	res := result{
		System:   sysInfo,
		Guessed:  guessed,
		Cloud:    cloud,
		Security: sec,
	}
	// parallel fetch of all requested items
	res.Reports = langfetch.FetchMany(args)

	// build dependency list for all keys we care about
	var deps []string
	for _, key := range append(guessed, args...) {
		deps = append(deps, langfetch.Dependencies(key)...)
	}
	// Services & ports only make sense on host; skip for pure container mode
	if opt.container == "" {
		res.Services = services.Check(deps)
		res.Ports = services.CheckPorts(deps, 300*time.Millisecond)
	}

	// OUTPUT SECTION  (json / md / plain)
	if opt.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(res)
	}

	if opt.md {
		printMarkdownOutput(res, guessed)
		return nil
	}

	printPlainOutput(res, guessed)
	return nil
}

func printPlainOutput(res result, guessed []string) {
	fmt.Println("=== System ===")
	fmt.Println(res.System)

	if len(guessed) > 0 {
		fmt.Printf("\nGuessed: %s\n", filepath.Base("."))
		fmt.Println("Detected items:", guessed)
	}

	for _, r := range res.Reports {
		if r.Err != nil {
			if len(guessed) == 0 {
				fmt.Fprintf(os.Stderr, "stackfetch: %s: %v\n", r.Key, r.Err)
			}
			continue
		}

		fmt.Printf("\n=== %s ===\n", r.Key)
		fmt.Println(r.Info)
		if depList := langfetch.Dependencies(r.Key); len(depList) > 0 {
			fmt.Println("  └─ depends on:")
			for _, d := range depList {
				st := services.StatusByName(d)
				printServiceLine(d, st)
			}
		}
	}

	if res.Cloud.Provider != "" {
		fmt.Println("\n=== Cloud ===")
		fmt.Printf("  Provider: %s\n", res.Cloud.Provider)
		if res.Cloud.InstanceID != "" {
			fmt.Printf("  Instance: %s\n", res.Cloud.InstanceID)
		}
		if res.Cloud.InstanceType != "" {
			fmt.Printf("  Instance Type: %s\n", res.Cloud.InstanceType)
		}
		if res.Cloud.Region != "" {
			fmt.Printf("  Region: %s\n", res.Cloud.Region)
		}
	}

	fmt.Println("\n=== Security ===")
	fmt.Printf("  Root: %t\n", res.Security.Root)
	fmt.Printf("  SELinux: %s\n", res.Security.SELinux)
	fmt.Printf("  SSH open on 22: %t\n", res.Security.SSHOpen)
	fmt.Printf("  Kernel EOL check: %t\n", res.Security.KernelEOL)

	if len(res.Ports) > 0 {
		fmt.Println("\n=== Ports ===")
		for _, ps := range res.Ports {
			status := "closed"
			if ps.Open {
				status = "open"
			}
			fmt.Printf("  %s:%d → %s\n", ps.Service, ps.Port, status)
		}
	}
}

func printMarkdownOutput(res result, guessed []string) {
	ui.Heading("System", 2)
	fmt.Printf("```text\n%s\n```\n\n", res.System)

	if len(guessed) > 0 {
		ui.Heading("Guess", 2)
		fmt.Printf("Guessed: `%s`\n", filepath.Base("."))
		fmt.Printf("Detected: `%s`\n\n", strings.Join(guessed, ", "))
	}

	for _, r := range res.Reports {
		if r.Err != nil {
			if len(guessed) == 0 {
				ui.Heading(r.Key, 3)
				fmt.Printf("`error`: %v\n\n", r.Err)
			}
			continue
		}
		ui.Heading(r.Key, 3)
		fmt.Printf("```text\n%s\n```\n\n", r.Info)
		if depList := langfetch.Dependencies(r.Key); len(depList) > 0 {
			fmt.Println("depends on:")
			for _, d := range depList {
				st := services.StatusByName(d)
				fmt.Printf(" - %s: installed=%t running=%t\n", d, st.Installed, st.Running)
			}
			fmt.Println()
		}
	}

	if res.Cloud.Provider != "" {
		ui.Heading("Cloud", 2)
		fmt.Printf("Provider: %s\n", res.Cloud.Provider)
		if res.Cloud.InstanceID != "" {
			fmt.Printf("Instance ID: %s\n", res.Cloud.InstanceID)
		}
		if res.Cloud.InstanceType != "" {
			fmt.Printf("Instance Type: %s\n", res.Cloud.InstanceType)
		}
		if res.Cloud.Region != "" {
			fmt.Printf("Region: %s\n", res.Cloud.Region)
		}
		fmt.Println()
	}

	ui.Heading("Security", 2)
	fmt.Printf("Root: %t\n", res.Security.Root)
	fmt.Printf("SELinux: %s\n", res.Security.SELinux)
	fmt.Printf("SSH open on 22: %t\n", res.Security.SSHOpen)
	fmt.Printf("Kernel EOL check: %t\n", res.Security.KernelEOL)
	fmt.Println()

	if len(res.Ports) > 0 {
		ui.Heading("Ports", 2)
		for _, ps := range res.Ports {
			status := "closed"
			if ps.Open {
				status = "open"
			}
			fmt.Printf("%s:%d → %s\n", ps.Service, ps.Port, status)
		}
		fmt.Println()
	}
}

// containerSystem fetches minimal sysinfo from inside a running container.
// It avoids host gopsutil calls and relies only on containerexec.Exec.
func containerSystem(rt containerexec.Runtime, cid string) sysinfo.Info {
	get := func(bin string, args ...string) string {
		out, _ := containerexec.Exec(rt, cid, bin, args...)
		return strings.TrimSpace(out)
	}

	osName := ""
	if etc := get("cat", "/etc/os-release"); etc != "" {
		for _, line := range strings.Split(etc, "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				osName = strings.Trim(line[len("PRETTY_NAME="):], `"`)
				break
			}
		}
	}

	return sysinfo.Info{
		OS:     osName,
		Kernel: get("uname", "-r"),
		Arch:   get("uname", "-m"),
	}
}

func printServiceLine(name string, s services.Status) {
	var line string
	switch {
	case !s.Installed:
		line = ui.RedString("✗ %s (not installed)", name)
	case s.Running:
		line = ui.GreenString("✔ %s (running)", name)
	default:
		line = ui.YellowString("! %s (installed, not running)", name)
	}
	fmt.Println(line)
}
