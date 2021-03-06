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
	"testing"
)

func Test_GetConfigHashAndGenerationFromConfig(t *testing.T) {

	// just testing actual extraction of ConfigHash - not ini-file reading
	testini := `
	;
	; this is a header section
	;
	;ConfigHash=asdasdlkajhhnxh=?   
	;               notice this ^
	[system]
	ConfigGenerationNumber=4711
	`

	hash, generation, err := GetConfigHashAndGenerationFromConfig(testini)

	if err != nil {
		t.Errorf("extracting hash or generation failed : %s", err)
	}

	if hash != "asdasdlkajhhnxh=?" {
		t.Fail()
	}
	if generation != 4711 {
		t.Errorf("Wrong generation :" + string(generation))
	}
}
