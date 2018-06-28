// +build integration

// Space above here matters
// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/hyperledger/burrow/core"
	"github.com/hyperledger/burrow/core/integration"
	"github.com/hyperledger/burrow/execution/pbtransactor"
)

var privateAccounts = integration.MakePrivateAccounts(5) // make keys
var genesisDoc = integration.TestGenesisDoc(privateAccounts)
var inputAccount = &pbtransactor.InputAccount{Address: privateAccounts[0].Address().Bytes()}
var kern *core.Kernel

// Needs to be in a _test.go file to be picked up
func TestMain(m *testing.M) {
	returnValue := integration.TestWrapper(privateAccounts, genesisDoc, func(k *core.Kernel) int {
		kern = k
		return m.Run()
	})

	time.Sleep(3 * time.Second)
	os.Exit(returnValue)
}
