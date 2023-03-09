package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackDnsRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackDnsRecord,
				Check:  resource.ComposeTestCheckFunc(

				//testAccCheckApsaraStackDataSourceID("data.apsarastack_dns_records.default"),
				//resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.record_id"),
				//resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.domain_id"),
				//resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.host_record"),
				//resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.type"),
				//resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.rr_set"),
				//resource.TestCheckNoResourceAttr("data.apsarastack_dns_records.default", "records.ttl"),
				),
			},
		},
	})
}

const dataSourceApsaraStackDnsRecord = `

resource "apsarastack_dns_domain" "default" {
 domain_name = "testdummy."
 remark = "test_dummy_1"
}
resource "apsarastack_dns_record" "default" {
 domain_id   = apsarastack_dns_domain.default.id
lba_strategy = "ALL_RR",
 name = "testrecord"
 type        = "A"
 ttl         = 300
 rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

data "apsarastack_dns_records" "default"{
 zone_id         = apsarastack_dns_record.default.zone_id
 name = apsarastack_dns_record.default.name
}
`
