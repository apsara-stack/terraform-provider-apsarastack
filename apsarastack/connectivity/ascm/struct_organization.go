package ascm

type Organization struct {
	Id               string `json:"Id" xml:"Id"`
	RegionId         string `json:"RegionId" xml:"RegionId"`
	Status           string `json:"Status" xml:"Status"`
	Name             string `json:"Name" xml:"Name"`
	PersonNum        string `json:"PersonNum" name:"PersonNum"`
	ResourceGroupNum string `json:"ResourceGroupNum" name:"ResourceGroupNum"`
	CreationTime     string `json:"CreationTime" xml:"CreationTime"`
	ResourceGroupId  string `json:"ResourceGroupId" xml:"ResourceGroupId"`
}
