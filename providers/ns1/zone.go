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
	"log"
	"strings"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

type ZoneGenerator struct {
	NS1Service
}

func (z ZoneGenerator) createResources(zoneList []*dns.Zone) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, zone := range zoneList {
		name := strings.ReplaceAll(zone.Zone, ".", "_")
		resources = append(resources, terraform_utils.NewResource(
			zone.ID,
			name,
			"ns1_zone",
			"ns1",
			map[string]string{
				"zone": zone.Zone,
			},
			[]string{},
			map[string]interface{}{},
		))
	}
	log.Println(resources)
	return resources
}

func (z *ZoneGenerator) InitResources() error {
	client := z.generateClient()
	zones, _, err := client.Zones.List()
	if err != nil {
		log.Println(err)
		return err
	}
	z.Resources = z.createResources(zones)
	return nil
}