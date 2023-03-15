---
subcategory: "DNS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_dns_records"
sidebar_current: "docs-apsarastack-datasource-dns-records"
description: |-
    Provides a list of records available to the dns.
---

# apsarastack\_dns\_records

This data source provides a list of DNS Domain Records in an ApsaraStack Cloud account according to the specified filters.

## Example Usage

```
resource "apsarastack_dns_domain" "default" {
  domain_name = "domaintest."
  remark = "testing Domain"
}

# Create a new Domain record
resource "apsarastack_dns_record" "default" {
   zone_id   = apsarastack_dns_domain.default.domain_id
  name = "testing_record"
  type        = "A"
  remark = "testing Record"
  ttl         = 300
  lba_strategy = "ALL_RR"
  rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

data "apsarastack_dns_records" "default"{
 zone_id         = alibabacloudstack_dns_record.default.zone_id
 name = alibabacloudstack_dns_record.default.name
}
output "records" {
  value = data.apsarastack_dns_records.default.*
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The domain Id associated to the records.
* `name` - (Optional) Host record regex. 
* `value_regex` - (Optional) Host record value regex. 
* `type` - (Optional) Record type. Valid items are `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `REDIRECT_URL`, `FORWORD_URL` .
* `ids` - (Optional) A list of record IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of record IDs. 
* `records` - A list of records. Each element contains the following attributes:
  * `record_id` - ID of the record.
  * `zone_id` - ID of the domain the record belongs to.
  * `name` - Host record of the domain.
  * `type` - Type of the record.
  * `ttl` - TTL of the record.
  * `remark` - Description of the record.
  * `rr_set` - RrSet for the record.
