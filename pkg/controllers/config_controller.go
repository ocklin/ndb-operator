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

package controllers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ocklin/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	"github.com/ocklin/ndb-operator/pkg/constants"
)

type ConfigControl struct {
}

func getMgmdHostname(ndb *v1alpha1.Ndb, count int) string {
	dnsZone := fmt.Sprintf("%s.svc.cluster.local", ndb.Namespace)
	mgmHostname := fmt.Sprintf("%s-%d.%s.%s", ndb.Name+"-mgmd", count, ndb.Name, dnsZone)
	return mgmHostname
}

func getNdbdHostname(ndb *v1alpha1.Ndb, count int) string {
	dnsZone := fmt.Sprintf("%s.svc.cluster.local", ndb.Namespace)
	mgmHostname := fmt.Sprintf("%s-%d.%s.%s", ndb.Name+"-ndbd", count, ndb.Name, dnsZone)
	return mgmHostname
}

func (cc ConfigControl) WriteConfig(ndb *v1alpha1.Ndb) error {

	//filename := constants.DataDir + "/config.ini"
	filename := "config.ini"

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

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

	noofrepl := fmt.Sprintf("%d", *ndb.Spec.Ndbd.NoOfReplicas)
	//noofrepl := strconv.Itoa(*ndb.Spec.Ndbd.NoOfReplicas)
	configString := strings.ReplaceAll(defaultSections, "{{$noofreplicas}}", noofrepl)
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

	for i := 0; i < int(*ndb.Spec.Ndbd.NodeCount); i++ {

		ns := ndbdSection
		ns = strings.ReplaceAll(ns, "{{$nodeId}}", strconv.Itoa(nodeId))
		ns = strings.ReplaceAll(ns, "{{$hostname}}", getNdbdHostname(ndb, i))
		ns = strings.ReplaceAll(ns, "{{$datadir}}", constants.DataDir)

		configString += ns
		nodeId++
		configString += "\n"
	}

	/* pure estetics - trim whitespace from lines */
	s := strings.Split(configString, "\n")
	configString = ""
	for _, line := range s {
		configString += strings.TrimSpace(line) + "\n"
	}

	_, err = f.WriteString(configString)
	if err != nil {
		return err
	}

	return nil
}
