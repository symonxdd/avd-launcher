package models

type AvdInfo struct {
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	Path           string `json:"path"`
	DiskUsage      string `json:"diskUsage"`
	Running        bool   `json:"running"`
	ApiLevel        string `json:"apiLevel"`
	AndroidVersion  string `json:"androidVersion"`
	AndroidCodename string `json:"androidCodename"`
	Abi            string `json:"abi"`
	RamSize        string `json:"ramSize"`
	Resolution     string `json:"resolution"`
	HasGooglePlay  bool   `json:"hasGooglePlay"`
}
