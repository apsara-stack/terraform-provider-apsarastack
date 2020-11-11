package ascm

// EndpointMap Endpoint Data
var EndpointMap map[string]string

// EndpointType regional or central
var EndpointType = "regional"

// GetEndpointMap Get Endpoint Data Map
func GetEndpointMap() map[string]string {
	if EndpointMap == nil {
		EndpointMap = map[string]string{
			"cn-shanghai-internal-test-1": "ascm.aliyuncs.com",
			"cn-beijing-gov-1":            "ascm.aliyuncs.com",
			"cn-shenzhen-su18-b01":        "ascm.aliyuncs.com",
			"cn-beijing":                  "ascm.aliyuncs.com",
			"cn-shanghai-inner":           "ascm.aliyuncs.com",
			"cn-shenzhen-st4-d01":         "ascm.aliyuncs.com",
			"cn-haidian-cm12-c01":         "ascm.aliyuncs.com",
			"cn-hangzhou-internal-prod-1": "ascm.aliyuncs.com",
			"cn-north-2-gov-1":            "ascm.aliyuncs.com",
			"cn-yushanfang":               "ascm.aliyuncs.com",
			"cn-qingdao":                  "ascm.aliyuncs.com",
			"cn-hongkong-finance-pop":     "ascm.aliyuncs.com",
			"cn-qingdao-nebula":           "ascm-nebula.cn-qingdao-nebula.aliyuncs.com",
			"cn-shanghai":                 "ascm.aliyuncs.com",
			"cn-shanghai-finance-1":       "ascm.aliyuncs.com",
			"cn-hongkong":                 "ascm.aliyuncs.com",
			"cn-beijing-finance-pop":      "ascm.aliyuncs.com",
			"cn-wuhan":                    "ascm.aliyuncs.com",
			"us-west-1":                   "ascm.aliyuncs.com",
			"cn-shenzhen":                 "ascm.aliyuncs.com",
			"cn-zhengzhou-nebula-1":       "ascm-nebula.cn-qingdao-nebula.aliyuncs.com",
			"rus-west-1-pop":              "ascm.aliyuncs.com",
			"cn-shanghai-et15-b01":        "ascm.aliyuncs.com",
			"cn-hangzhou-bj-b01":          "ascm.aliyuncs.com",
			"cn-hangzhou-internal-test-1": "ascm.aliyuncs.com",
			"eu-west-1-oxs":               "ascm-nebula.cn-shenzhen-cloudstone.aliyuncs.com",
			"cn-zhangbei-na61-b01":        "ascm.aliyuncs.com",
			"cn-beijing-finance-1":        "ascm.aliyuncs.com",
			"cn-hangzhou-internal-test-3": "ascm.aliyuncs.com",
			"cn-shenzhen-finance-1":       "ascm.aliyuncs.com",
			"cn-hangzhou-internal-test-2": "ascm.aliyuncs.com",
			"cn-hangzhou-test-306":        "ascm.aliyuncs.com",
			"cn-shanghai-et2-b01":         "ascm.aliyuncs.com",
			"cn-hangzhou-finance":         "ascm.aliyuncs.com",
			"ap-southeast-1":              "ascm.aliyuncs.com",
			"cn-beijing-nu16-b01":         "ascm.aliyuncs.com",
			"cn-edge-1":                   "ascm-nebula.cn-qingdao-nebula.aliyuncs.com",
			"us-east-1":                   "ascm.aliyuncs.com",
			"cn-fujian":                   "ascm.aliyuncs.com",
			"ap-northeast-2-pop":          "ascm.aliyuncs.com",
			"cn-shenzhen-inner":           "ascm.aliyuncs.com",
			"cn-zhangjiakou-na62-a01":     "ascm.cn-zhangjiakou.aliyuncs.com",
			"cn-hangzhou":                 "ascm.aliyuncs.com",
		}
	}
	return EndpointMap
}

// GetEndpointType Get Endpoint Type Value
func GetEndpointType() string {
	return EndpointType
}
