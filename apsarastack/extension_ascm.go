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
type AddRoleList struct {
	LoginName  string   `json:"loginName"`
	RoleIDList []string `json:"roleIdList"`
}
type AscmUser struct {
	DisplayName      string   `json:"displayName"`
	Email            string   `json:"email"`
	LoginPolicyID    int      `json:"loginPolicyId"`
	MobileNationCode string   `json:"mobileNationCode"`
	PolicyID         int      `json:"policyId"`
	OrganizationID   string   `json:"organizationId"`
	LoginName        string   `json:"loginName"`
	FullName         string   `json:"fullName"`
	RoleIDList       []string `json:"roleIdList"`
	CellphoneNum     string   `json:"cellphoneNum"`
	UserEmail        string   `json:"userEmail"`
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
		Deleted            bool     `json:"deleted"`
		DisplayName        string   `json:"displayName"`
		Email              string   `json:"email"`
		EnableDingTalk     bool     `json:"enableDingTalk"`
		EnableEmail        bool     `json:"enableEmail"`
		EnableShortMessage bool     `json:"enableShortMessage"`
		ID                 int      `json:"id"`
		RoleIDList         []string `json:"roleIdList"`
		LastLoginTime      int64    `json:"lastLoginTime"`
		LoginName          string   `json:"loginName"`
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

type DeletedUser struct {
	Redirect       bool   `json:"redirect"`
	AsapiSuccess   bool   `json:"asapiSuccess"`
	Code           string `json:"code"`
	Cost           int    `json:"cost"`
	AsapiRequestID string `json:"asapiRequestId"`
	Data           []struct {
		CellphoneNum       string `json:"cellphoneNum"`
		Default            bool   `json:"default"`
		Deleted            bool   `json:"deleted"`
		DisplayName        string `json:"displayName"`
		EnableDingTalk     bool   `json:"enableDingTalk"`
		LoginName          string `json:"loginName"`
		ID                 int    `json:"id"`
		MobileNationCode   string `json:"mobileNationCode"`
		EnableEmail        bool   `json:"enableEmail"`
		Email              string `json:"email"`
		EnableShortMessage bool   `json:"enableShortMessage"`
		Status             string `json:"status"`
	} `json:"data"`
	Success  bool `json:"success"`
	PageInfo struct {
		Total       int `json:"total"`
		TotalPage   int `json:"totalPage"`
		PageSize    int `json:"pageSize"`
		CurrentPage int `json:"currentPage"`
	} `json:"pageInfo"`
	PureListData bool   `json:"pureListData"`
	Message      string `json:"message"`
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

type RamRole struct {
	Redirect       bool   `json:"redirect"`
	AsapiSuccess   bool   `json:"asapiSuccess"`
	Code           string `json:"code"`
	Cost           int    `json:"cost"`
	AsapiRequestID string `json:"asapiRequestId"`
	Data           []struct {
		Product                  string `json:"product"`
		AssumeRolePolicyDocument string `json:"assumeRolePolicyDocument"`
		OrganizationName         string `json:"organizationName"`
		RoleID                   string `json:"roleId"`
		Description              string `json:"description"`
		RoleType                 string `json:"roleType"`
		AliyunUserID             int    `json:"aliyunUserId"`
		OrganizationID           int    `json:"organizationId"`
		RoleName                 string `json:"roleName"`
		Ctime                    int64  `json:"ctime"`
		ID                       int    `json:"id"`
		Arn                      string `json:"arn"`
		Region                   string `json:"region"`
		CuserID                  string `json:"cuserId"`
	} `json:"data"`
	Success  bool `json:"success"`
	PageInfo struct {
		Total       int `json:"total"`
		TotalPage   int `json:"totalPage"`
		PageSize    int `json:"pageSize"`
		CurrentPage int `json:"currentPage"`
	} `json:"pageInfo"`
	PureListData bool   `json:"pureListData"`
	Message      string `json:"message"`
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
		RoleID                 int    `json:"roleId"`
		NewRoleName            string `json:"newRoleName"`
		NewDescription         string `json:"newDescription"`
		RoleName               string `json:"roleName,newRoleName"`
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

type LoginPolicy struct {
	Redirect       bool   `json:"redirect"`
	AsapiSuccess   bool   `json:"asapiSuccess"`
	Code           string `json:"code"`
	Cost           int    `json:"cost"`
	AsapiRequestID string `json:"asapiRequestId"`
	Data           []struct {
		MuserID    string `json:"muserId"`
		TimeRanges []struct {
			LoginPolicyID int    `json:"loginPolicyId"`
			StartTime     string `json:"startTime"`
			EndTime       string `json:"endTime"`
		} `json:"timeRanges"`
		Description string `json:"description"`
		Rule        string `json:"rule"`
		LpID        string `json:"lpId"`
		IPRanges    []struct {
			Protocol      string `json:"protocol"`
			IPRange       string `json:"ipRange"`
			LoginPolicyID int    `json:"loginPolicyId"`
		} `json:"ipRanges"`
		Default                bool   `json:"default"`
		UserCount              int    `json:"userCount"`
		OwnerOrganizationID    int    `json:"ownerOrganizationId"`
		Enable                 bool   `json:"enable"`
		Name                   string `json:"name"`
		ID                     int    `json:"id"`
		CuserID                string `json:"cuserId"`
		OrganizationVisibility string `json:"organizationVisibility"`
	} `json:"data"`
	Success  bool `json:"success"`
	PageInfo struct {
		Total       int `json:"total"`
		TotalPage   int `json:"totalPage"`
		PageSize    int `json:"pageSize"`
		CurrentPage int `json:"currentPage"`
	} `json:"pageInfo"`
	PureListData bool   `json:"pureListData"`
	Message      string `json:"message"`
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
	AsapiSuccess   bool   `json:"asapiSuccess"`
	Code           int    `json:"code"`
	AsapiRequestID string `json:"asapiRequestId"`
	Data           []struct {
		GmtModified string `json:"gmtModified"`
		Creator     string `json:"creator"`
		SeriesName  string `json:"seriesName"`
		Modifier    string `json:"modifier"`
		PageSize    int    `json:"pageSize"`
		OrderBy     struct {
			ID string `json:"id"`
		} `json:"orderBy"`
		GmtCreate       string `json:"gmtCreate"`
		SeriesID        string `json:"seriesId"`
		PageOrder       string `json:"pageOrder"`
		Deleted         bool   `json:"deleted"`
		IsDeleted       string `json:"isDeleted"`
		PageSort        string `json:"pageSort"`
		PageStart       int    `json:"pageStart"`
		SeriesNameLabel string `json:"seriesNameLabel"`
		ResourceType    string `json:"resourceType"`
	} `json:"data"`
	HTTPOk  bool   `json:"httpOk"`
	Success bool   `json:"success"`
	Message string `json:"message"`
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

type AscmQuota struct {
	Redirect        bool   `json:"redirect"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	Code            string `json:"code"`
	Cost            int    `json:"cost"`
	AsapiRequestID  string `json:"asapiRequestId"`
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	Data            struct {
		QuotaTypeID                 int    `json:"quotaTypeId"`
		QuotaBody                   string `json:"quotaBody"`
		QuotaType                   string `json:"quotaType"`
		RegionID                    string `json:"regionId"`
		ProductName                 string `json:"productName"`
		RegionName                  string `json:"regionName"`
		AllocateDiskCloudSsd        int    `json:"allocateDisk_cloud_ssd"`
		Cluster                     string `json:"cluster"`
		TotalMem                    int    `json:"totalMem"`
		TotalDisk                   int    `json:"totalDisk"`
		TotalDiskCloudEfficiency    int    `json:"totalDisk_cloud_efficiency"`
		AllocateGpu                 int    `json:"allocateGpu"`
		TargetType                  string `json:"targetType"`
		UsedMem                     int    `json:"usedMem"`
		AllocateCPU                 int    `json:"allocateCpu"`
		UsedDiskCloudEfficiency     int    `json:"usedDisk_cloud_efficiency"`
		TotalDiskCloudSsd           int    `json:"totalDisk_cloud_ssd"`
		UsedDiskCloudSsd            int    `json:"usedDisk_cloud_ssd"`
		TotalCPU                    int    `json:"totalCpu"`
		TotalCU                     int    `json:"totalCu"`
		DtFlag                      bool   `json:"dtFlag"`
		Ctime                       int64  `json:"ctime"`
		ID                          int    `json:"id"`
		UsedGpu                     int    `json:"usedGpu"`
		Region                      string `json:"region"`
		AllocateMem                 int    `json:"allocateMem"`
		AllocateDiskCloudEfficiency int    `json:"allocateDisk_cloud_efficiency"`
		TotalGpu                    int    `json:"totalGpu"`
		UsedCPU                     int    `json:"usedCpu"`
		TotalVPC                    int    `json:"totalVPC"`
		TotalEIP                    int    `json:"totalEIP"`
		UsedVPC                     int    `json:"usedVPC"`
		AllocateVPC                 int    `json:"allocateVPC"`
		TotalAmount                 int    `json:"totalAmount"`
		UsedVipInternal             int    `json:"usedVipInternal"`
		AllocateVipPublic           int    `json:"allocateVipPublic"`
		TotalVipPublic              int    `json:"totalVipPublic"`
		TotalVipInternal            int    `json:"totalVipInternal"`
		UsedVipPublic               int    `json:"usedVipPublic"`
		AllocateVipInternal         int    `json:"allocateVipInternal"`
		UsedAmount                  int    `json:"usedAmount"`
		AllocateAmount              int    `json:"allocateAmount"`
		UsedDisk                    int    `json:"usedDisk"`
		AllocateDisk                int    `json:"allocateDisk"`
	} `json:"data"`
	RequestID      string `json:"requestId"`
	Success        bool   `json:"success"`
	PureListData   bool   `json:"pureListData"`
	Message        string `json:"message"`
	ServerRole     string `json:"serverRole"`
	Domain         string `json:"domain"`
	API            string `json:"api"`
	AsapiErrorCode string `json:"asapiErrorCode"`
}
