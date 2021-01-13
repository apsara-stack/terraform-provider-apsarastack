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
