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
	root.AddCommand(newQACmd(mdOut))
	root.AddCommand(newKeysCmd())

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
	printSectionHeader("System")
	fmt.Println(res.System)

	if len(guessed) > 0 {
		fmt.Printf("\n%s Guessed project context: %s\n", ui.CyanString(ui.Icon("lightbulb")), filepath.Base("."))
		fmt.Printf("%s %s\n", ui.GrayString("Detected:"), strings.Join(guessed, ", "))
	}

	for _, r := range res.Reports {
		if r.Err != nil {
			if len(guessed) == 0 {
				fmt.Fprintf(os.Stderr, "stackfetch: %s: %v\n", ui.RedString(r.Key), r.Err)
			}
			continue
		}

		fmt.Printf("\n%s\n", sectionHeader(r.Key))
		fmt.Println(r.Info)
		if depList := langfetch.Dependencies(r.Key); len(depList) > 0 {
			fmt.Printf("%s Dependency health\n", ui.GrayString("  "+ui.Icon("branch")))
			deps := collectDependencyStatus(depList)
			printDependencySummary(deps)
			fmt.Printf("%s depends on:\n", ui.GrayString("  "+ui.Icon("leaf")))
			for _, st := range deps {
				printServiceLine(st.Name, st)
			}
			fmt.Printf("%s\n", readinessBadge(readinessFromStatuses(deps)))
		}
	}

	if res.Cloud.Provider != "" {
		printSectionHeader("Cloud")
		fmt.Printf("%s %s\n", ui.CyanString(ui.Icon("satellite")+"  provider"), ui.GreenString(res.Cloud.Provider))
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

	printSectionHeader("Security")
	fmt.Printf("  %s %s %s\n", securityIndicator(res.Security.Root), ui.GreenString("Root user:"), boolState(res.Security.Root))
	fmt.Printf("  %s %s %s\n", ui.CyanString(ui.Icon("lock")), ui.GreenString("SELinux:"), securityLabel(res.Security.SELinux))
	fmt.Printf("  %s %s %s\n", securityIndicator(!res.Security.SSHOpen), ui.GreenString("SSH/22: blocked:"), boolState(!res.Security.SSHOpen))
	fmt.Printf("  %s %s %s\n", securityIndicator(!res.Security.KernelEOL), ui.GreenString("Kernel EOL:"), boolState(!res.Security.KernelEOL))
	fmt.Printf("  %s %s %s\n", ui.GrayString(ui.Icon("flag")), ui.GreenString("Posture score:"), riskProfile(res.Security))

	if len(res.Ports) > 0 {
		printPortMatrix(res.Ports)
	}
	if len(res.Services) > 0 {
		printSectionHeader("Readiness")
		fmt.Printf("%s %s\n", ui.Icon("shield"), readinessBanner(readinessFromStatuses(res.Services)))
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
			deps := collectDependencyStatus(depList)
			installed, running := dependencySummary(deps)
			fmt.Printf("Depends on: %d/%d installed, %d/%d running\n", installed, len(depList), running, len(depList))
			for _, st := range deps {
				fmt.Printf(" - %s [%s / %s]\n", st.Name, depState(st.Installed, true), depState(st.Running, true))
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
	fmt.Printf("Posture score: %s\n", riskProfile(res.Security))
	fmt.Printf("Root: %t\n", res.Security.Root)
	fmt.Printf("SELinux: %s\n", res.Security.SELinux)
	fmt.Printf("SSH open on 22: %t\n", res.Security.SSHOpen)
	fmt.Printf("Kernel EOL check: %t\n", res.Security.KernelEOL)
	fmt.Println()

	if len(res.Ports) > 0 {
		ui.Heading("Ports", 2)
		fmt.Printf("Summary: %s\n", readinessPortSummary(res.Ports))
		openCount := 0
		for _, ps := range res.Ports {
			if ps.Open {
				openCount++
			}
		}
		fmt.Printf("Total open: %d / %d\n", openCount, len(res.Ports))
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
	fmt.Printf("    %s %s\n", depGlyph(s), depStatusLine(name, s))
}

func printSectionHeader(title string) {
	fmt.Printf("\n=== %s ===\n", ui.BlueString(title))
}

func sectionHeader(title string) string {
	return ui.CyanString(fmt.Sprintf("=== %s ===", title))
}

func collectDependencyStatus(depList []string) []services.Status {
	statusList := make([]services.Status, 0, len(depList))
	for _, dep := range depList {
		statusList = append(statusList, services.StatusByName(dep))
	}
	return statusList
}

func dependencySummary(statuses []services.Status) (installed int, running int) {
	for _, st := range statuses {
		if st.Installed {
			installed++
		}
		if st.Running {
			running++
		}
	}
	return
}

func printDependencySummary(statuses []services.Status) {
	installed, running := dependencySummary(statuses)
	total := len(statuses)
	installedBadge := ui.GreenString("%d/%d", installed, total)
	runningBadge := ui.YellowString("%d/%d", running, total)
	depOk := ui.GreenString(ui.Icon("check"))
	depWarn := ui.YellowString(ui.Icon("warning"))

	if total == 0 {
		fmt.Printf("  %s deps: none declared\n", ui.GrayString("◦"))
		return
	}

	if installed == total && running == total {
		fmt.Printf("  %s installed %s, running %s\n", depOk, installedBadge, ui.GreenString("%d/%d", running, total))
		return
	}

	if installed == total {
		fmt.Printf("  %s installed %s, running %s\n", depWarn, installedBadge, runningBadge)
		return
	}

	fmt.Printf("  %s installed %s, running %s\n", depWarn, installedBadge, runningBadge)
}

func depGlyph(s services.Status) string {
	switch {
	case !s.Installed:
		return ui.RedString(ui.Icon("cross"))
	case s.Running:
		return ui.GreenString(ui.Icon("running"))
	default:
		return ui.YellowString(ui.Icon("partial"))
	}
}

func depStatusLine(name string, s services.Status) string {
	if !s.Installed {
		return ui.RedString("%s (not installed)", name)
	}
	if s.Running {
		return ui.GreenString("%s (running)", name)
	}
	return ui.YellowString("%s (installed, not running)", name)
}

func securityLabel(v string) string {
	switch strings.ToLower(v) {
	case "enforcing":
		return ui.GreenString(v)
	case "disabled":
		return ui.RedString(v)
	case "permissive":
		return ui.YellowString(v)
	default:
		return ui.GrayString(v)
	}
}

func boolState(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func securityIndicator(_value bool) string {
	if _value {
		return ui.GreenString(ui.Icon("check"))
	}
	return ui.RedString(ui.Icon("cross"))
}

type readiness struct {
	installed int
	running   int
	total     int
	score     int
}

func readinessFromStatuses(statuses []services.Status) readiness {
	r := readiness{}
	r.total = len(statuses)
	if r.total == 0 {
		return r
	}
	for _, st := range statuses {
		if st.Installed {
			r.installed++
		}
		if st.Running {
			r.running++
		}
	}
	r.score = int(((2*r.installed + r.running) * 100) / (3 * r.total))
	return r
}

func readinessBadge(r readiness) string {
	if r.total == 0 {
		return ui.GrayString("Readiness: independent")
	}
	note := fmt.Sprintf("Readiness: %d%% (installed %d/%d, running %d/%d)", r.score, r.installed, r.total, r.running, r.total)
	if r.score >= 90 {
		return ui.GreenString(note)
	}
	if r.score >= 65 {
		return ui.YellowString(note)
	}
	return ui.RedString(note)
}

func readinessBanner(r readiness) string {
	if r.total == 0 {
		return "independent"
	}
	if r.score >= 90 {
		return fmt.Sprintf("excellent (%d%%)", r.score)
	}
	if r.score >= 65 {
		return fmt.Sprintf("warning (%d%%)", r.score)
	}
	return fmt.Sprintf("critical (%d%%)", r.score)
}

func readinessPortSummary(ports []services.PortStatus) string {
	if len(ports) == 0 {
		return "ports: none tracked"
	}
	open := 0
	for _, ps := range ports {
		if ps.Open {
			open++
		}
	}
	return fmt.Sprintf("%d/%d ports open", open, len(ports))
}

func riskProfile(sec security.Report) string {
	risk := 0
	if sec.KernelEOL {
		risk++
	}
	if sec.SSHOpen {
		risk++
	}
	if strings.EqualFold(sec.SELinux, "disabled") {
		risk++
	}

	if risk == 0 {
		return ui.GreenString("Green")
	}
	if risk == 1 {
		return ui.YellowString("Moderate")
	}
	return ui.RedString("High")
}

func depState(v bool, positive bool) string {
	if v == positive {
		return "+"
	}
	return "-"
}

func printPortMatrix(ports []services.PortStatus) {
	printSectionHeader("Ports")
	openCount := 0
	for _, ps := range ports {
		if ps.Open {
			openCount++
		}
	}
	fmt.Printf("  %s %d/%d open\n", ui.GrayString(ui.Icon("goal")+" exposure"), openCount, len(ports))
	for _, ps := range ports {
		status := ui.RedString("closed")
		glyph := ui.RedString(ui.Icon("blank"))
		if ps.Open {
			status = ui.GreenString("open")
			glyph = ui.GreenString(ui.Icon("running"))
		}
		fmt.Printf("  %s %s:%d %s %s\n", glyph, ps.Service, ps.Port, ui.Icon("arrow"), status)
	}
}
