package models

type AvdInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	DiskUsage string `json:"diskUsage"`
	Running   bool   `json:"running"`
}
