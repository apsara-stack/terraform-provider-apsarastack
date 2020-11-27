//

package main

import (
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: apsarastack.Provider,
	})
}
