package models

type AccelInfo struct {
	Status     string `json:"status"`     // "available" or "unavailable"
	Hypervisor string `json:"hypervisor"` // e.g. "WHPX", "KVM", "Hypervisor.Framework"
	Details    string `json:"details"`    // Human-readable detail line
}
