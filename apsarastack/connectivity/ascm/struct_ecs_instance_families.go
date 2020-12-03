package ascm

type EcsInstanceFamily struct {
	Success bool `json:"success"`
	Data    struct {
		InstanceTypeFamilies []struct {
			InstanceTypeFamilyID string `json:"instanceTypeFamilyId"`
			Generation           string `json:"generation"`
		} `json:"instanceTypeFamilies"`
	} `json:"data"`
	Message   string      `json:"message"`
	Code      int         `json:"code"`
	HTTPCode  interface{} `json:"httpCode"`
	IP        interface{} `json:"ip"`
	RequestID interface{} `json:"requestId"`
	HTTPOk    bool        `json:"httpOk"`
}
