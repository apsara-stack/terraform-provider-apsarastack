---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_listener"
sidebar_current: "docs-apsarastack-resource-slb-listener"
description: |-
  Provides an Application Load Balancer resource.
---

# apsarastack\_slb\_listener

Provides an Application Load Balancer Listener resource.

For information about slb and how to use it, see [What is Server Load Balancer](https://apsarastackdocument.oss-cn-hangzhou.aliyuncs.com/01_ApsaraStackEnterprise/V3.11.0-intl-en/Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf?spm=a3c0i.214467.3807842930.7.61e76bdb1JWVyX&file=Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf).


## Example Usage

```
data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "apsarastack_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "apsarastack_vswitch" "vsw" {
  name       = "vsw"
  vpc_id            = apsarastack_vpc.vpc.id
  cidr_block        = apsarastack_vpc.vpc.cidr_block
  availability_zone =  "${data.apsarastack_zones.default.zones.0.id}"
}
resource "apsarastack_slb" "slb" {
  name          = "slb"
  vswitch_id    = apsarastack_vswitch.vsw.id
}

resource "apsarastack_slb_server_certificate" "servercertificate" {
  name               = "slbservercertificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}

resource "apsarastack_slb_listener" "listener" {
load_balancer_id          = apsarastack_slb.slb.id
backend_port              = 80
frontend_port             = 80
protocol                  = "http"
bandwidth                 = 10
server_certificate_id     =apsarastack_slb_server_certificate.servercertificate.id
sticky_session            = "on"
sticky_session_type       = "insert"
cookie_timeout            = 86400
cookie                    = "testslblistenercookie"
health_check              = "on"
health_check_domain       = "ali.com"
health_check_uri          = "/cons"
health_check_connect_port = 20
healthy_threshold         = 8
unhealthy_threshold       = 8
health_check_timeout      = 8
health_check_interval     = 5
health_check_http_code    = "http_2xx,http_3xx"
x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new listener.
* `frontend_port` - (Required, ForceNew) Port used by the Server Load Balancer instance frontend. Valid value range: [1-65535].
* `backend_port` - (Required, ForceNew) Port used by the Server Load Balancer instance backend. Valid value range: [1-65535].
* `protocol` - (Optional, ForceNew) The protocol to listen on. Valid values are [`http`, `https`, `tcp`, `udp`].
* `bandwidth` - (Required) Bandwidth peak of Listener. For the public network instance charged per traffic consumed, the Bandwidth on Listener can be set to -1, indicating the bandwidth peak is unlimited. Valid values are [-1, 1-5000] in Mbps.
* `description` - (Optional) The description of slb listener. This description can have a string of 1 to 80 characters. Default value: null.
* `scheduler` - (Optional) Scheduling algorithm, Valid values are `wrr`, `rr` and `wlc`.  Default to "wrr".
* `sticky_session` - (Required) Whether to enable session persistence, Valid values are `on` and `off`. Default to `off`.
* `sticky_session_type` - (Optional) Mode for handling the cookie. If `sticky_session` is "on", it is mandatory. Otherwise, it will be ignored. Valid values are `insert` and `server`. `insert` means it is inserted from Server Load Balancer; `server` means the Server Load Balancer learns from the backend server.
* `cookie_timeout` - (Optional) Cookie timeout. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "insert". Otherwise, it will be ignored. Valid value range: [1-86400] in seconds.
* `cookie` - (Optional) The cookie configured on the server. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "server". Otherwise, it will be ignored. Valid value：String in line with RFC 2965, with length being 1- 200. It only contains characters such as ASCII codes, English letters and digits instead of the comma, semicolon or spacing, and it cannot start with $.
* `persistence_timeout` - (Optional) Timeout of connection persistence. Valid value range: [0-3600] in seconds. Default to 0 and means closing it.
* `health_check` - (Required) Whether to enable health check. Valid values are`on` and `off`. TCP and UDP listener's HealthCheck is always on, so it will be ignore when launching TCP or UDP listener.
* `health_check_type` - (Optional) Type of health check. Valid values are: `tcp` and `http`. Default to `tcp` . TCP supports TCP and HTTP health check mode, you can select the particular mode depending on your application.
* `health_check_domain` - (Optional) Domain name used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and only characters such as letters, digits, ‘-‘ and ‘.’ are allowed. When it is not set or empty,  Server Load Balancer uses the private network IP address of each backend server as Domain used for health check.
* `health_check_uri` - (Optional) URI used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and it must start with /. Only characters such as letters, digits, ‘-’, ‘/’, ‘.’, ‘%’, ‘?’, #’ and ‘&’ are allowed.
* `health_check_connect_port` - (Optional) Port used for health check. Valid value range: [1-65535]. Default to "None" means the backend server port is used.
* `healthy_threshold` - (Optional) Threshold determining the result of the health check is success. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `unhealthy_threshold` - (Optional) Threshold determining the result of the health check is fail. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `health_check_timeout` - (Optional) Maximum timeout of each health check response. It is required when `health_check` is on. Valid value range: [1-300] in seconds. Default to 5. Note: If `health_check_timeout` < `health_check_interval`, its will be replaced by `health_check_interval`.
* `health_check_interval` - (Optional) Time interval of health checks. It is required when `health_check` is on. Valid value range: [1-50] in seconds. Default to 2.
* `health_check_http_code` - (Optional) Regular health check HTTP status code. Multiple codes are segmented by “,”. It is required when `health_check` is on. Default to `http_2xx`.  Valid values are: `http_2xx`,  `http_3xx`, `http_4xx` and `http_5xx`.
* `health_check_method` - (Optional) HealthCheckMethod used for health check.`http` and `https` support regions ap-northeast-1, ap-southeast-1, ap-southeast-2, ap-southeast-3, us-east-1, us-west-1, eu-central-1, ap-south-1, me-east-1, cn-huhehaote, cn-zhangjiakou, ap-southeast-5, cn-shenzhen, cn-hongkong, cn-qingdao, cn-chengdu, eu-west-1, cn-hangzhou", cn-beijing, cn-shanghai.This function does not support the TCP protocol .
* `server_certificate_id` - (Required) SLB Server certificate ID. It is required when `protocol` is `https`.
* `gzip` - (Optional) Whether to enable "Gzip Compression". If enabled, files of specific file types will be compressed, otherwise, no files will be compressed. Default to true.
* `x_forwarded_for` - (Optional) Whether to set additional HTTP Header field "X-Forwarded-For" (documented below).
* `established_timeout` - (Optional) Timeout of tcp listener established connection idle timeout. Valid value range: [10-900] in seconds. Default to 900.
* `server_group_id` - (Optional) the id of server group to be apply on the listener, is the id of resource `apsarastack_slb_server_group`.
* `listener_forward` - (Optional, ForceNew) Whether to enable http redirect to https, Valid values are `on` and `off`. Default to `off`.
* `forward_port` - (Optional, ForceNew) The port that http redirect to https.
* `health_check_method` - (Optional, ForceNew, Available in 1.70.0+) The method of health check. Valid values: ["head", "get"].
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

-> **NOTE:** Once enable the http redirect to https function, any parameters excepted forward_port,listener_forward,load_balancer_id,frontend_port,protocol will be ignored. More info, please refer to [Redirect http to https](https://apsarastackdocument.oss-cn-hangzhou.aliyuncs.com/01_ApsaraStackEnterprise/V3.11.0-intl-en/Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf?spm=a3c0i.214467.3807842930.7.61e76bdb1JWVyX&file=Alibaba%20Cloud%20Apsara%20Stack%20Enterprise%202001%2C%20Internal_%20V3.11.0%20Developer%20Guide%20-%20Cloud%20Essentials%20and%20Security%2020200513.pdf).


### Block x_forwarded_for

The x_forwarded_for mapping supports the following:

* `retrive_slb_ip` - (Optional) Whether to use the XForwardedFor_SLBIP header to obtain the public IP address of the SLB instance. Default to false.
* `retrive_slb_id` - (Optional) Whether to use the XForwardedFor header to obtain the ID of the SLB instance. Default to false.
* `retrive_slb_proto` - (Optional) Whether to use the XForwardedFor_proto header to obtain the protocol used by the listener. Default to false.

## Listener fields and protocol mapping

load balance support 4 protocal to listen on, they are `http`,`https`,`tcp`,`udp`, the every listener support which portocal following:

listener parameter | support protocol | value range |
------------- | ------------- | ------------- | 
backend_port | http & https & tcp & udp | 1-65535 | 
frontend_port | http & https & tcp & udp | 1-65535 |
protocol | http & https & tcp & udp |
bandwidth | http & https & tcp & udp | -1 / 1-5000 |
scheduler | http & https & tcp & udp | wrr rr or wlc |
sticky_session | http & https | on or off |
sticky_session_type | http & https | insert or server | 
cookie_timeout | http & https | 1-86400  | 
cookie | http & https |   | 
persistence_timeout | tcp & udp | 0-3600 | 
health_check | http & https | on or off | 
health_check_type | tcp | tcp or http | 
health_check_domain | http & https & tcp | 
health_check_method | http & https & tcp | 
health_check_uri | http & https & tcp |  | 
health_check_connect_port | http & https & tcp & udp | 1-65535 or -520 | 
healthy_threshold | http & https & tcp & udp | 1-10 | 
unhealthy_threshold | http & https & tcp & udp | 1-10 | 
health_check_timeout | http & https & tcp & udp | 1-300 |
health_check_interval | http & https & tcp & udp | 1-50 |
health_check_http_code | http & https & tcp | http_2xx,http_3xx,http_4xx,http_5xx | 
server_certificate_id | https |  |
gzip | http & https | true or false  |
x_forwarded_for | http & https |  |
established_timeout | tcp       | 10-900|
server_group_id    | http & https & tcp & udp | the id of resource apsarastack_slb_server_group |

The listener mapping supports the following:

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer listener. Its format as `<load_balancer_id>:<protocol>:<frontend_port>`. Before verson 1.57.1, the foramt as `<load_balancer_id>:<frontend_port>`.
* `load_balancer_id` - The Load Balancer ID which is used to launch a new listener.
* `frontend_port` - Port used by the Server Load Balancer instance frontend.
* `backend_port` - Port used by the Server Load Balancer instance backend.
* `protocol` - The protocol to listen on.
* `bandwidth` - Bandwidth peak of Listener.
* `scheduler` - Scheduling algorithm.
* `sticky_session` - Whether to enable session persistence.
* `sticky_session_type` - Mode for handling the cookie.
* `cookie_timeout` - Cookie timeout.
* `cookie` - The cookie configured on the server.
* `persistence_timeout` - Timeout of connection persistence.
* `health_check` - Whether to enable health check.
* `health_check_type` - Type of health check.
* `health_check_domain` - Domain name used for health check.
* `health_check_method` - HealthCheckMethod used for health check.
* `health_check_uri` - URI used for health check.
* `health_check_connect_port` - Port used for health check.
* `healthy_threshold` - Threshold determining the result of the health check is success.
* `unhealthy_threshold` - Threshold determining the result of the health check is fail.
* `health_check_timeout` - Maximum timeout of each health check response.
* `health_check_interval` - Time interval of health checks.
* `health_check_http_code` - Regular health check HTTP status code.
* `server_certificate_id` - (Optional) Security certificate ID.

