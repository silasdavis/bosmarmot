package test

import (
	"context"
	"strings"
	"testing"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/rpc/rpctransact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func CreateContract(t testing.TB, cli rpctransact.TransactClient, inputAccount *rpctransact.InputAccount) *rpcevents.EventDataCall {
	t.Helper()

	create, err := cli.TransactAndHold(context.Background(), &rpctransact.TransactParam{
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

func CallAddEvent(t testing.TB, cli rpctransact.TransactClient, inputAccount *rpctransact.InputAccount, contractAddress []byte, name, description string) (call *rpcevents.EventDataCall) {
	t.Helper()

	call, err := cli.TransactAndHold(context.Background(), &rpctransact.TransactParam{
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

func NewTransactClient(t testing.TB) rpctransact.TransactClient {
	t.Helper()

	conn, err := grpc.Dial(rpc.DefaultGRPCConfig().ListenAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return rpctransact.NewTransactClient(conn)
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
