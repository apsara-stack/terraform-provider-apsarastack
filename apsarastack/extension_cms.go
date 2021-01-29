package apsarastack

const (
	Average          = "Average"
	Minimum          = "Minimum"
	Maximum          = "Maximum"
	ErrorCodeMaximum = "ErrorCodeMaximum"
)

const (
	MoreThan        = ">"
	MoreThanOrEqual = ">="
	LessThan        = "<"
	LessThanOrEqual = "<="
	Equal           = "=="
	NotEqual        = "!="
)

const (
	SiteMonitorHTTP = "HTTP"
	SiteMonitorPing = "Ping"
	SiteMonitorTCP  = "TCP"
	SiteMonitorUDP  = "UDP"
	SiteMonitorDNS  = "DNS"
	SiteMonitorSMTP = "SMTP"
	SiteMonitorPOP3 = "POP3"
	SiteMonitorFTP  = "FTP"
)

type CmsContact struct {
	Code string `json:"Code"`
	Cost int    `json:"Cost"`
	Data []struct {
		Cid  string `json:"Cid"`
		Name string `json:"Name"`
	} `json:"Data"`
	Message  string `json:"Message"`
	Redirect bool   `json:"Redirect"`
	Success  bool   `json:"Success"`
}

type AlaramRul struct {
	Redirect       bool   `json:"redirect"`
	TotalCount     int    `json:"TotalCount"`
	AsapiSuccess   bool   `json:"asapiSuccess"`
	Code           string `json:"code"`
	Cost           int    `json:"cost"`
	AsapiRequestID string `json:"asapiRequestId"`
	PageSize       int    `json:"PageSize"`
	PageNumber     int    `json:"PageNumber"`
	Success        bool   `json:"success"`
	Alarms         struct {
		Alarm []struct {
			GroupName           string `json:"GroupName"`
			NoEffectiveInterval string `json:"NoEffectiveInterval"`
			SilenceTime         int    `json:"SilenceTime"`
			ContactGroups       string `json:"ContactGroups"`
			MailSubject         string `json:"MailSubject"`
			SourceType          string `json:"SourceType"`
			RuleID              string `json:"RuleId"`
			Period              int    `json:"Period"`
			Dimensions          string `json:"Dimensions"`
			EffectiveInterval   string `json:"EffectiveInterval"`
			AlertState          string `json:"AlertState"`
			Namespace           string `json:"Namespace"`
			GroupID             string `json:"GroupId"`
			MetricName          string `json:"MetricName"`
			Department          int    `json:"Department"`
			EnableState         bool   `json:"EnableState"`
			Escalations         struct {
				Critical struct {
					ComparisonOperator string `json:"ComparisonOperator"`
					Times              int    `json:"Times"`
					Statistics         string `json:"Statistics"`
					Threshold          string `json:"Threshold"`
				} `json:"Critical"`
				Info struct {
				} `json:"Info"`
				Warn struct {
				} `json:"Warn"`
			} `json:"Escalations"`
			DepartmentName    string `json:"DepartmentName"`
			Webhook           string `json:"Webhook"`
			Resources         string `json:"Resources"`
			RegionID          string `json:"RegionId"`
			RuleName          string `json:"RuleName"`
			ResourceGroup     int    `json:"ResourceGroup"`
			ResourceGroupName string `json:"ResourceGroupName"`
		} `json:"Alarm"`
	} `json:"Alarms"`
	PureListData bool   `json:"pureListData"`
	Message      string `json:"message"`
}
type AlarmRules []struct {
	Action string `json:"action"`
	Params struct {
		ResourceGroupID   interface{} `json:"resourceGroupId"`
		Namespace         string      `json:"Namespace"`
		GroupName         string      `json:"GroupName"`
		GroupID           string      `json:"GroupId"`
		SilenceTime       int         `json:"SilenceTime"`
		EffectiveInterval string      `json:"EffectiveInterval"`
		Webhook           string      `json:"Webhook"`
		ContactGroups     string      `json:"ContactGroups"`
		EmailSubject      string      `json:"EmailSubject"`
		Resources         []struct {
			BucketName string `json:"BucketName"`
		} `json:"Resources"`
		EscalationsCriticalThreshold          int    `json:"Escalations.Critical.Threshold"`
		EscalationsCriticalComparisonOperator string `json:"Escalations.Critical.ComparisonOperator"`
		EscalationsCriticalStatistics         string `json:"Escalations.Critical.Statistics"`
		EscalationsCriticalTimes              int    `json:"Escalations.Critical.Times"`
		EscalationsWarnThreshold              int    `json:"Escalations.Warn.Threshold"`
		EscalationsWarnComparisonOperator     string `json:"Escalations.Warn.ComparisonOperator"`
		EscalationsWarnStatistics             string `json:"Escalations.Warn.Statistics"`
		EscalationsWarnTimes                  int    `json:"Escalations.Warn.Times"`
		EscalationsInfoThreshold              int    `json:"Escalations.Info.Threshold"`
		EscalationsInfoComparisonOperator     string `json:"Escalations.Info.ComparisonOperator"`
		EscalationsInfoStatistics             string `json:"Escalations.Info.Statistics"`
		EscalationsInfoTimes                  int    `json:"Escalations.Info.Times"`
		RuleName                              string `json:"RuleName"`
		RuleID                                string `json:"RuleId"`
		MetricName                            string `json:"MetricName"`
		Period                                string `json:"Period"`
		Unit                                  string `json:"Unit"`
		InstanceID                            string `json:"instanceID"`
	} `json:"params"`
}
