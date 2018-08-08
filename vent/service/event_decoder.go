package service

import (
	"fmt"
	"strings"

	"github.com/hyperledger/burrow/rpc/rpcevents"
)

type EventLogDecoder func(*rpcevents.EventDataLog, map[string]string)

func DecodeEvent(header *rpcevents.EventHeader, log *rpcevents.EventDataLog, logDecoders map[string]EventLogDecoder) map[string]string {
	data := make(map[string]string)

	// decode header
	data["index"] = fmt.Sprintf("%v", header.GetIndex())
	data["height"] = fmt.Sprintf("%v", header.GetHeight())
	data["eventType"] = header.GetEventType()
	data["txHash"] = string(header.GetTxHash())

	// decode data log
	if len(log.Topics) > 1 {
		eventName := strings.Trim(string(log.Topics[1]), "\x00")

		data["eventName"] = eventName

		if logDecoder, ok := logDecoders[eventName]; ok {
			logDecoder(log, data)
		}
	}

	return data
}
