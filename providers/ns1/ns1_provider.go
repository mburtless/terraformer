// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package ns1

import (
	"os"
	"errors"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
		"github.com/GoogleCloudPlatform/terraformer/terraform_utils/provider_wrapper"
)

type NS1Provider struct {
	terraform_utils.Provider
	apiKey string
}

func (p *NS1Provider) Init(args []string) error {
	if k := os.Getenv("NS1_APIKEY"); k != "" {
		p.apiKey = k
	} else {
		return errors.New("set NS1_APIKEY env var")
	}

	return nil
}

func (p *NS1Provider) GetName() string {
	return "ns1"
}

func (p *NS1Provider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			"ns1": map[string]interface{}{
				"version": provider_wrapper.GetProviderVersion(p.GetName()),
				"apikey": p.apiKey,
			},
		},
	}
}

func (NS1Provider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{}
}


func (p *NS1Provider) GetSupportedService() map[string]terraform_utils.ServiceGenerator {
	return map[string]terraform_utils.ServiceGenerator{
		"zone": &ZoneGenerator{},
	}
}

func (p *NS1Provider) InitService(serviceName string) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New("ns1: " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"apikey": p.apiKey,
	})
	return nil
}