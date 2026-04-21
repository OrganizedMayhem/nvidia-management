package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/OrganizedMayhem/nvidia-management/internal/color"
	"github.com/OrganizedMayhem/nvidia-management/internal/systemd"
	"github.com/spf13/cobra"
)

// services is the canonical list of NVIDIA systemd units managed by this tool.
var services = []string{
	"nvidia-hibernate.service",
	"nvidia-powerd.service",
	"nvidia-resume.service",
	"nvidia-suspend-then-hibernate.service",
	"nvidia-suspend.service",
}

// colour helpers
var (
	bold   = color.New(color.Bold)
	cyan   = color.New(color.FgCyan, color.Bold)
	green  = color.New(color.FgGreen)
	yellow = color.New(color.FgYellow)
	red    = color.New(color.FgRed)
	dim    = color.New(color.Faint)
)

var rootCmd = &cobra.Command{
	Use:   "nvidia-management",
	Short: "Manage NVIDIA systemd services",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if os.Getuid() != 0 {
			red.Fprintf(os.Stderr, "✗  Must be run as root (sudo).\n")
			os.Exit(1)
		}
	},
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable and start all NVIDIA services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		cyan.Println("── Enabling & starting NVIDIA services ──")
		fmt.Println()

		for _, svc := range services {
			doAction(svc, "enable", systemd.IsEnabled, systemd.Enable, "already enabled", "Enabled")
		}
		fmt.Println()
		for _, svc := range services {
			doAction(svc, "start", systemd.IsActive, systemd.Start, "already running", "Started")
		}
		fmt.Println()
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Stop and disable all NVIDIA services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		cyan.Println("── Stopping & disabling NVIDIA services ──")
		fmt.Println()

		for _, svc := range services {
			doAction(svc, "stop", func(s string) bool { return !systemd.IsActive(s) }, systemd.Stop, "already stopped", "Stopped")
		}
		fmt.Println()
		for _, svc := range services {
			doAction(svc, "disable", func(s string) bool { return !systemd.IsEnabled(s) }, systemd.Disable, "already disabled", "Disabled")
		}
		fmt.Println()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current state of all NVIDIA services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		cyan.Println("── NVIDIA service status ──")
		fmt.Println()

		bold.Printf("  %-45s %-10s %s\n", "Service", "Active", "Enabled")
		fmt.Printf("  %s\n", strings.Repeat("─", 65))

		for _, svc := range services {
			active := systemd.IsActive(svc)
			enabled := systemd.IsEnabled(svc)

			activeStr, aColor := "inactive", dim
			if active {
				activeStr, aColor = "active", green
			}

			enabledStr, eColor := "disabled", dim
			if enabled {
				enabledStr, eColor = "enabled", green
			}

			fmt.Printf("  %-45s ", svc)
			aColor.Printf("%-10s", activeStr)
			fmt.Print(" ")
			eColor.Printf("%s", enabledStr)
			fmt.Println()
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(onCmd, offCmd, statusCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// doAction is a generic helper for enable/disable/start/stop operations.
func doAction(svc, verb string, skip func(string) bool, act func(string) error, skipMsg, doneMsg string) {
	if skip(svc) {
		yellow.Printf("  ⚠  %s %s — skipping\n", svc, skipMsg)
		return
	}
	if err := act(svc); err != nil {
		red.Printf("  ✗  Failed to %s %s: %v\n", verb, svc, err)
		return
	}
	green.Printf("  ✔  %-9s %s\n", doneMsg, svc)
}
