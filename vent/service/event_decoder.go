package service

import (
	"fmt"
	"strings"

	"github.com/hyperledger/burrow/execution/exec"
)

type EventLogDecoder func(*exec.LogEvent, map[string]string)

// DecodeEvent decodes standard & non standard event data
// using provided logDecoders
func DecodeEvent(header *exec.Header, log *exec.LogEvent, logDecoders map[string]EventLogDecoder) map[string]string {
	data := make(map[string]string)

	// decode header
	data["index"] = fmt.Sprintf("%v", header.GetIndex())
	data["height"] = fmt.Sprintf("%v", header.GetHeight())
	data["eventType"] = header.GetEventType().String()
	data["txHash"] = string(header.TxHash)

	// decode data log
	if len(log.Topics) > 1 {
		eventName := strings.Trim(log.Topics[1].String(), "\x00")

		data["eventName"] = eventName

		if logDecoder, ok := logDecoders[eventName]; ok {
			logDecoder(log, data)
		}
	}

	return data
}
