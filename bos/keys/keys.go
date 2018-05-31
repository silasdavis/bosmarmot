package keys

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hyperledger/burrow/keys"
	"github.com/hyperledger/burrow/logging"
	log "github.com/sirupsen/logrus"
)

type LocalKeyClient struct {
	keys.KeyClient
}

const DefaultKeysHost = "localhost"
const DefaultKeysPort = "4767"

var keysTimeout = 5 * time.Second

func DefaultKeysURL() string {
	return fmt.Sprintf("http://%s:%s", DefaultKeysHost, DefaultKeysPort)
}

func NewKeyClient(keysUrl string) *LocalKeyClient {
	return &LocalKeyClient{
		keys.NewKeyClient(keysUrl, logging.NewNoopLogger()),
	}
}

// Returns an initialized key client to a docker container
// running the keys server
// Adding the Ip address is optional and should only be used
// for passing data
func InitKeyClient(keysUrl string) (*LocalKeyClient, error) {
	aliveCh := make(chan struct{})
	localKeyClient := NewKeyClient(keysUrl)
	err := localKeyClient.HealthCheck()

	go func() {
		for err != nil {
			err = localKeyClient.HealthCheck()
		}
		aliveCh <- struct{}{}
	}()
	select {
	case <-time.After(keysTimeout):
		return nil, fmt.Errorf("keys instance did not become responsive after %s: %v", keysTimeout, err)
	case <-aliveCh:
		return localKeyClient, nil
	}
}

// Keyclient returns a list of keys that it is aware of.
// params:
// host - search for keys on the host
// container - search for keys on the container
// quiet - don't print output, just return the list you find
func (keys *LocalKeyClient) ListKeys(keysPath string, quiet bool) ([]string, error) {
	var result []string
	addrs, err := ioutil.ReadDir(keysPath)
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		result = append(result, addr.Name())
	}
	if !quiet {
		if len(addrs) == 0 {
			log.Warn("No keys found on host")
		} else {
			// First key.
			log.WithField("=>", result[0]).Warn("The keys on your host kind marmot")
			// Subsequent keys.
			if len(result) > 1 {
				for _, addr := range result[1:] {
					log.WithField("=>", addr).Warn()
				}
			}
		}
	}

	return result, nil
}
