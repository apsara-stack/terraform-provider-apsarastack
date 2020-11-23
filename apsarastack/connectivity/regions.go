package connectivity

// Region represents ECS region
type Region string

const (
	Hangzhou    = Region("cn-hangzhou")
	Qingdao     = Region("cn-qingdao")
	Beijing     = Region("cn-beijing")
	Hongkong    = Region("cn-hongkong")
	Shenzhen    = Region("cn-shenzhen")
	Shanghai    = Region("cn-shanghai")
	Zhangjiakou = Region("cn-zhangjiakou")
	Huhehaote   = Region("cn-huhehaote")
	ChengDu     = Region("cn-chengdu")
	HeYuan      = Region("cn-heyuan")
	WuLanChaBu  = Region("cn-wulanchabu")

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")

	APSouth1 = Region("ap-south-1")

	USWest1 = Region("us-west-1")
	USEast1 = Region("us-east-1")

	MEEast1 = Region("me-east-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	ShenZhenFinance = Region("cn-shenzhen-finance-1")
	ShanghaiFinance = Region("cn-shanghai-finance-1")
)

var ValidRegions = []Region{
	Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, WuLanChaBu,
	USWest1, USEast1,
	APNorthEast1, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5,
	APSouth1,
	MEEast1,
	EUCentral1, EUWest1,
}

var EcsClassicSupportedRegions = []Region{Shenzhen, Shanghai, Beijing, Qingdao, Hangzhou, Hongkong, USWest1, APSouthEast1}
var EcsSpotNoSupportedRegions = []Region{APSouth1}
var RdsClassicNoSupportedRegions = []Region{APSouth1, APSouthEast2, APSouthEast3, APNorthEast1, EUCentral1, EUWest1, MEEast1}
var RdsMultiAzNoSupportedRegions = []Region{Qingdao, APNorthEast1, APSouthEast5, MEEast1}
var RdsPPASNoSupportedRegions = []Region{Qingdao, USEast1, APNorthEast1, EUCentral1, MEEast1, APSouthEast2, APSouthEast3, APSouth1, APSouthEast5, ChengDu, EUWest1}
var RouteTableNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var SlbClassicNoSupportedRegions = []Region{APNorthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, MEEast1, EUCentral1, EUWest1, Huhehaote, Zhangjiakou}
var HttpHttpsHealthCheckMehtodSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, EUWest1, ChengDu, Qingdao, Hongkong, Shenzhen, APSouthEast5, Zhangjiakou, Huhehaote, MEEast1, APSouth1, EUCentral1, USWest1, APSouthEast3, APSouthEast2, APSouthEast1, APNorthEast1}
var OssVersioningSupportedRegions = []Region{APSouth1}
var OssSseSupportedRegions = []Region{Qingdao, Hangzhou, Beijing, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouth1, USEast1}
var CRNoSupportedRegions = []Region{Beijing, Hangzhou, Qingdao, Huhehaote, Zhangjiakou}
var KmsSkippedRegions = []Region{}
var GpdbClassicNoSupportedRegions = []Region{APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1, EUCentral1}
var MongoDBClassicNoSupportedRegions = []Region{Huhehaote, Zhangjiakou, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1}
var MongoDBMultiAzSupportedRegions = []Region{Hangzhou, Beijing, Shenzhen, EUCentral1}
