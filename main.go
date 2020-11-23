//

package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: apsarastack.Provider,
	})
}
