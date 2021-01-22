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
type MetaList struct {
	TotalCount int    `json:"TotalCount"`
	RequestID  string `json:"RequestId"`
	Resources  struct {
		Resource []struct {
			MetricName  string `json:"MetricName"`
			Periods     string `json:"Periods"`
			Description string `json:"Description"`
			Dimensions  string `json:"Dimensions"`
			Labels      string `json:"Labels"`
			Unit        string `json:"Unit"`
			Statistics  string `json:"Statistics"`
			Namespace   string `json:"Namespace"`
		} `json:"Resource"`
	} `json:"Resources"`
	Code    int  `json:"Code"`
	Success bool `json:"Success"`
}
