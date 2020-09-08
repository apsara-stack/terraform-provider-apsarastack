---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_server_certificate"
sidebar_current: "docs-apsarastack-resource-slb-server-certificate"
description: |-
  Provides a Load Banlancer Server Certificate resource.
---

# apsarastack\_slb\_server\_certificate

A Load Balancer Server Certificate is an ssl Certificate used by the listener of the protocol https.

For information about slb and how to use it, see [What is Server Load Balancer](https://apsarastackdocument.oss-cn-hangzhou.aliyuncs.com/01_ApsaraStackEnterprise/V3.11.0-intl-en/Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf?spm=a3c0i.214467.3807842930.7.61e76bdb1JWVyX&file=Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf).

For information about Server Certificate and how to use it, see [Configure Server Certificate](https://apsarastackdocument.oss-cn-hangzhou.aliyuncs.com/01_ApsaraStackEnterprise/V3.11.0-intl-en/Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf?spm=a3c0i.214467.3807842930.7.61e76bdb1JWVyX&file=Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf).


## Example Usage

* using server_certificate/private content as string example

```
# create a server certificate
resource "apsarastack_slb_server_certificate" "foo" {
  name               = "slbservercertificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
```

* using server_certificate/private file example

```
# create a server certificate
resource "apsarastack_slb_server_certificate" "foo" {
  name               = "slbservercertificate"
  server_certificate = "${file("${path.module}/server_certificate.pem")}"
  private_key        = "${file("${path.module}/private_key.pem")}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the Server Certificate.
* `server_certificate` - (Optional, ForceNew) the content of the ssl certificate. where `apsarastack_certificate_id` is null, it is required, otherwise it is ignored.
* `private_key` - (Optional, ForceNew) the content of privat key of the ssl certificate specified by `server_certificate`. where `apsarastack_certificate_id` is null, it is required, otherwise it is ignored.
* `resource_group_id` - (Optional, ForceNew, Available in 1.58.0+) The Id of resource group which the slb server certificate belongs.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.
## Attributes Reference

The following attributes are exported:

* `id` - The Id of Server Certificate (SSL Certificate).

## Import

Server Load balancer Server Certificate can be imported using the id, e.g.

```
$ terraform import apsarastack_slb_server_certificate.example abc123456
```
