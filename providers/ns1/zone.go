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
	"fmt"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

type DNSGenerator struct {
	NS1Service
}

// TODO: DELETE ME
func (z DNSGenerator) createZonesResourceOld(zoneList []*dns.Zone) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, zone := range zoneList {
		name := strings.ReplaceAll(zone.Zone, ".", "_")
		r := terraform_utils.NewResource(
			zone.ID,
			name,
			"ns1_zone",
			"ns1",
			map[string]string{
				"zone": zone.Zone,
			},
			[]string{},
			map[string]interface{}{},
			// add directly to state, not config
			//map[string]interface{}{"autogenerate_ns_record": true},
		)
		// remove secondaries networks from config
		//r.IgnoreKeys = append(r.IgnoreKeys, "networks")
		resources = append(resources, r)
	}
	//log.Println(resources)
	return resources
}

func (z DNSGenerator) createZonesResource(client *api.Client, zone *dns.Zone) ([]terraform_utils.Resource, error) {
	name := strings.ReplaceAll(zone.Zone, ".", "_")
	r := terraform_utils.NewResource(
		zone.ID,
		name,
		"ns1_zone",
		"ns1",
		map[string]string{
			"zone": zone.Zone,
		},
		[]string{},
		map[string]interface{}{},
		// add directly to state, not config
		//map[string]interface{}{"autogenerate_ns_record": true},
	)
	return []terraform_utils.Resource{r}, nil
}

func (z DNSGenerator) createRecordsResources(client *api.Client, zone *dns.Zone) ([]terraform_utils.Resource, error) {
	resources := []terraform_utils.Resource{}
	rTypesToIgnore := []string{"RRSIG", "DNSKEY"}

	// get zone for list of records
	zoneDetails, _, err := client.Zones.Get(zone.Zone)
	if err != nil {
		log.Println(err)
		return resources, err
	}

RECORDS:
	for _, record := range zoneDetails.Records {
		for _, t := range rTypesToIgnore {
			if record.Type == t {
				continue RECORDS
			}
		}

		name := strings.ReplaceAll(record.Domain, ".", "_")
		r := terraform_utils.NewResource(
			record.ID,
			fmt.Sprintf("%s_%s", record.Type, name),
			"ns1_record",
			"ns1",
			map[string]string{
				"zone":   zone.Zone,
				"domain": record.Domain,
				"type":   record.Type,
			},
			[]string{},
			map[string]interface{}{},
		)
		resources = append(resources, r)
	}

	return resources, nil
}

func (z *DNSGenerator) InitResources() error {
	client := z.generateClient()
	zones, _, err := client.Zones.List()
	if err != nil {
		log.Println(err)
		return err
	}
	//z.Resources = z.createZonesResource(zones)

	funcs := []func(*api.Client, *dns.Zone) ([]terraform_utils.Resource, error){
		z.createZonesResource,
		z.createRecordsResources,
	}

	//loop through zones
	for _, zone := range zones {
		for _, f := range funcs {
			tmpRes, err := f(client, zone)
			if err != nil {
				log.Println(err)
				return err
			}
			z.Resources = append(z.Resources, tmpRes...)
		}
	}
	return nil
}

func (z *DNSGenerator) PostConvertHook() error {
	for i, r := range z.Resources {
		if r.InstanceInfo.Type == "ns1_zone" {
			// delete networks from any secondaries in config
			if secondaries, ok := z.Resources[i].Item["secondaries"]; ok {
				for _, secondary := range secondaries.([]interface{}) {
					delete(secondary.(map[string]interface{}), "networks")
				}
			}

			// add autogenerate_ns_record default to state
			z.Resources[i].InstanceState.Attributes["autogenerate_ns_record"] = "true"
		}
	}
	return nil
}
