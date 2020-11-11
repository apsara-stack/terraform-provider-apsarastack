package ascm

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"reflect"
)

type Client struct {
	sdk.Client
}

func SetClientProperty(client *Client, propertyName string, propertyValue interface{}) {
	v := reflect.ValueOf(client).Elem()
	if v.FieldByName(propertyName).IsValid() && v.FieldByName(propertyName).CanSet() {
		v.FieldByName(propertyName).Set(reflect.ValueOf(propertyValue))
	}
}
func SetEndpointDataToClient(client *Client) {
	SetClientProperty(client, "EndpointMap", GetEndpointMap())
	SetClientProperty(client, "EndpointType", GetEndpointType())
}

func NewClientWithOptions(regionId string, config *sdk.Config, credential auth.Credential) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithOptions(regionId, config, credential)
	SetEndpointDataToClient(client)
	return
}

func NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret string) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithAccessKey(regionId, accessKeyId, accessKeySecret)
	return
}
