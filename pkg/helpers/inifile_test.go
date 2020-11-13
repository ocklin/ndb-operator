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

package helpers

import (
	"fmt"
	"os"
	"testing"
)

func TestReadInifile(t *testing.T) {

	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	testini := `
	;
	; this is a header section
	;
	; ConfigHash=asdasdlkajhhnxh=?   
	;               notice this ^

	; should not create 2nd header with just empty line

	[ndbd default]
	NoOfReplicas=2
	DataMemory=80M
	ServerPort=2202
	StartPartialTimeout=15000
	StartPartitionedTimeout=0
	
	[tcp default]
	AllowUnresolvedHostnames=1
	
	# more comments to be ignored

	[ndb_mgmd]
	NodeId=0
	Hostname=example-ndb-0.example-ndb.svc.default-namespace.com
	DataDir=/var/lib/ndb
	
	# key=value
	# comment with key value pair here should be ignored
	[ndbd]
	NodeId=1
	Hostname=example-ndb-0.example-ndb.svc.default-namespace.com
	DataDir=/var/lib/ndb
	ServerPort=1186
	
	[mysqld]
	NodeId=1
	Hostname=example-ndb-0.example-ndb.svc.default-namespace.com
	
	[mysqld]`

	l, err := f.WriteString(testini)
	if err != nil {
		t.Error(err)
		f.Close()
		return
	}

	if l == 0 {
		f.Close()
		t.Fail()
		return
	}
	f.Close()

	c, err := ParseFile("test.txt")

	if err != nil {
		t.Error(err)
		f.Close()
		return
	}

	if c == nil {
		t.Fail()
		return
	}

	t.Log("Iterating")
	for s, grp := range c.Groups {
		for _, sec := range grp {
			t.Log("[" + s + "]")
			for key, val := range sec {
				t.Log(key + ": " + val)
			}
		}
	}

	t.Fail()
}
