package ascm

type RegionsByProduct struct {
	Body struct {
		RegionList []struct {
			RegionID   string `json:"RegionId"`
			RegionType string `json:"RegionType"`
		} `json:"RegionList"`
	} `json:"body"`
	Code            int  `json:"code"`
	SuccessResponse bool `json:"successResponse"`
}
