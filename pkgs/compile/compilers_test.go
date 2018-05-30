package compile

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/monax/bosmarmot/pkgs/util"
	"github.com/stretchr/testify/assert"
)

func TestRequestCreation(t *testing.T) {
	os.Chdir(testContractPath()) // important to maintain relative paths
	var err error
	contractCode := `pragma solidity ^0.4.0;

contract c {
    function f() {
        uint8[5] memory foo3 = [1, 1, 1, 1, 1];
    }
}`

	var testMap = map[string]*util.IncludedFiles{
		"27fbf28c5dfb221f98526c587c5762cdf4025e85809c71ba871caa2ca42a9d85.sol": {
			ObjectNames: []string{"c"},
			Script:      []byte(contractCode),
		},
	}

	req, err := CreateRequest("simpleContract.sol", "", false)

	if err != nil {
		t.Fatal(err)
	}
	if req.Libraries != "" {
		t.Errorf("Expected empty libraries, got %s", req.Libraries)
	}
	if req.Language != "sol" {
		t.Errorf("Expected Solidity file, got %s", req.Language)
	}
	if req.Optimize != false {
		t.Errorf("Expected false optimize, got true")
	}
	if !reflect.DeepEqual(req.Includes, testMap) {
		t.Errorf("Got incorrect Includes map, expected %v, got %v", testMap, req.Includes)
	}

}

func TestLocalMulti(t *testing.T) {
	ClearCache(util.SolcScratchPath)
	expectedSolcResponse := util.BlankSolcResponse()
	actualOutput, err := exec.Command("solc", "--combined-json", "bin,abi", "contractImport1.sol").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	warning, responseJSON := extractWarningJSON(strings.TrimSpace(string(actualOutput)))
	err = json.Unmarshal([]byte(responseJSON), expectedSolcResponse)

	respItemArray := make([]ResponseItem, 0)

	for contract, item := range expectedSolcResponse.Contracts {
		respItem := ResponseItem{
			Objectname: objectName(strings.TrimSpace(contract)),
			Bytecode:   trimAuxdata(strings.TrimSpace(item.Bin)),
			ABI:        strings.TrimSpace(item.Abi),
		}
		respItemArray = append(respItemArray, respItem)
	}
	expectedResponse := &Response{
		Objects: respItemArray,
		Warning: warning,
		Version: "",
		Error:   "",
	}
	ClearCache(util.SolcScratchPath)
	resp, err := RequestCompile("contractImport1.sol", false, "")
	if err != nil {
		t.Fatal(err)
	}
	fixupCompilersResponse(resp, "contractImport1.sol")
	allClear := true
	for _, object := range expectedResponse.Objects {
		if !contains(resp.Objects, object) {
			allClear = false
		}
	}
	if !allClear {
		t.Errorf("Got incorrect response, expected %v, \n\n got %v", expectedResponse, resp)
	}
	ClearCache(util.SolcScratchPath)
}

func TestLocalSingle(t *testing.T) {
	ClearCache(util.SolcScratchPath)
	expectedSolcResponse := util.BlankSolcResponse()

	shellCmd := exec.Command("solc", "--combined-json", "bin,abi", "simpleContract.sol")
	actualOutput, err := shellCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	warning, responseJSON := extractWarningJSON(strings.TrimSpace(string(actualOutput)))
	err = json.Unmarshal([]byte(responseJSON), expectedSolcResponse)

	respItemArray := make([]ResponseItem, 0)

	for contract, item := range expectedSolcResponse.Contracts {
		respItem := ResponseItem{
			Objectname: objectName(strings.TrimSpace(contract)),
			Bytecode:   trimAuxdata(strings.TrimSpace(item.Bin)),
			ABI:        strings.TrimSpace(item.Abi),
		}
		respItemArray = append(respItemArray, respItem)
	}
	expectedResponse := &Response{
		Objects: respItemArray,
		Warning: warning,
		Version: "",
		Error:   "",
	}
	ClearCache(util.SolcScratchPath)
	resp, err := RequestCompile("simpleContract.sol", false, "")
	if err != nil {
		t.Fatal(err)
	}
	fixupCompilersResponse(resp, "simpleContract.sol")
	assert.Equal(t, expectedResponse, resp)
	ClearCache(util.SolcScratchPath)
}

func TestFaultyContract(t *testing.T) {
	ClearCache(util.SolcScratchPath)
	var expectedSolcResponse Response

	actualOutput, err := exec.Command("solc", "--combined-json", "bin,abi", "faultyContract.sol").CombinedOutput()
	err = json.Unmarshal(actualOutput, expectedSolcResponse)
	t.Log(expectedSolcResponse.Error)
	resp, err := RequestCompile("faultyContract.sol", false, "")
	t.Log(resp.Error)
	if err != nil {
		if expectedSolcResponse.Error != resp.Error {
			t.Errorf("Expected %v got %v", expectedSolcResponse.Error, resp.Error)
		}
	}
	output := strings.TrimSpace(string(actualOutput))
	err = json.Unmarshal([]byte(output), expectedSolcResponse)
}

func testContractPath() string {
	baseDir, _ := os.Getwd()
	return filepath.Join(baseDir, "..", "..", "tests", "compilers_fixtures")
}

// The solidity 0.4.21 compiler appends something called auxdata to the end of the bin file (this is visible with
// solc --asm). This is a swarm hash of the metadata, and it's always at the end. This includes the path of the
// solidity source file, so it will differ.
func trimAuxdata(bin string) string {
	return bin[:len(bin)-86]
}

func extractWarningJSON(output string) (warning string, json string) {
	jsonBeginsCertainly := strings.Index(output, `{"contracts":`)

	if jsonBeginsCertainly > 0 {
		warning = output[:jsonBeginsCertainly]
		json = output[jsonBeginsCertainly:]
	} else {
		json = output
	}
	return
}

func fixupCompilersResponse(resp *Response, filename string) {
	for i := range resp.Objects {
		resp.Objects[i].Bytecode = trimAuxdata(resp.Objects[i].Bytecode)
	}
	// compilers changes the filename, change it back again in the warning
	re := regexp.MustCompile("[0-9a-f]+\\.sol")
	resp.Warning = re.ReplaceAllString(resp.Warning, filename)
}

func contains(s []ResponseItem, e ResponseItem) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
