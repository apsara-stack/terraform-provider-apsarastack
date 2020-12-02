package ascm

//type ResourceGroup struct {
//	OrganizationId                  			string              `json:"OrganizationId" xml:"OrganizationId"`
//	RegionId               			string              `json:"RegionId" xml:"RegionId"`
//	ExpiredTime                   string                        `json:"ExpiredTime" xml:"ExpiredTime"`
//	Status                 			string              `json:"Status" xml:"Status"`
//	ResourceGroupName                			string              `json:"Name" xml:"ResourceGroupName"`
//	CreationTime           			string               `json:"CreationTime" xml:"CreationTime"`
//	ResourceGroupId        			string               `json:"ResourceGroupId" xml:"ResourceGroupId"`
//}

type ResourceGroup struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		Status            string `json:"Status"`
		Creator           string `json:"creator"`
		GmtCreated        int64  `json:"gmtCreated"`
		ResourceGroupID   int    `json:"id"`
		OrganizationID    int    `json:"organizationID"`
		OrganizationName  string `json:"organizationName"`
		ResourceGroupName string `json:"resourceGroupName"`
		RsID              string `json:"rsId"`
		GmtModified       int64  `json:"gmtModified,omitempty"`
		ResourceGroupType int    `json:"resourceGroupType,omitempty"`
	} `json:"data"`
	Message  string `json:"message"`
	PageInfo struct {
		CurrentPage int   `json:"currentPage"`
		PageSize    int64 `json:"pageSize"`
		Total       int   `json:"total"`
		TotalPage   int   `json:"totalPage"`
	} `json:"pageInfo"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
	Status       string `json:"Status"`
}
