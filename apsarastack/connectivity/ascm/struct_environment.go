package ascm

type EnvironmentProduct struct {
	Code    int      `json:"code"`
	Result  []string `json:"result"`
	Success bool     `json:"success"`
}
