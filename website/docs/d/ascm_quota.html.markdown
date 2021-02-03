---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_quota"
sidebar_current: "docs-apsarastack-datasource-ascm-quota"
description: |-
    Provides a list of quota to the user.
---

# apsarastack\_ascm_quota

This data source provides the quota of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_quota" "default" {
  quota_type = "organization"
  quota_type_id = "54437"
  product_name = "SLB"
  output_file = "quota"
}
output "quota" {
  value = data.apsarastack_ascm_quota.default.*
}
```

## Argument Reference

The following arguments are supported:

  * `product_name` - (Required) The name of the service. Valid values: ECS, OSS, VPC, RDS, SLB, ODPS, and EIP.
  * `quota_type` - (Required) The type of the quota. Valid values: organization and resourceGroup.
  * `quota_type_id` - (Required) The ID of the quota type. Specify an organization ID when the QuotaType parameter is set to organization. Specify a resource set ID when the QuotaType parameter is set to resourceGroup.
  * `cluster_name` - (Optional) The name of the cluster. This reserved parameter is optional and can be left empty.
  * `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `id` - ID of the quota.
  * `quota_type` - Name of an organization, or a Resource Group.
  * `quota_type_id` - ID of an organization, or a Resource Group.
  * `total_vip_internal` - Total vip internal.
  * `total_vip_public` - Total vip public.
  * `region` - name of the region where product belong.
 
