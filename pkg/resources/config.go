// Copyright 2020 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ocklin/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	"github.com/ocklin/ndb-operator/pkg/constants"
	"github.com/ocklin/ndb-operator/pkg/helpers"
)

func GetConfigHashAndGenerationFromConfig(configStr string) (string, int64, error) {

	config, err := helpers.ParseString(configStr)

	var generation int64

	if err != nil {
		return "", generation, err
	}

	configHash := helpers.GetValueFromSingleSectionGroup(config, "header", "ConfigHash")
	generationStr := helpers.GetValueFromSingleSectionGroup(config, "system", "ConfigGenerationNumber")

	generation, _ = strconv.ParseInt(generationStr, 10, 64)

	return configHash, generation, nil
}

func getMgmdHostname(ndb *v1alpha1.Ndb, count int) string {
	dnsZone := fmt.Sprintf("%s.svc.cluster.local", ndb.Namespace)
	mgmHostname := fmt.Sprintf("%s-%d.%s.%s", ndb.Name+"-mgmd", count, ndb.GetManagementServiceName(), dnsZone)
	return mgmHostname
}

func getNdbdHostname(ndb *v1alpha1.Ndb, count int) string {
	dnsZone := fmt.Sprintf("%s.svc.cluster.local", ndb.Namespace)
	mgmHostname := fmt.Sprintf("%s-%d.%s.%s", ndb.Name+"-ndbd", count, ndb.GetDataNodeServiceName(), dnsZone)
	return mgmHostname
}

func GetConfigString(ndb *v1alpha1.Ndb) (string, error) {

	header := `
	# auto generated config.ini - do not edit
	#
	# ConfigHash={{$confighash}}
	`

	systemSection := `
	[system]
	ConfigGenerationNumber={{$configgeneration}}
	Name={{$clustername}}
	`

	defaultSections := `
  [ndbd default]
  NoOfReplicas={{$noofreplicas}}
  DataMemory=80M

  [tcp default]
  AllowUnresolvedHostnames=1`

	mgmdSection := `	
  [ndb_mgmd]
  NodeId={{$nodeId}}
  Hostname={{$hostname}}
  DataDir={{$datadir}}`

	ndbdSection := `
  [ndbd]
  NodeId={{$nodeId}}
  Hostname={{$hostname}}
  DataDir={{$datadir}}
  ServerPort=1186`

	// START of generation
	configString := ""

	// header
	hash := ndb.Status.ReceivedConfigHash
	configString += strings.ReplaceAll(header, "{{$confighash}}", hash)
	configString += "\n"

	// system section
	// TODO - this is wrong - needs to be reeived generation (as received config hash)
	generation := fmt.Sprintf("%d", ndb.ObjectMeta.Generation)
	syss := systemSection
	syss = strings.ReplaceAll(syss, "{{$configgeneration}}", generation)
	syss = strings.ReplaceAll(syss, "{{$clustername}}", ndb.Name)
	configString += syss + "\n"

	// ndbd default
	noofrepl := fmt.Sprintf("%d", *ndb.Spec.Ndbd.NoOfReplicas)
	//noofrepl := strconv.Itoa(*ndb.Spec.Ndbd.NoOfReplicas)
	configString += strings.ReplaceAll(defaultSections, "{{$noofreplicas}}", noofrepl)
	configString += "\n"

	/*
		TODO - how about hostname/nodeid stability when patching existing config?
	*/
	nodeId := int(1)
	for i := 0; i < int(*ndb.Spec.Mgmd.NodeCount); i++ {

		ms := mgmdSection
		ms = strings.ReplaceAll(ms, "{{$nodeId}}", strconv.Itoa(nodeId))
		ms = strings.ReplaceAll(ms, "{{$hostname}}", getMgmdHostname(ndb, i))
		ms = strings.ReplaceAll(ms, "{{$datadir}}", constants.DataDir)

		configString += ms
		nodeId++
		configString += "\n"
	}

	// data node sections
	for i := 0; i < int(*ndb.Spec.Ndbd.NodeCount); i++ {

		ns := ndbdSection
		ns = strings.ReplaceAll(ns, "{{$nodeId}}", strconv.Itoa(nodeId))
		ns = strings.ReplaceAll(ns, "{{$hostname}}", getNdbdHostname(ndb, i))
		ns = strings.ReplaceAll(ns, "{{$datadir}}", constants.DataDir)

		configString += ns
		nodeId++
		configString += "\n"
	}
	configString += "\n"

	// mysqld sections
	// at least 1 must be there in order to not fail ndb_mgmd start
	mysqlSections := 1
	if ndb.Spec.Mysqld.NodeCount != nil && int(*ndb.Spec.Mysqld.NodeCount) > 1 {
		// we alloc one section more than needed for internal purposes
		mysqlSections = int(*ndb.Spec.Mysqld.NodeCount) + 1
	}

	for i := 0; i < mysqlSections; i++ {
		configString += "[mysqld]\n"
	}

	/* pure estetics - trim whitespace from lines */
	s := strings.Split(configString, "\n")
	configString = ""
	for _, line := range s {
		configString += strings.TrimSpace(line) + "\n"
	}

	//klog.Infof("Config string: \n %s", configString)
	return configString, nil
}
