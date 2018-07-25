package jobs

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/monax/bosmarmot/bos/abi"
	"github.com/monax/bosmarmot/bos/def"
	"github.com/monax/bosmarmot/bos/util"
	log "github.com/sirupsen/logrus"
)

func QueryContractJob(query *def.QueryContract, do *def.Packages) (string, []*def.Variable, error) {
	// Preprocess variables. We don't preprocess data as it is processed by ReadAbiFormulateCall
	query.Source, _ = util.PreProcess(query.Source, do)
	query.Destination, _ = util.PreProcess(query.Destination, do)
	query.ABI, _ = util.PreProcess(query.ABI, do)

	var queryDataArray []string
	var err error
	query.Function, queryDataArray, err = util.PreProcessInputData(query.Function, query.Data, do, false)
	if err != nil {
		return "", nil, err
	}

	// Get the packed data from the ABI functions
	var data string
	var packedBytes []byte
	if query.ABI == "" {
		packedBytes, err = abi.ReadAbiFormulateCall(query.Destination, query.Function, queryDataArray, do)
		data = hex.EncodeToString(packedBytes)
	} else {
		packedBytes, err = abi.ReadAbiFormulateCall(query.ABI, query.Function, queryDataArray, do)
		data = hex.EncodeToString(packedBytes)
	}
	if err != nil {
		var str, err = util.ABIErrorHandler(do, err, nil, query)
		return str, nil, err
	}

	// Call the client
	txe, err := do.QueryContract(def.QueryArg{
		Input:   query.Source,
		Address: query.Destination,
		Data:    data,
	})
	if err != nil {
		return "", nil, err
	}

	// Formally process the return
	log.WithField("res", txe.Result.Return).Debug("Decoding Raw Result")
	if query.ABI == "" {
		log.WithField("abi", query.Destination).Debug()
		query.Variables, err = abi.ReadAndDecodeContractReturn(query.Destination, query.Function, txe.Result.Return, do)
	} else {
		log.WithField("abi", query.ABI).Debug()
		query.Variables, err = abi.ReadAndDecodeContractReturn(query.ABI, query.Function, txe.Result.Return, do)
	}
	if err != nil {
		return "", nil, err
	}

	result2 := util.GetReturnValue(query.Variables)
	// Finalize
	if result2 != "" {
		log.WithField("=>", result2).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result2, query.Variables, nil
}

func QueryAccountJob(query *def.QueryAccount, do *def.Packages) (string, error) {
	// Preprocess variables
	query.Account, _ = util.PreProcess(query.Account, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Perform Query
	arg := fmt.Sprintf("%s:%s", query.Account, query.Field)
	log.WithField("=>", arg).Info("Querying Account")

	result, err := util.AccountsInfo(query.Account, query.Field, do)
	if err != nil {
		return "", err
	}

	// Result
	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func QueryNameJob(query *def.QueryName, do *def.Packages) (string, error) {
	// Preprocess variables
	query.Name, _ = util.PreProcess(query.Name, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	log.WithFields(log.Fields{
		"name":  query.Name,
		"field": query.Field,
	}).Info("Querying")
	result, err := util.NamesInfo(query.Name, query.Field, do)
	if err != nil {
		return "", err
	}

	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func QueryValsJob(query *def.QueryVals, do *def.Packages) (string, error) {
	var result string

	// Preprocess variables
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	log.WithField("=>", query.Field).Info("Querying Vals")
	result, err := util.ValidatorsInfo(query.Field, do)
	if err != nil {
		return "", err
	}

	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func AssertJob(assertion *def.Assert, do *def.Packages) (string, error) {
	var result string
	// Preprocess variables
	assertion.Key, _ = util.PreProcess(assertion.Key, do)
	assertion.Relation, _ = util.PreProcess(assertion.Relation, do)
	assertion.Value, _ = util.PreProcess(assertion.Value, do)

	// Switch on relation
	log.WithFields(log.Fields{
		"key":      assertion.Key,
		"relation": assertion.Relation,
		"value":    assertion.Value,
	}).Info("Assertion =>")

	switch assertion.Relation {
	case "==", "eq":
		/*log.Debug("Compare", strings.Compare(assertion.Key, assertion.Value))
		log.Debug("UTF8?: ", utf8.ValidString(assertion.Key))
		log.Debug("UTF8?: ", utf8.ValidString(assertion.Value))
		log.Debug("UTF8?: ", utf8.RuneCountInString(assertion.Key))
		log.Debug("UTF8?: ", utf8.RuneCountInString(assertion.Value))*/
		if assertion.Key == assertion.Value {
			return assertPass("==", assertion.Key, assertion.Value)
		} else {
			return assertFail("==", assertion.Key, assertion.Value)
		}
	case "!=", "ne":
		if assertion.Key != assertion.Value {
			return assertPass("!=", assertion.Key, assertion.Value)
		} else {
			return assertFail("!=", assertion.Key, assertion.Value)
		}
	case ">", "gt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k > v {
			return assertPass(">", assertion.Key, assertion.Value)
		} else {
			return assertFail(">", assertion.Key, assertion.Value)
		}
	case ">=", "ge":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k >= v {
			return assertPass(">=", assertion.Key, assertion.Value)
		} else {
			return assertFail(">=", assertion.Key, assertion.Value)
		}
	case "<", "lt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k < v {
			return assertPass("<", assertion.Key, assertion.Value)
		} else {
			return assertFail("<", assertion.Key, assertion.Value)
		}
	case "<=", "le":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k <= v {
			return assertPass("<=", assertion.Key, assertion.Value)
		} else {
			return assertFail("<=", assertion.Key, assertion.Value)
		}
	default:
		return "", fmt.Errorf("Error: Bad assert relation: \"%s\" is not a valid relation. See documentation for more information.", assertion.Relation)
	}

	return result, nil
}

func bulkConvert(key, value string) (int, int, error) {
	k, err := strconv.Atoi(key)
	if err != nil {
		return 0, 0, err
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, err
	}
	return k, v, nil
}

func assertPass(typ, key, val string) (string, error) {
	log.WithField("=>", fmt.Sprintf("%s %s %s", key, typ, val)).Warn("Assertion Succeeded")
	return "passed", nil
}

func assertFail(typ, key, val string) (string, error) {
	log.WithField("=>", fmt.Sprintf("%s %s %s", key, typ, val)).Warn("Assertion Failed")
	return "failed", fmt.Errorf("assertion failed")
}

func convFail() (string, error) {
	return "", fmt.Errorf("The Key of your assertion cannot be converted into an integer.\nFor string conversions please use the equal or not equal relations.")
}
