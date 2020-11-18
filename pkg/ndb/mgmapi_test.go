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
	"fmt"
	"testing"
)

func TestGetStatus(t *testing.T) {
	api := &Mgmclient{}

	err := api.Connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.Disconnect()

	nodeStatus, err := api.GetStatus()
	if err != nil {
		t.Errorf("get status failed: %s", err)
		return
	}

	for s, v := range *nodeStatus {
		fmt.Println(s, v)
	}
}

func TestGetOwnNodeId(t *testing.T) {

	api := &Mgmclient{}

	err := api.Connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.Disconnect()

	nodeid, err := api.GetOwnNodeId()
	if err != nil {
		t.Errorf("get status failed: %s", err)
		return
	}

	fmt.Printf("Own nodeid: %d\n", nodeid)
}

func TestStopNodes(t *testing.T) {
	api := &Mgmclient{}

	err := api.Connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.Disconnect()

	nodeIds := []int{2}
	disconnect, err := api.StopNodes(&nodeIds)
	if err != nil {
		t.Errorf("stop failed : %s", err)
		return
	}

	if disconnect {
		fmt.Println("Disconnect")
	}

}

func TestGetConfig(t *testing.T) {
	api := &Mgmclient{}

	err := api.Connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.Disconnect()

	_, err = api.GetConfig()
	if err != nil {
		t.Errorf("getting config failed : %s", err)
		return
	}

	_, err = api.GetConfigFromNode(2)
	if err != nil {
		t.Errorf("getting config failed : %s", err)
		return
	}
}

func TestShowConfig(t *testing.T) {
	api := &Mgmclient{}

	err := api.Connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.Disconnect()

	err = api.showConfig()
	if err != nil {
		t.Errorf("getting config failed : %s", err)
		return
	}
}

func TestShowVariables(t *testing.T) {
	api := &Mgmclient{}

	err := api.Connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.Disconnect()

	nodeid, err := api.GetOwnNodeId()
	if err != nil {
		t.Errorf("show variables failed: %s", err)
		return
	}

	if nodeid == 0 {
		t.Errorf("show variables failed with wrong or unknown node id: %d", nodeid)
	}
}

func TestConnectWantedNodeId(t *testing.T) {
	api := &Mgmclient{}

	wantedNodeId := 2
	err := api.ConnectToNodeId(wantedNodeId)
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.Disconnect()

	nodeid, err := api.GetOwnNodeId()
	if err != nil {
		t.Errorf("show variables failed: %s", err)
		return
	}

	if nodeid != wantedNodeId {
		t.Errorf("Connecting to wanted node id %d failed with wrong or unknown node id: %d", wantedNodeId, nodeid)
	}
}
