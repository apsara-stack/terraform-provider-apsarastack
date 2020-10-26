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
and each project can create multiple Logstores.

## Example Usage

Basic Usage

```
resource "apsarastack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "apsarastack_log_store" "example" {
  project               = "${apsarastack_log_project.example.name}"
  name                  = "tf-log-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
```


## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `name` - (Required, ForceNew) The log store, which is unique in the same project.
* `retention_period` - (Optional) The data retention time (in days). Valid values: [1-3650]. Default to 30. Log store data will be stored permanently when the value is "3650".
* `shard_count` - (Optional) The number of shards in this log store. Default to 2. You can modify it by "Split" or "Merge" operations. 
* `auto_split` - (Optional) Determines whether to automatically split a shard. Default to true.
* `max_split_shard_count` - (Optional) The maximum number of shards for automatic split, which is in the range of 1 to 64. You must specify this parameter when autoSplit is true.
* `append_meta` - (Optional) Determines whether to append log meta automatically. The meta includes log receive time and client IP address. Default to true.
* `enable_web_tracking` - (Optional) Determines whether to enable Web Tracking. Default false.

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

