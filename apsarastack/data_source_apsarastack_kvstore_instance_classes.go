package apsarastack

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceApsaraStackKVStoreInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackKVStoreAvailableResourceRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(KVStoreMemcache),
					string(KVStoreRedis),
				}, false),
				Default: string(KVStoreRedis),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "cluster", "rwsplit"}, false),
			},
			"edition_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Community", "Enterprise"}, false),
			},
			"series_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"enhanced_performance_type", "hybrid_storage"}, false),
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"double", "single", "readone", "readthree", "readfive"}, false),
			},
			"shard_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 8, 16, 32, 64, 128, 256}),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PrePaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Price"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackKVStoreAvailableResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := r_kvstore.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.ZoneId = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	request.InstanceChargeType = instanceChargeType
	request.Engine = d.Get("engine").(string)
	var response = &r_kvstore.DescribeAvailableResourceResponse{}
	err := resource.Retry(time.Minute*5, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAvailableResource(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*r_kvstore.DescribeAvailableResourceResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_kvstore_instance_classes", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	var instanceClasses []string
	var ids []string

	engine, engineGot := d.GetOk("engine")
	engine = strings.ToLower(engine.(string))
	engineVersion, engineVersionGot := d.GetOk("engine_version")
	architecture, architectureGot := d.GetOk("architecture")
	editionType, editionTypeGot := d.GetOk("edition_type")
	seriesType, seriesTypeGot := d.GetOk("series_type")
	shardNumber, shardNumberGot := d.GetOk("shard_number")
	nodeType, nodeTypeGot := d.GetOk("node_type")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		zondId := AvailableZone.ZoneId
		ids = append(ids, zondId)
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEditionType := range SupportedEngine.SupportedEditionTypes.SupportedEditionType {
				if editionTypeGot && editionType.(string) != SupportedEditionType.EditionType {
					continue
				}
				ids = append(ids, SupportedEditionType.EditionType)
				for _, SupportedSeriesType := range SupportedEditionType.SupportedSeriesTypes.SupportedSeriesType {
					if seriesTypeGot && seriesType.(string) != SupportedSeriesType.SeriesType {
						continue
					}
					for _, SupportedEngineVersion := range SupportedSeriesType.SupportedEngineVersions.SupportedEngineVersion {
						if engineVersionGot && engineVersion.(string) != SupportedEngineVersion.Version {
							continue
						}
						for _, SupportedArchitectureType := range SupportedEngineVersion.SupportedArchitectureTypes.SupportedArchitectureType {
							if architectureGot && architecture.(string) != SupportedArchitectureType.Architecture {
								continue
							}
							for _, SupportedShardNumber := range SupportedArchitectureType.SupportedShardNumbers.SupportedShardNumber {
								number, _ := strconv.Atoi(SupportedShardNumber.ShardNumber)
								if shardNumberGot && shardNumber.(int) != number {
									continue
								}
								for _, SupportedNodeType := range SupportedShardNumber.SupportedNodeTypes.SupportedNodeType {
									if nodeTypeGot && nodeType.(string) != SupportedNodeType.SupportedNodeType {
										continue
									}
									for _, AvailableResource := range SupportedNodeType.AvailableResources.AvailableResource {
										instanceClasses = append(instanceClasses, AvailableResource.InstanceClass)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))

	var instanceClassPrices []map[string]interface{}
	sortedBy := d.Get("sorted_by").(string)
	if sortedBy == "Price" && len(instanceClasses) > 0 {
		bssopenapiService := BssopenapiService{client}
		priceList, err := getKVStoreInstanceClassPrice(bssopenapiService, instanceChargeType, instanceClasses)
		if err != nil {
			return WrapError(err)
		}
		for i, instanceClass := range instanceClasses {
			classPrice := map[string]interface{}{
				"instance_class": instanceClass,
				"price":          fmt.Sprintf("%.4f", priceList[i]),
			}
			instanceClassPrices = append(instanceClassPrices, classPrice)
		}
		sort.SliceStable(instanceClassPrices, func(i, j int) bool {
			iPrice, _ := strconv.ParseFloat(instanceClassPrices[i]["price"].(string), 64)
			jPrice, _ := strconv.ParseFloat(instanceClassPrices[j]["price"].(string), 64)
			return iPrice < jPrice
		})

		err = d.Set("classes", instanceClassPrices)
		if err != nil {
			return WrapError(err)
		}

		instanceClasses = instanceClasses[:0]
		for _, instanceClass := range instanceClassPrices {
			instanceClasses = append(instanceClasses, instanceClass["instance_class"].(string))
		}
	}

	err = d.Set("instance_classes", instanceClasses)
	if err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), instanceClassPrices)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func getKVStoreInstanceClassPrice(bssopenapiService BssopenapiService, instanceChargeType string, instanceClasses []string) ([]float64, error) {
	client := bssopenapiService.client
	var modules interface{}
	moduleCode := "InstanceClass"
	var payAsYouGo []bssopenapi.GetPayAsYouGoPriceModuleList
	var subsciption []bssopenapi.GetSubscriptionPriceModuleList
	for _, instanceClass := range instanceClasses {
		config := fmt.Sprintf("InstanceClass:%s,Region:%s", instanceClass, client.Region)
		if instanceChargeType == string(PostPaid) {
			payAsYouGo = append(payAsYouGo, bssopenapi.GetPayAsYouGoPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
				PriceType:  "Hour",
			})
		} else {
			subsciption = append(subsciption, bssopenapi.GetSubscriptionPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
			})

		}
	}

	if len(payAsYouGo) != 0 {
		modules = payAsYouGo
	} else {
		modules = subsciption
	}

	return bssopenapiService.GetInstanceTypePrice("redisa", "", modules)
}
