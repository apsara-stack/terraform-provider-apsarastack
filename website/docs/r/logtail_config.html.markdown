---
subcategory: "Log Service (SLS)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_logtail_config"
sidebar_current: "docs-apsarastack-resource-logtail-config"
description: |-
  Provides a Apsarastack logtail config resource.
---

# apsarastack\_logtail\_config

The Logtail access service is a log collection agent provided by Log Service. 
You can use Logtail to collect logs from servers such as Apsarastack Cloud Elastic
Compute Service (ECS) instances in real time in the Log Service console.
)

## Example Usage

Basic Usage

```
resource "apsarastack_log_project" "example" {
  name        = "test-tf"
  description = "create by terraform"
}
resource "apsarastack_log_store" "example" {
  project               = "${apsarastack_log_project.example.name}"
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "apsarastack_logtail_config" "example" {
  project      = "${apsarastack_log_project.example.name}"
  logstore     = "${apsarastack_log_store.example.name}"
  input_type   = "file"
  name         = "tf-log-config"
  output_type  = "LogService"
  input_detail = "${file("config.json")}"
}
```


## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `logstore` - (Required, ForceNew) The log store name to the query index belongs.
* `input_type` - (Required) The input type. Currently only two types of files and plugin are supported.
* `name` - (Required, ForceNew) The Logtail configuration name, which is unique in the same project.
* `output_type` - (Required) The output type. Currently, only LogService is supported.
* `input_detail` - (Required) The logtail configure the required JSON files. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log store index. It formats of `<project>:<logstore>:<config_name>`.

