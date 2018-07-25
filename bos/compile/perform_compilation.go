package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// New Request object from script and map of include files
func CompilerRequest(file string,
	includes map[string]*IncludedFiles, libs map[string]string, optimize bool) *Request {
	if includes == nil {
		includes = make(map[string]*IncludedFiles)
	}
	return &Request{
		Includes:  includes,
		Libraries: libs,
		Optimize:  optimize,
	}
}

// Compile request object
type Request struct {
	ScriptName string                    `json:"name"`
	Includes   map[string]*IncludedFiles `json:"includes"`  // our required files and metadata
	Libraries  map[string]string         `json:"libraries"` // string of libName:LibAddr separated by comma
	Optimize   bool                      `json:"optimize"`  // run with optimize flag
}

type BinaryRequest struct {
	BinaryFile string `json:"binary"`
	Libraries  string `json:"libraries"`
}

type SolidityInput struct {
	Language string                         `json:"language"`
	Sources  map[string]SolidityInputSource `json:"sources"`
	Settings struct {
		Libraries map[string]map[string]string `json:"libraries"`
		Optimizer struct {
			Enabled bool `json:"enabled"`
		} `json:"optimizer"`
		OutputSelection struct {
			File struct {
				OutputType []string `json:"*"`
			} `json:"*"`
		} `json:"outputSelection"`
	} `json:"settings"`
}

type SolidityInputSource struct {
	Content string   `json:"content,omitempty"`
	Urls    []string `json:"urls,omitempty"`
}

type SolidityOutput struct {
	Contracts map[string]map[string]SolidityOutputContract
	Errors    []struct {
		Component        string
		FormattedMessage string
		Message          string
		Severity         string
		Type             string
	}
}

type SolidityOutputContract struct {
	Abi json.RawMessage
	Evm struct {
		Bytecode struct {
			Object         string
			Opcodes        string
			LinkReferences json.RawMessage
		}
	}
	Devdoc   json.RawMessage
	Userdoc  json.RawMessage
	Metadata string
}

type Response struct {
	Objects []ResponseItem `json:"objects"`
	Warning string         `json:"warning"`
	Version string         `json:"version"`
	Error   string         `json:"error"`
}

type BinaryResponse struct {
	Binary string `json:"binary"`
	Error  string `json:"error"`
}

// Compile response object
type ResponseItem struct {
	Objectname string                 `json:"objectname"`
	Binary     SolidityOutputContract `json:"binary"`
}

const (
	AddressLength    = 40
	RelocationLength = 20
)

func RequestBinaryLinkage(file string, libraries map[string]string) (*BinaryResponse, error) {
	//Create Binary Request, send it off
	codeB, err := ioutil.ReadFile(file)
	if err != nil {
		return &BinaryResponse{}, err
	}
	contract := SolidityOutputContract{}
	err = json.Unmarshal(codeB, &contract)
	if err != nil {
		return &BinaryResponse{}, err
	}
	bin := contract.Evm.Bytecode.Object
	if !strings.Contains(bin, "_") {
		return &BinaryResponse{
			Binary: bin,
			Error:  "",
		}, nil
	}
	var links map[string]map[string]struct{ Start, Length int }
	err = json.Unmarshal(contract.Evm.Bytecode.LinkReferences, &links)
	if err != nil {
		return &BinaryResponse{}, err
	}
	for _, f := range links {
		for name, relo := range f {
			addr, ok := libraries[name]
			if !ok {
				return &BinaryResponse{}, fmt.Errorf("library %s is not defined", name)
			}
			if relo.Length != RelocationLength {
				return &BinaryResponse{}, fmt.Errorf("linkReference should be %d bytes long, not %d", RelocationLength, relo.Length)
			}
			if len(addr) != AddressLength {
				return &BinaryResponse{}, fmt.Errorf("address %s should be %d character long, not %d", addr, AddressLength, len(addr))
			}
			start := relo.Start * 2
			end := relo.Start*2 + AddressLength
			if bin[start+1] != '_' || bin[end-1] != '_' {
				return &BinaryResponse{}, fmt.Errorf("relocation dummy not found at %d in %s ", relo.Start, bin)
			}
			bin = bin[:start] + addr + bin[end:]
		}
	}

	return &BinaryResponse{
		Binary: bin,
		Error:  "",
	}, nil
}

//todo: Might also need to add in a map of library names to addrs
func RequestCompile(file string, optimize bool, libraries map[string]string) (*Response, error) {
	request, err := CreateRequest(file, libraries, optimize)
	if err != nil {
		return nil, err
	}

	/*for k, v := range request.Includes {
		log.WithFields(log.Fields{
			"k": k,
			"v": string(v.Script),
		}).Debug("check request loop")
	}*/

	resp := compile(request)

	PrintResponse(*resp, false)

	return resp, nil
}

// Compile takes a dir and some code, replaces all includes, compiles
func compile(req *Request) *Response {
	input := SolidityInput{Language: "Solidity", Sources: make(map[string]SolidityInputSource)}

	for k, v := range req.Includes {
		input.Sources[k] = SolidityInputSource{Content: string(v.Script)}
	}

	input.Settings.Optimizer.Enabled = req.Optimize
	input.Settings.OutputSelection.File.OutputType = []string{"abi", "evm.bytecode.linkReferences", "metadata", "bin", "devdoc"}
	input.Settings.Libraries = make(map[string]map[string]string)
	input.Settings.Libraries[""] = make(map[string]string)

	if req.Libraries != nil {
		for l, a := range req.Libraries {
			input.Settings.Libraries[""][l] = "0x" + a
		}
	}

	command, err := json.Marshal(input)
	if err != nil {
		return &Response{Error: err.Error()}
	}

	log.WithField("Command: ", command).Debug("Command Input")
	result, err := runSolidity(string(command))
	log.WithField("Command Result: ", result).Debug("Command Output")

	output := SolidityOutput{}
	err = json.Unmarshal([]byte(result), &output)
	if err != nil {
		return &Response{Error: err.Error()}
	}

	respItemArray := make([]ResponseItem, 0)

	for _, s := range output.Contracts {
		for contract, item := range s {
			respItem := ResponseItem{
				Objectname: objectName(contract),
				Binary:     item,
			}
			respItemArray = append(respItemArray, respItem)
		}
	}

	warnings := ""
	errors := ""
	for _, msg := range output.Errors {
		if msg.Type == "Warning" {
			warnings += msg.FormattedMessage
		} else {
			errors += msg.FormattedMessage
		}
	}

	for _, re := range respItemArray {
		log.WithFields(log.Fields{
			"name": re.Objectname,
			"bin":  re.Binary.Evm.Bytecode,
			"abi":  re.Binary.Abi,
		}).Debug("Response formulated")
	}

	return &Response{
		Objects: respItemArray,
		Warning: warnings,
		Error:   errors,
	}
}

func objectName(contract string) string {
	if contract == "" {
		return ""
	}
	parts := strings.Split(strings.TrimSpace(contract), ":")
	return parts[len(parts)-1]
}

func runSolidity(jsonCmd string) (string, error) {
	buf := bytes.NewBufferString(jsonCmd)
	shellCmd := exec.Command("solc", "--standard-json")
	shellCmd.Stdin = buf
	output, err := shellCmd.CombinedOutput()
	s := string(output)
	return s, err
}

func CreateRequest(file string, libraries map[string]string, optimize bool) (*Request, error) {
	var includes = make(map[string]*IncludedFiles)

	code, err := ioutil.ReadFile(file)
	if err != nil {
		return &Request{}, err
	}
	dir := path.Dir(file)
	//log.Debug("Before parsing includes =>\n\n%s", string(code))
	err = FindIncludes(code, dir, file, includes)
	if err != nil {
		return &Request{}, err
	}

	return CompilerRequest(file, includes, libraries, optimize), nil
}

func PrintResponse(resp Response, cli bool) {
	if resp.Error != "" {
		log.Warn(resp.Error)
	} else {
		for _, r := range resp.Objects {
			message := log.WithFields((log.Fields{
				"name": r.Objectname,
				"bin":  r.Binary.Evm.Bytecode,
				"abi":  string(r.Binary.Abi[:]),
				"link": string(r.Binary.Evm.Bytecode.LinkReferences[:]),
			}))
			if cli {
				message.Warn("Response")
			} else {
				message.Info("Response")
			}
		}
	}
}

// this handles all of our imports
type IncludedFiles struct {
	ObjectNames []string `json:"objectNames"` //objects in the file
	Script      []byte   `json:"script"`      //actual code
}

// Find all matches to the include regex
// Replace filenames with hashes
func FindIncludes(code []byte, dir, file string, includes map[string]*IncludedFiles) error {
	// find includes, load those as well
	regexPattern := `import (.+?)??("|')(.+?)("|')(as)?(.+)?;`
	var regExpression *regexp.Regexp
	var err error
	if regExpression, err = regexp.Compile(regexPattern); err != nil {
		return err
	}
	// Find all includes of included imports
	// do it recursively
	for _, s := range regExpression.FindAll(code, -1) {
		log.WithField("=>", string(s)).Debug("Include FindAll result")
		err := includeFinder(regExpression, s, dir, includes)
		if err != nil {
			log.Error("ERR!:", err)
		}
	}

	includeFile := &IncludedFiles{
		Script: code,
	}

	includes[file] = includeFile

	return nil
}

// read the included file, hash it; if we already have it, return include replacement
// if we don't, run replaceIncludes on it (recursive)
// modifies the "includes" map
func includeFinder(r *regexp.Regexp, originCode []byte, dir string, included map[string]*IncludedFiles) error {
	// regex look for strings that would match the import statement
	m := r.FindStringSubmatch(string(originCode))
	match := m[3]
	log.WithField("=>", match).Debug("Match")
	// load the file
	newFilePath := path.Join(dir, match)
	incl_code, err := ioutil.ReadFile(newFilePath)
	if err != nil {
		log.Errorln("failed to read include file", err)
		return fmt.Errorf("Failed to read include file: %s", err.Error())
	}

	if _, ok := included[match]; ok {
		return nil
	}

	// recursively replace the includes for this file
	this_dir := path.Dir(newFilePath)
	return FindIncludes(incl_code, this_dir, newFilePath, included)
}

func extractObjectNames(script []byte) ([]string, error) {
	regExpression, err := regexp.Compile("(contract|library) (.+?) (is)?(.+?)?({)")
	if err != nil {
		return nil, err
	}
	objectNamesList := regExpression.FindAllSubmatch(script, -1)
	var objects []string
	for _, objectNames := range objectNamesList {
		objects = append(objects, string(objectNames[2]))
	}
	return objects, nil
}
