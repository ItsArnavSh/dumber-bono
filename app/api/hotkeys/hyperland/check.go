package hyperland

import (
	"os"
	"strings"
)

func IsUsingHyperland() bool {
	// Hyprland sets XDG_CURRENT_DESKTOP to "Hyprland" (case-insensitive check for safety)
	desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	if strings.EqualFold(desktop, "Hyprland") {
		return true
	}

	// Fallback check: HYPRLAND_INSTANCE_SIGNATURE is always present in a Hyprland session
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return true
	}
	return false
}
