package jobs

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/hyperledger/burrow/txs/payload"
	"github.com/monax/bosmarmot/bos/def"
	"github.com/monax/bosmarmot/bos/util"
	log "github.com/sirupsen/logrus"
)

func SendJob(send *def.Send, do *def.Packages) (string, error) {
	// Process Variables
	send.Source, _ = util.PreProcess(send.Source, do)
	send.Destination, _ = util.PreProcess(send.Destination, do)
	send.Amount, _ = util.PreProcess(send.Amount, do)

	// Use Default
	send.Source = useDefault(send.Source, do.Package.Account)

	// Formulate tx
	log.WithFields(log.Fields{
		"source":      send.Source,
		"destination": send.Destination,
		"amount":      send.Amount,
	}).Info("Sending Transaction")

	tx, err := do.Send(def.SendArg{
		Input:    send.Source,
		Output:   send.Destination,
		Amount:   send.Amount,
		Sequence: send.Sequence,
	})
	if err != nil {
		return "", util.MintChainErrorHandler(do, err)
	}

	// Sign, broadcast, display
	return txFinalize(do, tx)
}

func RegisterNameJob(name *def.RegisterName, do *def.Packages) (string, error) {
	// Process Variables
	name.DataFile, _ = util.PreProcess(name.DataFile, do)

	// If a data file is given it should be in csv format and
	// it will be read first. Once the file is parsed and sent
	// to the chain then a single nameRegTx will be sent if that
	// has been populated.
	if name.DataFile != "" {
		// open the file and use a reader
		fileReader, err := os.Open(name.DataFile)
		if err != nil {
			return "", err
		}

		defer fileReader.Close()
		r := csv.NewReader(fileReader)

		// loop through the records
		for {
			// Read the record
			record, err := r.Read()

			// Catch the errors
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", err
			}

			// Sink the Amount into the third slot in the record if
			// it doesn't exist
			if len(record) <= 2 {
				record = append(record, name.Amount)
			}

			// Send an individual Tx for the record
			// [TODO]: move these to async using goroutines?
			r, err := registerNameTx(&def.RegisterName{
				Source:   name.Source,
				Name:     record[0],
				Data:     record[1],
				Amount:   record[2],
				Fee:      name.Fee,
				Sequence: name.Sequence,
			}, do)

			if err != nil {
				return "", err
			}

			n := fmt.Sprintf("%s:%s", record[0], record[1])

			// TODO: write smarter
			if err = WriteJobResultCSV(n, r); err != nil {
				return "", err
			}
		}
	}

	// If the data field is populated then there is a single
	// nameRegTx to send. So do that *now*.
	if name.Data != "" {
		return registerNameTx(name, do)
	} else {
		return "data_file_parsed", nil
	}
}

// Runs an individual nametx.
func registerNameTx(name *def.RegisterName, do *def.Packages) (string, error) {
	// Process Variables
	name.Source, _ = util.PreProcess(name.Source, do)
	name.Name, _ = util.PreProcess(name.Name, do)
	name.Data, _ = util.PreProcess(name.Data, do)
	name.Amount, _ = util.PreProcess(name.Amount, do)
	name.Fee, _ = util.PreProcess(name.Fee, do)

	// Set Defaults
	name.Source = useDefault(name.Source, do.Package.Account)
	name.Fee = useDefault(name.Fee, do.DefaultFee)
	name.Amount = useDefault(name.Amount, do.DefaultAmount)

	// Formulate tx
	log.WithFields(log.Fields{
		"name":   name.Name,
		"data":   name.Data,
		"amount": name.Amount,
	}).Info("NameReg Transaction")

	//tx, err := rpc.Name(monaxNodeClient, monaxKeyClient, name.Source, name.Amount, name.Nonce, name.Fee, name.Name, name.Data)
	tx, err := do.Name(def.NameArg{
		Input:    name.Source,
		Sequence: name.Sequence,
		Name:     name.Name,
		Amount:   name.Amount,
		Data:     name.Data,
		Fee:      name.Fee,
	})
	if err != nil {
		return "", util.MintChainErrorHandler(do, err)
	}
	// Sign, broadcast, display
	return txFinalize(do, tx)
}

func PermissionJob(perm *def.Permission, do *def.Packages) (string, error) {
	// Process Variables
	perm.Source, _ = util.PreProcess(perm.Source, do)
	perm.Action, _ = util.PreProcess(perm.Action, do)
	perm.PermissionFlag, _ = util.PreProcess(perm.PermissionFlag, do)
	perm.Value, _ = util.PreProcess(perm.Value, do)
	perm.Target, _ = util.PreProcess(perm.Target, do)
	perm.Role, _ = util.PreProcess(perm.Role, do)

	// Set defaults
	perm.Source = useDefault(perm.Source, do.Package.Account)

	log.Debug("Target: ", perm.Target)
	log.Debug("Marmots Deny: ", perm.Role)
	log.Debug("Action: ", perm.Action)
	// Populate the transaction appropriately

	// Formulate tx
	//tx, err := rpc.Permissions(monaxNodeClient, monaxKeyClient, perm.Source, perm.Nonce, perm.Action,
	//	perm.Target, perm.PermissionFlag, perm.Role, perm.Value)
	tx, err := do.Permissions(def.PermArg{
		Input:      perm.Source,
		Sequence:   perm.Sequence,
		Action:     perm.Action,
		Target:     perm.Target,
		Permission: perm.PermissionFlag,
		Role:       perm.Role,
		Value:      perm.Value,
	})
	if err != nil {
		return "", util.MintChainErrorHandler(do, err)
	}

	log.Debug("What are the args returned in transaction: ", tx.PermArgs)

	// Sign, broadcast, display
	return txFinalize(do, tx)
}

func txFinalize(do *def.Packages, tx payload.Payload) (string, error) {
	txe, err := do.SignAndBroadcast(tx)
	if err != nil {
		return "", util.MintChainErrorHandler(do, err)
	}

	util.ReadTxSignAndBroadcast(txe, err)
	if err != nil {
		return "", err
	}

	return txe.Receipt.TxHash.String(), nil
}

func useDefault(thisOne, defaultOne string) string {
	if thisOne == "" {
		return defaultOne
	}
	return thisOne
}
