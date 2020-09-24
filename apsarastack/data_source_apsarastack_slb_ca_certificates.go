package apsarastack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceApsaraStackSlbCACertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackSlbCACertificatesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"tags": tagsSchema(),
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackSlbCACertificatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := slb.CreateDescribeCACertificatesRequest()
	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		Tags := make([]slb.DescribeCACertificatesTag, 0, len(tags))
		for k, v := range tags {
			certificatesTag := slb.DescribeCACertificatesTag{
				Key:   k,
				Value: v.(string),
			}
			Tags = append(Tags, certificatesTag)
		}
		request.Tag = &Tags
	}
	request.RegionId = client.RegionId
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeCACertificates(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_slb_ca_certificates", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeCACertificatesResponse)
	var filteredTemp []slb.CACertificate
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, certificate := range response.CACertificates.CACertificate {
			if r != nil && !r.MatchString(certificate.CACertificateName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[certificate.CACertificateId]; !ok {
					continue
				}
			}

			filteredTemp = append(filteredTemp, certificate)
		}
	} else {
		filteredTemp = response.CACertificates.CACertificate
	}

	return slbCACertificatesDescriptionAttributes(d, filteredTemp, meta)
}

func caCertificateTagsMappings(d *schema.ResourceData, id string, meta interface{}) map[string]string {
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := SlbService{client}
	tags, err := slbService.DescribeTags(id, nil, TagResourceCertificate)

	if err != nil {
		return nil
	}

	return slbTagsToMap(tags)
}

func slbCACertificatesDescriptionAttributes(d *schema.ResourceData, certificates []slb.CACertificate, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, certificate := range certificates {

		mapping := map[string]interface{}{
			"id":                certificate.CACertificateId,
			"name":              certificate.CACertificateName,
			"fingerprint":       certificate.Fingerprint,
			"created_time":      certificate.CreateTime,
			"created_timestamp": certificate.CreateTimeStamp,
			"region_id":         certificate.RegionId,
			"tags":              caCertificateTagsMappings(d, certificate.CACertificateId, meta),
		}
		ids = append(ids, certificate.CACertificateId)
		names = append(names, certificate.CACertificateName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("certificates", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
