package ascm

type SpecificField struct {
	Success   bool        `json:"success"`
	Data      []string    `json:"data"`
	Message   string      `json:"message"`
	Code      int         `json:"code"`
	HTTPCode  interface{} `json:"httpCode"`
	IP        interface{} `json:"ip"`
	RequestID interface{} `json:"requestId"`
	HTTPOk    bool        `json:"httpOk"`
}
