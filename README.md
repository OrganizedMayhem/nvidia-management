# nvidia-management

A CLI utility for managing NVIDIA systemd services on AMD CPU/GPU systems that use an external NVIDIA GPU.

## Why?

On systems where the primary CPU and GPU are AMD, NVIDIA's systemd services (power management, suspend/resume hooks, etc.) are unnecessary overhead when no external NVIDIA card is connected. However, when you plug in an external NVIDIA GPU (e.g. via Thunderbolt/USB4 eGPU enclosure), those services need to be running for proper power management and suspend/resume behavior.

This tool gives you a single command to toggle all NVIDIA services on or off, rather than managing each unit individually with `systemctl`.

## Managed Services

- `nvidia-hibernate.service`
- `nvidia-powerd.service`
- `nvidia-resume.service`
- `nvidia-suspend-then-hibernate.service`
- `nvidia-suspend.service`

## Usage

Requires root privileges (`sudo`).

```
sudo nvidia-management on       # Enable and start all NVIDIA services
sudo nvidia-management off      # Stop and disable all NVIDIA services
sudo nvidia-management status   # Show current state of all NVIDIA services
```

## Building

```
go build .
```
