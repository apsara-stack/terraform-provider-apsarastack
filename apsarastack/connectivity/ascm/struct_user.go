package ascm

type User struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		AccessKeys []struct {
			AccesskeyID string `json:"accesskeyId"`
			Ctime       int64  `json:"ctime"`
			CuserID     string `json:"cuserId"`
			ID          int    `json:"id"`
			Region      string `json:"region"`
			Status      string `json:"status"`
		} `json:"accessKeys"`
		CellphoneNum string `json:"cellphoneNum"`
		Default      bool   `json:"default"`
		DefaultRole  struct {
			Active                 bool   `json:"active"`
			ArID                   string `json:"arId"`
			Code                   string `json:"code"`
			CuserID                string `json:"cuserId"`
			Default                bool   `json:"default"`
			Description            string `json:"description"`
			Enable                 bool   `json:"enable"`
			ID                     int    `json:"id"`
			MuserID                string `json:"muserId"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			RAMRole                bool   `json:"rAMRole"`
			RoleLevel              int64  `json:"roleLevel"`
			RoleName               string `json:"roleName"`
			RoleRange              string `json:"roleRange"`
			RoleType               string `json:"roleType"`
		} `json:"defaultRole"`
		Deleted            bool   `json:"deleted"`
		DisplayName        string `json:"displayName"`
		Email              string `json:"email"`
		EnableDingTalk     bool   `json:"enableDingTalk"`
		EnableEmail        bool   `json:"enableEmail"`
		EnableShortMessage bool   `json:"enableShortMessage"`
		ID                 int    `json:"id"`
		LastLoginTime      int64  `json:"lastLoginTime"`
		LoginName          string `json:"loginName"`
		LoginPolicy        struct {
			CuserID  string `json:"cuserId"`
			Default  bool   `json:"default"`
			Enable   bool   `json:"enable"`
			ID       int    `json:"id"`
			IPRanges []struct {
				IPRange       string `json:"ipRange"`
				LoginPolicyID int    `json:"loginPolicyId"`
				Protocol      string `json:"protocol"`
			} `json:"ipRanges"`
			LpID                   string `json:"lpId"`
			MuserID                string `json:"muserId"`
			Name                   string `json:"name"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			Rule                   string `json:"rule"`
			TimeRanges             []struct {
				EndTime       string `json:"endTime"`
				LoginPolicyID int    `json:"loginPolicyId"`
				StartTime     string `json:"startTime"`
			} `json:"timeRanges"`
		} `json:"loginPolicy"`
		MobileNationCode string `json:"mobileNationCode"`
		Organization     struct {
			Alias             string        `json:"alias"`
			Ctime             int64         `json:"ctime"`
			CuserID           string        `json:"cuserId"`
			ID                int           `json:"id"`
			Internal          bool          `json:"internal"`
			Level             string        `json:"level"`
			Mtime             int64         `json:"mtime"`
			MultiCloudStatus  string        `json:"multiCloudStatus"`
			MuserID           string        `json:"muserId"`
			Name              string        `json:"name"`
			ParentID          int           `json:"parentId"`
			SupportRegionList []interface{} `json:"supportRegionList"`
			UUID              string        `json:"uuid"`
		} `json:"organization,omitempty"`
		ParentPk   string `json:"parentPk"`
		PrimaryKey string `json:"primaryKey"`
		Roles      []struct {
			Active                 bool   `json:"active"`
			ArID                   string `json:"arId"`
			Code                   string `json:"code"`
			Default                bool   `json:"default"`
			Description            string `json:"description"`
			Enable                 bool   `json:"enable"`
			ID                     int    `json:"id"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			RAMRole                bool   `json:"rAMRole"`
			RoleLevel              int64  `json:"roleLevel"`
			RoleName               string `json:"roleName"`
			RoleRange              string `json:"roleRange"`
			RoleType               string `json:"roleType"`
		} `json:"roles"`
		Status         string        `json:"status"`
		UserGroupRoles []interface{} `json:"userGroupRoles"`
		UserGroups     []interface{} `json:"userGroups"`
		UserRoles      []struct {
			Active                 bool   `json:"active"`
			ArID                   string `json:"arId"`
			Code                   string `json:"code"`
			CuserID                string `json:"cuserId"`
			Default                bool   `json:"default"`
			Description            string `json:"description"`
			Enable                 bool   `json:"enable"`
			ID                     int    `json:"id"`
			MuserID                string `json:"muserId"`
			OrganizationVisibility string `json:"organizationVisibility"`
			OwnerOrganizationID    int    `json:"ownerOrganizationId"`
			RAMRole                bool   `json:"rAMRole"`
			RoleLevel              int64  `json:"roleLevel"`
			RoleName               string `json:"roleName"`
			RoleRange              string `json:"roleRange"`
			RoleType               string `json:"roleType"`
		} `json:"userRoles"`
	} `json:"data"`
	Message  string `json:"message"`
	PageInfo struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		Total       int `json:"total"`
		TotalPage   int `json:"totalPage"`
	} `json:"pageInfo"`
	PureListData bool `json:"pureListData"`
	Redirect     bool `json:"redirect"`
	Success      bool `json:"success"`
}
