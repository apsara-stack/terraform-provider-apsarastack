---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_ascm_service_clusters"
sidebar_current: "docs-apsarastack-datasource-ascm-service-clusters"
description: |-
    Provides a list of service cluster to the user.
---

# apsarastack\_ascm_service_clusters

This data source provides the service clusters of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_service_cluster" "cluster" {
  output_file = "cluster"
  product_name = "slb"
}

output "cluster" {
  value = data.apsarastack_ascm_service_cluster.cluster.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance family IDs.
* `product_name` - (Required) Filter the results by specifying name of the service.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `cluster_list` - A list of instance families. Each element contains the following attributes:
    * `cluster_by_region` - cluster by a region.
