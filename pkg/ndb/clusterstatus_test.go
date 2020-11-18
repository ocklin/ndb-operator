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

package ndb

import (
	"encoding/json"
	"fmt"
	"testing"
)

func nodeTypeFromNodeId(mgmNodeCount, dataNodeCount, apiNodeCount, nodeId int) int {

	if nodeId <= mgmNodeCount {
		return MgmNodeTypeId
	}
	if nodeId <= mgmNodeCount+dataNodeCount {
		return DataNodeTypeId
	}
	return ApiNodeTypeId
}

func Test_AddNodesByTLA(t *testing.T) {

	cs := NewClusterStatus(8)

	cs.SetNodeTypeFromTLA(1, "MGM")
	cs.SetNodeTypeFromTLA(2, "MGM")
	cs.SetNodeTypeFromTLA(3, "NDB")
	cs.SetNodeTypeFromTLA(4, "NDB")
	cs.SetNodeTypeFromTLA(5, "NDB")
	cs.SetNodeTypeFromTLA(6, "API")
	cs.SetNodeTypeFromTLA(7, "API")
	cs.SetNodeTypeFromTLA(8, "API")

	dnCnt := 0
	mgmCnt := 0
	apiCnt := 0
	for _, ns := range *cs {
		if ns.isApiNode() {
			apiCnt++
		}
		if ns.isDataNode() {
			dnCnt++
		}
		if ns.isMgmNode() {
			mgmCnt++
		}
	}

	if dnCnt != 3 && apiCnt != 3 && mgmCnt != 2 {
		t.Errorf("Wrong node type count")
	}

	err := cs.SetNodeTypeFromTLA(2, "MGM__")

	if err == nil {
		t.Errorf("Wrong node type string doesn't produce error")
	}
}

func Test_ClusterIsDegraded(t *testing.T) {

	cs := NewClusterStatus(8)

	// (!) start at 1
	for nodeId := 1; nodeId <= 8; nodeId++ {
		ns := &NodeStatus{
			NodeId:          nodeId,
			NodeType:        nodeTypeFromNodeId(2, 4, 2, nodeId),
			SoftwareVersion: "8.0.22",
			IsConnected:     true,
		}
		(*cs)[nodeId] = ns
	}

	for nodeId, ns := range *cs {
		s, _ := json.Marshal(ns)
		fmt.Printf("%d - %s\n", nodeId, s)
	}

	if cs.IsClusterDegraded() {
		t.Errorf("Cluster is not degraded but reported degraded.")
	}

	if ns, ok := (*cs)[7]; ok {
		(*ns).IsConnected = false
	} else {
		t.Errorf("Defined node id 1 not found.")
		return
	}
	if cs.IsClusterDegraded() {
		t.Errorf("Cluster is not degraded if 1 API Node is down but reported degraded.")
	}

	if ns, ok := (*cs)[1]; ok {
		(*ns).IsConnected = false
	} else {
		t.Errorf("Defined node id 1 not found.")
		return
	}

	if !cs.IsClusterDegraded() {
		t.Errorf("Cluster is degraded but reported not degraded.")
	}

}
