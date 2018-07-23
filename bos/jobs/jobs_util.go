package jobs

import (
	"github.com/monax/bosmarmot/bos/def"
	"github.com/monax/bosmarmot/bos/util"
	log "github.com/sirupsen/logrus"
)

func SetAccountJob(account *def.Account, do *def.Packages) (string, error) {
	var result string

	// Preprocess
	account.Address, _ = util.PreProcess(account.Address, do)

	// Set the Account in the Package & Announce
	do.Package.Account = account.Address
	log.WithField("=>", do.Package.Account).Info("Setting Account")

	// Set result and return
	result = account.Address
	return result, nil
}

func SetValJob(set *def.SetJob, do *def.Packages) (string, error) {
	var result string
	set.Value, _ = util.PreProcess(set.Value, do)
	log.WithField("=>", set.Value).Info("Setting Variable")
	result = set.Value
	return result, nil
}
