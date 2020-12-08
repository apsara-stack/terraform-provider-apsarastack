package apsarastack

type ResourceGroup struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		GmtCreated        int64  `json:"gmtCreated"`
		ID                int    `json:"id"`
		OrganizationID    int    `json:"organizationID"`
		OrganizationName  string `json:"organizationName"`
		ResourceGroupName string `json:"resourceGroupName"`
		RsID              string `json:"rsId"`
		Creator           string `json:"creator,omitempty"`
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
	PureListData bool `json:"pureListData"`
	Redirect     bool `json:"redirect"`
	Success      bool `json:"success"`
}

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

type PasswordPolicy struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data struct {
		ID                           string `json:"id"`
		HardExpiry                   bool   `json:"hardExpiry"`
		MaxLoginAttemps              int    `json:"maxLoginAttemps"`
		MaxPasswordAge               int    `json:"maxPasswordAge"`
		MinimumPasswordLength        int    `json:"minimumPasswordLength"`
		PasswordErrorCaptchaTime     int    `json:"passwordErrorCaptchaTime"`
		PasswordErrorLockPeriod      int    `json:"passwordErrorLockPeriod"`
		PasswordErrorTolerancePeriod int    `json:"passwordErrorTolerancePeriod"`
		PasswordReusePrevention      int    `json:"passwordReusePrevention"`
		RequireLowercaseCharacters   bool   `json:"requireLowercaseCharacters"`
		RequireNumbers               bool   `json:"requireNumbers"`
		RequireSymbols               bool   `json:"requireSymbols"`
		RequireUppercaseCharacters   bool   `json:"requireUppercaseCharacters"`
	} `json:"data"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

type Organization struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		Alias             string        `json:"alias"`
		CuserID           string        `json:"cuserId"`
		ID                int           `json:"id"`
		Internal          bool          `json:"internal"`
		Level             string        `json:"level"`
		MultiCloudStatus  string        `json:"multiCloudStatus"`
		MuserID           string        `json:"muserId"`
		Name              string        `json:"name"`
		ParentID          int           `json:"parentId"`
		SupportRegionList []interface{} `json:"supportRegionList"`
		UUID              string        `json:"uuid"`
		SupportRegions    string        `json:"supportRegions,omitempty"`
		Mtime             int64         `json:"mtime,omitempty"`
	} `json:"data"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	RequestID    string `json:"requestId"`
	Success      bool   `json:"success"`
}

type LoginPolicy struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
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
		UserCount   int    `json:"userCount"`
		Description string `json:"description,omitempty"`
		Mtime       int64  `json:"mtime,omitempty"`
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

type Roles struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		Active                 bool   `json:"active"`
		ArID                   string `json:"arId"`
		Code                   string `json:"code"`
		Default                bool   `json:"default"`
		Description            string `json:"description,omitempty"`
		Enable                 bool   `json:"enable"`
		ID                     int    `json:"id"`
		OrganizationVisibility string `json:"organizationVisibility"`
		OwnerOrganizationID    int    `json:"ownerOrganizationId"`
		RAMRole                bool   `json:"rAMRole"`
		RoleLevel              int64  `json:"roleLevel"`
		RoleName               string `json:"roleName"`
		RoleRange              string `json:"roleRange"`
		RoleType               string `json:"roleType"`
		UserCount              int    `json:"userCount"`
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
type InstanceFamily struct {
	Success bool `json:"success"`
	Data    []struct {
		ID          interface{} `json:"id"`
		GmtCreate   string      `json:"gmtCreate"`
		GmtModified string      `json:"gmtModified"`
		Creator     string      `json:"creator"`
		Modifier    string      `json:"modifier"`
		IsDeleted   string      `json:"isDeleted"`
		PageStart   int         `json:"pageStart"`
		PageSize    int         `json:"pageSize"`
		PageSort    string      `json:"pageSort"`
		PageOrder   string      `json:"pageOrder"`
		OrderBy     struct {
			ID string `json:"id"`
		} `json:"orderBy"`
		RegionID        interface{} `json:"regionId"`
		SpecFrom        interface{} `json:"specFrom"`
		BaseVersion     interface{} `json:"baseVersion"`
		Status          interface{} `json:"status"`
		SeriesID        string      `json:"seriesId"`
		SeriesName      string      `json:"seriesName"`
		SeriesNameLabel string      `json:"seriesNameLabel"`
		ResourceType    string      `json:"resourceType"`
		Deleted         bool        `json:"deleted"`
	} `json:"data"`
	Message   string      `json:"message"`
	Code      int         `json:"code"`
	HTTPCode  interface{} `json:"httpCode"`
	IP        interface{} `json:"ip"`
	RequestID interface{} `json:"requestId"`
	HTTPOk    bool        `json:"httpOk"`
}
type EnvironmentProduct struct {
	Code    int      `json:"code"`
	Result  []string `json:"result"`
	Success bool     `json:"success"`
}

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
type ClustersByProduct struct {
	Body struct {
		ClusterList []struct {
			Region []string `json:"cn-neimeng-env30-d01"`
		} `json:"ClusterList"`
	} `json:"body"`
	Code            int  `json:"code"`
	SuccessResponse bool `json:"successResponse"`
}
