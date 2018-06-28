// +build integration

package integration

import (
	"context"
	"testing"

	"io"

	"strings"
	"time"

	"github.com/hyperledger/burrow/execution/events/pbevents"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/execution/pbtransactor"
	"github.com/hyperledger/burrow/rpc"
	"github.com/monax/bosmarmot/.gopath_burrow/src/github.com/hyperledger/burrow/binary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestGenerateEvents(t *testing.T) {
	tCli := NewTransactorClient(t)
	create := CreateContract(t, tCli, inputAccount)
	// Here is how we can generate an event
	name := "Momentus"
	description := "eventus"
	CallAddEvent(t, tCli, inputAccount, create.CallData.Callee, name, description)
	// This is a workaround for off-by-one on latest bound fixed in burrow
	time.Sleep(time.Second * 2)

	// Now we can query the execution events service
	eCli := NewExecutionEventsClient(t)
	request := &pbevents.GetEventsRequest{
		Query:      "EventType = 'LogEvent'",
		BlockRange: pbevents.NewBlockRange(pbevents.AbsoluteBound(0), pbevents.LatestBound()),
	}
	evs, err := eCli.GetEvents(context.Background(), request)
	require.NoError(t, err)

	// Grab the events
	var events []*pbevents.ExecutionEvent
	for {
		resp, err := evs.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}
		events = append(events, resp.Events...)
	}
	t.Log(events)

	// Here is how your consumer can get some basic data
	require.Len(t, events, 1)
	log := events[0].GetEventDataLog()
	assert.Equal(t, "TEST_EVENTS", TrimTopic(log.Topics[1]))
	assert.Equal(t, name, TrimTopic(log.Topics[2]))
	assert.Equal(t, description, TrimTopic(log.Topics[3]))
}

func CreateContract(t testing.TB, cli pbtransactor.TransactorClient,
	inputAccount *pbtransactor.InputAccount) *pbevents.EventDataCall {

	create, err := cli.TransactAndHold(context.Background(), &pbtransactor.TransactParam{
		InputAccount: inputAccount,
		Address:      nil,
		Data:         BytecodeEventsTest,
		Fee:          2,
		GasLimit:     10000,
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(0), create.StackDepth)
	return create
}

func CallAddEvent(t testing.TB, cli pbtransactor.TransactorClient, inputAccount *pbtransactor.InputAccount,
	contractAddress []byte, name, description string) (call *pbevents.EventDataCall) {

	call, err := cli.TransactAndHold(context.Background(), &pbtransactor.TransactParam{
		InputAccount: inputAccount,
		Address:      contractAddress,
		Data:         ABIPackAddEvent(name, description),
		Fee:          2,
		GasLimit:     10000,
	})
	require.NoError(t, err)
	if call.Exception != nil {
		t.Fatalf("call should not generate exception but returned: %v", call.Exception.Error())
	}
	return call
}

func NewTransactorClient(t testing.TB) pbtransactor.TransactorClient {
	conn, err := grpc.Dial(rpc.DefaultGRPCConfig().ListenAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pbtransactor.NewTransactorClient(conn)
}

func NewExecutionEventsClient(t testing.TB) pbevents.ExecutionEventsClient {
	conn, err := grpc.Dial(rpc.DefaultGRPCConfig().ListenAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pbevents.NewExecutionEventsClient(conn)
}

func TrimTopic(topic []byte) string {
	return strings.Trim(string(topic), "\x00")
}

// We will abstract away the below into some library calls in Burrow - but this gives us just enough for now
func ABIPackAddEvent(name, description string) []byte {
	// From AddEvents.sol
	// Pack the args
	namePacked := ABIStringPack(name)
	descriptionPacked := ABIStringPack(description)
	// The 'heads' are the offsets of each argument in the packed bytes, it's size is one word per argument
	headSize := uint64(64)
	descriptionOffset := headSize + uint64(len(namePacked))
	// Write function selector
	id := abi.FunctionID("addEvent(string,string)")
	// Write heads (where does 'name' start)
	bs := append(id[:], wordN(headSize)...)
	// Where does 'description start'
	bs = append(bs, wordN(descriptionOffset)...)
	// Write tails
	bs = append(bs, namePacked...)
	bs = append(bs, descriptionPacked...)
	return bs
}

func ABIStringPack(str string) []byte {
	// How many word256 needed for string Ceiling(len(str)/32)
	n := (len(str) + 31) / 32
	// Make the packed bytes of the correct multiple plus space for uint256 for length prefix
	bs := make([]byte, (n+1)*32)
	// Put the length in the first word256 by writing uint64 into last 8 bytes of first word
	copy(bs[:32], wordN(uint64(len(str))))
	// copy in the string
	copy(bs[32:], []byte(str))
	return bs
}

func TestABIStringPack(t *testing.T) {
	str := "01234567012345670123456701234567"
	assert.Equal(t, append(wordN(32), []byte(str)...), ABIStringPack(str))
	str = "0123456701234567012345670123456"
	assert.Equal(t, append(append(wordN(31), []byte(str)...), 0), ABIStringPack(str))
	str = "012345670123456701234567012345670"
	assert.Equal(t, append(append(wordN(33), []byte(str)...), wordN(0)[:31]...), ABIStringPack(str))
}

func wordN(i uint64) []byte {
	word := make([]byte, 32)
	binary.PutUint64BE(word[24:32], i)
	return word
}
