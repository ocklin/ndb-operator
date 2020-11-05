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
	"testing"
	"time"
)

func TestGetStatus(t *testing.T) {
	api := &mgmclient{}

	err := api.connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.disconnect()

	for {
		err = api.getStatus()
		if err != nil {
			t.Errorf("get status failed: %s", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func TestRestart(t *testing.T) {
	api := &mgmclient{}

	err := api.connect()
	if err != nil {
		t.Errorf("Connection failed: %s", err)
		return
	}
	defer api.disconnect()

	err = api.restart()
	if err != nil {
		t.Errorf("restart failed : %s", err)
		return
	}
}
