package ascm

type ClustersByProduct struct {
	Body struct {
		ClusterList []struct {
			Region []string `json:"cn-neimeng-env30-d01"`
		} `json:"ClusterList"`
	} `json:"body"`
	Code            int  `json:"code"`
	SuccessResponse bool `json:"successResponse"`
}
