package listener

type NewBeacon struct {
	Id       uint32 `json:"id"`
	Os       uint32 `json:"os"`
	Arch     uint32 `json:"arch"`
	Sleep    uint32 `json:"sleep"`
	Jitter   uint32 `json:"jitter"`
	Caps     uint32 `json:"caps"`
	Hostname string `json:"hostname"`
	Username string `json:"username"`
}
