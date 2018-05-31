package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Compiler struct {
	Config LangConfig
	Lang   string
}

// New Request object from script and map of include files
func (c *Compiler) CompilerRequest(file string,
	includes map[string]*IncludedFiles, libs string, optimize bool,
	hashFileReplacement map[string]string) *Request {
	if includes == nil {
		includes = make(map[string]*IncludedFiles)
	}
	return &Request{
		Language:        c.Lang,
		Includes:        includes,
		Libraries:       libs,
		Optimize:        optimize,
		FileReplacement: hashFileReplacement,
	}
}

// Compile request object
type Request struct {
	ScriptName      string                    `json:"name"`
	Language        string                    `json:"language"`
	Includes        map[string]*IncludedFiles `json:"includes"`  // our required files and metadata
	Libraries       string                    `json:"libraries"` // string of libName:LibAddr separated by comma
	Optimize        bool                      `json:"optimize"`  // run with optimize flag
	FileReplacement map[string]string         `json:"replacement"`
}

type BinaryRequest struct {
	BinaryFile string `json:"binary"`
	Libraries  string `json:"libraries"`
}

// this handles all of our imports
type IncludedFiles struct {
	ObjectNames []string `json:"objectNames"` //objects in the file
	Script      []byte   `json:"script"`      //actual code
}

const (
	SOLIDITY = "sol"
)

type LangConfig struct {
	CacheDir     string   `json:"cache"`
	IncludeRegex string   `json:"regex"`
	CompileCmd   []string `json:"cmd"`
}

// Fill in the filename and return the command line args
func (l LangConfig) Cmd(includes []string, libraries string, optimize bool) (args []string) {
	for _, s := range l.CompileCmd {
		if s == "_" {
			if optimize {
				args = append(args, "--optimize")
			}
			if libraries != "" {
				for _, l := range strings.Split(libraries, " ") {
					if len(l) > 0 {
						args = append(args, "--libraries")
						args = append(args, libraries)
					}
				}
			}
			args = append(args, includes...)
		} else {
			args = append(args, s)
		}
	}
	return
}

// todo: add indexes for where to find certain parts in submatches (quotes, filenames, etc.)
// Global variable mapping languages to their configs
var Languages = map[string]LangConfig{
	SOLIDITY: {
		CacheDir:     SolcScratchPath,
		IncludeRegex: `import (.+?)??("|')(.+?)("|')(as)?(.+)?;`,
		CompileCmd: []string{
			"solc",
			"--combined-json", "bin,abi",
			"_",
		},
	},
}

// individual contract items
type SolcItem struct {
	Bin string `json:"bin"`
	Abi string `json:"abi"`
}

// full solc response object
type SolcResponse struct {
	Contracts map[string]*SolcItem `mapstructure:"contracts" json:"contracts"`
	Version   string               `mapstructure:"version" json:"version"` // json encoded
}

func BlankSolcItem() *SolcItem {
	return &SolcItem{}
}

func BlankSolcResponse() *SolcResponse {
	return &SolcResponse{
		Version:   "",
		Contracts: make(map[string]*SolcItem),
	}
}

// Find all matches to the include regex
// Replace filenames with hashes
func (c *Compiler) ReplaceIncludes(code []byte, dir, file string,
	includes map[string]*IncludedFiles,
	hashFileReplacement map[string]string) ([]byte, error) {
	// find includes, load those as well
	regexPattern := c.IncludeRegex()
	var regExpression *regexp.Regexp
	var err error
	if regExpression, err = regexp.Compile(regexPattern); err != nil {
		return nil, err
	}
	OriginObjectNames, err := c.extractObjectNames(code)
	if err != nil {
		return nil, err
	}
	// replace all includes with hash of included imports
	// make sure to return hashes of includes so we can cache check them too
	// do it recursively
	code = regExpression.ReplaceAllFunc(code, func(s []byte) []byte {
		log.WithField("=>", string(s)).Debug("Include Replacer result")
		s, err := c.includeReplacer(regExpression, s, dir, includes, hashFileReplacement)
		if err != nil {
			log.Error("ERR!:", err)
		}
		return s
	})

	originHash := sha256.Sum256(code)
	origin := hex.EncodeToString(originHash[:])
	origin += "." + c.Lang

	includeFile := &IncludedFiles{
		ObjectNames: OriginObjectNames,
		Script:      code,
	}

	includes[origin] = includeFile
	hashFileReplacement[origin] = file

	return code, nil
}

// read the included file, hash it; if we already have it, return include replacement
// if we don't, run replaceIncludes on it (recursive)
// modifies the "includes" map
func (c *Compiler) includeReplacer(r *regexp.Regexp, originCode []byte, dir string, included map[string]*IncludedFiles, hashFileReplacement map[string]string) ([]byte, error) {
	// regex look for strings that would match the import statement
	m := r.FindStringSubmatch(string(originCode))
	match := m[3]
	log.WithField("=>", match).Debug("Match")
	// load the file
	newFilePath := path.Join(dir, match)
	incl_code, err := ioutil.ReadFile(newFilePath)
	if err != nil {
		log.Errorln("failed to read include file", err)
		return nil, fmt.Errorf("Failed to read include file: %s", err.Error())
	}

	// take hash before replacing includes to see if we've already parsed this file
	hash := sha256.Sum256(incl_code)
	includeHash := hex.EncodeToString(hash[:])
	log.WithField("=>", includeHash).Debug("Included Code's Hash")
	if _, ok := included[includeHash]; ok {
		//then replace
		fullReplacement := strings.SplitAfter(m[0], m[2])
		fullReplacement[1] = includeHash + "." + c.Lang + "\""
		ret := strings.Join(fullReplacement, "")
		return []byte(ret), nil
	}

	// recursively replace the includes for this file
	this_dir := path.Dir(newFilePath)
	incl_code, err = c.ReplaceIncludes(incl_code, this_dir, newFilePath, included, hashFileReplacement)
	if err != nil {
		return nil, err
	}

	// compute hash
	hash = sha256.Sum256(incl_code)
	h := hex.EncodeToString(hash[:])

	//Starting with full regex string,
	//Split strings from the quotation mark and then,
	//assuming 3 array cells, replace the middle one.
	fullReplacement := strings.SplitAfter(m[0], m[2])
	fullReplacement[1] = h + "." + c.Lang + m[4]
	ret := []byte(strings.Join(fullReplacement, ""))
	return ret, nil
}

// Return the regex string to match include statements
func (c *Compiler) IncludeRegex() string {
	return c.Config.IncludeRegex
}

func (c *Compiler) extractObjectNames(script []byte) ([]string, error) {
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

var (
	BosRoot         = ResolveBosRoot()
	SolcScratchPath = filepath.Join(BosRoot, "sol")
)

func ResolveBosRoot() string {
	var bos string
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		bos = filepath.Join(home, ".bos")
	} else {
		bos = filepath.Join(os.Getenv("HOME"), ".bos")
	}
	return bos
}

// InitScratchDir creates an Monax directory hierarchy under BosRoot dir.
func InitScratchDir() (err error) {
	for _, d := range []string{
		BosRoot,
		SolcScratchPath,
	} {
		if _, err := os.Stat(d); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(d, 0777); err != nil {
					return err
				}
			}
		}
	}
	return
}
