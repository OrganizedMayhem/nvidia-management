// Package systemd provides helpers for querying and controlling systemd services
// via the systemctl binary. It intentionally avoids D-Bus bindings to keep the
// dependency surface minimal.
package systemd

import (
	"os/exec"
)

// IsActive returns true when the service is currently active (running).
func IsActive(svc string) bool {
	return exec.Command("systemctl", "is-active", "--quiet", svc).Run() == nil
}

// IsEnabled returns true when the service is enabled (starts at boot).
func IsEnabled(svc string) bool {
	return exec.Command("systemctl", "is-enabled", "--quiet", svc).Run() == nil
}

// Enable enables the service (symlinks it into the appropriate target).
func Enable(svc string) error {
	return exec.Command("systemctl", "enable", svc).Run()
}

// Disable disables the service.
func Disable(svc string) error {
	return exec.Command("systemctl", "disable", svc).Run()
}

// Start starts the service immediately (without enabling it).
func Start(svc string) error {
	return exec.Command("systemctl", "start", svc).Run()
}

// Stop stops the service immediately (without disabling it).
func Stop(svc string) error {
	return exec.Command("systemctl", "stop", svc).Run()
}
