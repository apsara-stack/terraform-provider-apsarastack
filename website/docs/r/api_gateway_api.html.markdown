---
subcategory: "API Gateway"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_api_gateway_api"
sidebar_current: "docs-apsarastack-resource-api-gateway-api"
description: |-
  Provides a Apsarastack Api Gateway Api Resource.
---

# apsarastack_api_gateway_api

Provides an api resource.When you create an API, you must enter the basic information about the API, and define the API request information, the API backend service and response information.

For information about Api Gateway Api and how to use it, see [Create an API](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/create-an-api-dev1.html?spm=a2c4g.14484438.10001.154)

-> **NOTE:** Terraform will auto build api while it uses `apsarastack_api_gateway_api` to build api.

## Example Usage

Basic Usage

```

variable "name" {
	  default = "tf_testAccApiGatewayApi_7875290"
	}

	variable "apigateway_group_description_test" {
	  default = "tf_testAcc_api group description"
	}
	
	resource "apsarastack_api_gateway_group" "default" {
	  name = "${var.name}"
	  description = "${var.apigateway_group_description_test}"
	}
	

resource "apsarastack_api_gateway_api" "default" {
  name = "${apsarastack_api_gateway_group.default.name}"
  group_id = "${apsarastack_api_gateway_group.default.id}"
  description = "tf_testAcc_api description"
  auth_type = "APP"
  request_config {
    protocol = "HTTP"
    method = "GET"
    path = "/test/path/vpc"
    mode = "MAPPING"
  }
  
  service_type = "FunctionCompute"
  fc_service_config {
    region = "cn-qingdao-env17-d01"
    function_name = "tf_testAccApiGatewayApi_7875290Func"
    service_name = "tf_testAccApiGatewayApi_7875290"
    timeout = "20"
    arn_role = "cloudapi-openapi"
  }
  
  request_parameters {
    name = "testparam"
    type = "STRING"
    required = "OPTIONAL"
    in = "QUERY"
    in_service = "QUERY"
    name_service = "testparams"
  }
  
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the api gateway api. Defaults to null.
* `group_id` - (Required, ForcesNew) The api gateway that the api belongs to. Defaults to null.
* `description` - (Required) The description of the api. Defaults to null.
* `auth_type` - (Required) The authorization Type including APP and ANONYMOUS. Defaults to null.
* `request_config` - (Required, Type: list) Request_config defines how users can send requests to your API.
* `service_type` - (Required) The type of backend service. Type including HTTP,VPC and MOCK. Defaults to null.
* `fc_service_config` - (Optional, Type: list) fc_service_config defines the config when service_type selected 'FunctionCompute'.
* `request_parameters` - (Required, Type: list) request_parameters defines the request parameters of the api.


### Block request_config

The request_config mapping supports the following:

* `protocol` - (Required) The protocol of api which supports values of 'HTTP','HTTPS' or 'HTTP,HTTPS'.
* `method` - (Required) The method of the api, including 'GET','POST','PUT' etc.
* `path` - (Required) The request path of the api.
* `mode` - (Required) The mode of the parameters between request parameters and service parameters, which support the values of 'MAPPING' and 'PASSTHROUGH'.



### Block fc_vpc_service_config

The fc_service_config mapping supports the following:

* `region` - (Required) The region that the function compute service belongs to.
* `function_name` - (Required) The function name of function compute service.
* `service_name` - (Required) The service name of function compute service.
* `arn_role` - (Optional) RAM role arn attached to the Function Compute service. This governs both who / what can invoke your Function, as well as what resources our Function has access to. See [User Permissions](https://www.alibabacloud.com/help/doc-detail/52885.htm) for more details.
* `timeout` - (Required) Backend service time-out time; unit: millisecond.



### Block request_parameters

The request_parameters mapping supports the following:

* `name` - (Required) Request's parameter name.
* `type` - (Required) Parameter type which supports values of 'STRING','INT','BOOLEAN','LONG',"FLOAT" and "DOUBLE".
* `required` - (Required) Parameter required or not; values: REQUIRED and OPTIONAL.
* `in` - (Required) Request's parameter location; values: BODY, HEAD, QUERY, and PATH.
* `in_service` - (Required) Backend service's parameter location; values: BODY, HEAD, QUERY, and PATH.
* `name_service` - (Required) Backend service's parameter name.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the api resource of api gateway.
* `api_id` - The ID of the api of api gateway.

## Import

Api gateway api can be imported using the id.Format to `<API Group Id>:<API Id>` e.g.

```
$ terraform import alicloud_api_gateway_api.example "ab2351f2ce904edaa8d92a0510832b91:e4f728fca5a94148b023b99a3e5d0b62"
```
