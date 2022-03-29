---
subcategory: "Log Service (SLS)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_log_store"
sidebar_current: "docs-apsarastack-resource-log-store"
description: |-
  Provides a Apsarastack log store resource.
---

# apsarastack\_log\_store

The log store is a unit in Log Service to collect, store, and query the log data. Each log store belongs to a project,
and each project can create multiple Logstores. [Refer to details](https://help.aliyun.com/apsara/enterprise/v_3_16_0_20220117/sls/enterprise-ascm-developer-guide/CreateLogstore.html?spm=a2c4g.14484438.10001.307)

## Example Usage

Basic Usage

```
resource "apsarastack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "apsarastack_log_store" "example" {
  project               = apsarastack_log_project.example.name
  name                  = "tf-log-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
```
Encrypt Usage
```
resource "apsarastack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "apsarastack_log_store" "example" {
  project               = apsarastack_log_project.example.name
  name                  = "tf-log-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
  encryption            = true
}
```

## Module Support

You can use the existing [sls module](https://registry.terraform.io/modules/terraform-alicloud-modules/sls/alicloud) 
to create SLS project, store and store index one-click, like ECS instances.

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `name` - (Required, ForceNew) The log store, which is unique in the same project.
* `retention_period` - (Optional) The data retention time (in days). Valid values: [1-3650]. Default to `30`. Log store data will be stored permanently when the value is `3650`.
* `shard_count` - (Optional) The number of shards in this log store. Default to 2. You can modify it by "Split" or "Merge" operations. [Refer to details](https://www.alibabacloud.com/help/doc-detail/28976.htm)
* `auto_split` - (Optional) Determines whether to automatically split a shard. Default to `false`.
* `max_split_shard_count` - (Optional) The maximum number of shards for automatic split, which is in the range of 1 to 64. You must specify this parameter when autoSplit is true.
* `append_meta` - (Optional) Determines whether to append log meta automatically. The meta includes log receive time and client IP address. Default to `true`.
* `enable_web_tracking` - (Optional) Determines whether to enable Web Tracking. Default `false`.
* `encryption` (ForceNew, Optional, Available in 1.124.0+) Determines whether to automatically encryption,Default to `false`, only supported at creation time.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It formats of `<project>:<name>`.
* `project` - The project name.
* `name` - Log store name.
* `retention_period` - The data retention time.
* `shard_count` - The number of shards.
* `auto_split` - Determines whether to automatically split a shard.
* `max_split_shard_count` - The maximum number of shards for automatic split.
* `append_meta` - Determines whether to append log meta automatically.
* `enable_web_tracking` - Determines whether to enable Web Tracking.
* `encryption` - Determines whether to automatically encryption.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create`  - (Defaults to 3 mins) Used when Creating LogStore. 
* `delete`  - (Defaults to 3 mins) Used when Deleting LogStore.
* `read`    - (Defaults to 2 mins) Used when Reading LogStore.

## Import

Log store can be imported using the id, e.g.

```
$ terraform import alicloud_log_store.example tf-log:tf-log-store
```
