---
subcategory: "EDAS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_edas_k8s_cluster"
sidebar_current: "docs-apsarastack-resource-edas-k8s-cluster"
description: |-
  Provides an EDAS K8s cluster resource.
---

# apsarastack\_edas\_k8s\_cluster

Provides an EDAS K8s cluster resource. For information about EDAS K8s Cluster and how to use it, see[What is EDAS K8s Cluster](https://www.alibabacloud.com/help/en/doc-detail/85108.htm).

-> **NOTE:** Available in 1.0.11+

## Example Usage

Basic Usage

```
resource "apsarastack_edas_k8s_cluster" "default" {
  cs_cluster_id = "xxxx-xxx-xxx"
}
```

## Argument Reference

The following arguments are supported:

* `cs_cluster_id` - (Required, ForceNew) The ID of the apsarastack container service kubernetes cluster that you want to import.
* `namespace_id` - (Optional, ForceNew) The ID of the namespace where you want to import. You can call the [ListUserDefineRegion](https://www.alibabacloud.com/help/en/doc-detail/149377.htm?spm=a2c63.p38356.879954.34.331054faK2yNvC#doc-api-Edas-ListUserDefineRegion) operation to query the namespace ID.


## Attributes Reference

The following attributes are exported:

* `cluster_name` - The name of the cluster that you want to create.
* `cluster_type` - The type of the cluster that you want to create. Valid values only: 5: K8s cluster. 
* `network_mode` - The network type of the cluster that you want to create. Valid values: 1: classic network. 2: VPC.
* `region_id` - The ID of the region.
* `vpc_id` - The ID of the Virtual Private Cloud (VPC) for the cluster.
* `cluster_import_status` - The import status of cluster: 
    `1`: success.
    `2`: failed.
    `3`: importing. 
    `4`: deleted.

## Import

EDAS cluster can be imported using the id, e.g.

```
$ terraform import apsarastack_edas_k8s_cluster.cluster cluster_id
```
