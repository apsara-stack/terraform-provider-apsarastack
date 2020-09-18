---
layout: "apsarastack"
page_title: "Provider: apsarastack"
sidebar_current: "docs-apsarastack-index"
description: |-
  The ApsaraStack provider is used to interact with many resources supported by ApsaraStack. The provider needs to be configured with the proper credentials before it can be used.
---

# ApsaraStack Cloud Provider

~> **News:** Currently, ApsaraStack Cloud has published [Terraform Module Web GUI](https://api.aliyun.com/#/cli?tool=Terraform) to
 help developers to use Terraform Module more simply and conveniently. Welcome to access it and let us know your more requirements!

The ApsaraStack Cloud provider is used to interact with the
many resources supported by ApsaraStack Cloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** From version 1.50.0, the provider start to support Terraform 0.12.x.


## Example Usage

```hcl
# Configure the ApsaraStack Provider
provider "apsarastack" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
  insecure    =  true
  proxy      = "${var.proxy}"
  domain     = "${var.domain}"
}


data "apsarastack_instance_types" "c2g4" {
  cpu_core_count = 2
  memory_size    = 4
}

data "apsarastack_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}

# Create a web server
resource "apsarastack_instance" "web" {
  image_id              = "${data.apsarastack_images.default.images.0.id}"
  internet_charge_type  = "PayByBandwidth"

  instance_type        = "${data.apsarastack_instance_types.c2g4.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  security_groups      = ["${apsarastack_security_group.default.id}"]
  instance_name        = "web"
  vswitch_id           = "vsw-abc12345"
}

# Create security group
resource "apsarastack_security_group" "default" {
  name        = "default"
  description = "default"
  vpc_id      = "vpc-abc12345"
}
```

## Authentication

The ApsaraStack provider accepts several ways to enter credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables
- ECS Role
- Assume role

### Static credentials

Static credentials can be provided by adding `access_key`, `secret_key` and `region` in-line in the
apsarastack provider block:

Usage:

```hcl
provider "apsarastack" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
  insecure    =  true
  proxy      = "${var.proxy}"
  domain     = "${var.domain}"
}

```

### Environment variables

You can provide your credentials via `APSARASTACK_ACCESS_KEY` and `APSARASTACK_SECRET_KEY`
environment variables, representing your ApsaraStack access key and secret key respectively.
`APSARASTACK_REGION` is also used, if applicable:

```hcl
provider "apsarastack" {}
```
Usage:

```shell
$ export APSARASTACK_ACCESS_KEY="anaccesskey"
$ export APSARASTACK_SECRET_KEY="asecretkey"
$ export APSARASTACK_REGION="cn-beijing"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the ApsaraStack Cloud
 `provider` block:

* `access_key` - This is the ApsaraStack access key. It must be provided, but
  it can also be sourced from the `APSARASTACK_ACCESS_KEY` environment variable, or via
  a dynamic access key if `ecs_role_name` is specified.

* `secret_key` - This is the ApsaraStack secret key. It must be provided, but
  it can also be sourced from the `APSARASTACK_SECRET_KEY` environment variable, or via
  a dynamic secret key if `ecs_role_name` is specified.
  
* `region` - This is the ApsaraStack region. It must be provided, but
  it can also be sourced from the `APSARASTACK_REGION` environment variables.

* `security_token` - ApsaraStack Security Token Service.
  It can be sourced from the `APSARASTACK_SECURITY_TOKEN` environment variable,  or via
  a dynamic security token if `ecs_role_name` is specified.

* `ecs_role_name` - "The RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the 'Access Control' section of the ApsaraStack Cloud console.",

* `skip_region_validation` - (Optional, Available in 1.52.0+) Skip static validation of region ID. Used by users of alternative ApsaraStackCloud-like APIs or users w/ access to regions that are not public (yet).

* `profile` - (Optional, Available in 1.49.0+) This is the ApsaraStack profile name as set in the shared credentials file. It can also be sourced from the `APSARASTACK_PROFILE` environment variable.

* `endpoints` - (Optional) An `endpoints` block (documented below) to support custom endpoints.

* `shared_credentials_file` - (Optional, Available in 1.49.0+) This is the path to the shared credentials file. It can also be sourced from the `APSARASTACK_SHARED_CREDENTIALS_FILE` environment variable. If this is not set and a profile is specified, ~/.aliyun/config.json will be used.

* `insecure` - (Optional) Use this to Trust self-signed certificates. It's typically used to allow insecure connections.

* `assume_role` - (Optional) An `assume_role` block (documented below). Only one `assume_role` block may be in the configuration.

* `protocol` - (Optional, Available in 1.72.0+) The Protocol of used by API request. Valid values: `HTTP` and `HTTPS`. Default to `HTTPS`.

* `configuration_source` - (Optional, Available in 1.56.0+) Use a string to mark a configuration file source, like `terraform-apsarastack-modules/terraform-apsarastack-ecs-instance` or `terraform-provider-apsarastack/examples/vpc`.
The length should not more than 64.

* `proxy` -  (Optional) Use this to set proxy connection.

* `domain` - (Optional) Use this to override the default domain. It's typically used to connect to custom domain.

## Testing

Credentials must be provided via the `APSARASTACK_ACCESS_KEY`, `APSARASTACK_SECRET_KEY` and `APSARASTACK_REGION` environment variables in order to run acceptance tests.
