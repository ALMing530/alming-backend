package entity

type SysInfo struct {
	Cpu  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
	Time string  `json:"time"`
}
