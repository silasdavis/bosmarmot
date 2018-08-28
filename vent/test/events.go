package test

import (
	"context"
	"testing"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc/rpctransact"
	"github.com/hyperledger/burrow/txs/payload"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func NewTransactClient(t testing.TB, listenAddress string) rpctransact.TransactClient {
	t.Helper()

	conn, err := grpc.Dial(listenAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return rpctransact.NewTransactClient(conn)
}

func CreateContract(t testing.TB, cli rpctransact.TransactClient, inputAddress crypto.Address) *exec.TxExecution {
	t.Helper()

	txe, err := cli.CallTxSync(context.Background(), &payload.CallTx{
		Input: &payload.TxInput{
			Address: inputAddress,
			Amount:  2,
		},
		Address:  nil,
		Data:     BytecodeEventsTest,
		Fee:      2,
		GasLimit: 10000,
	})
	require.NoError(t, err)

	return txe
}

func CallAddEvent(t testing.TB, cli rpctransact.TransactClient, inputAddress, contractAddress crypto.Address, name, description string) *exec.TxExecution {
	t.Helper()

	txe, err := cli.CallTxSync(context.Background(), &payload.CallTx{
		Input: &payload.TxInput{
			Address: inputAddress,
			Amount:  2,
		},
		Address:  &contractAddress,
		Data:     abiPackAddEvent(name, description),
		Fee:      2,
		GasLimit: 1000000,
	})
	require.NoError(t, err)

	if txe.Exception != nil {
		t.Fatalf("call should not generate exception but returned: %v", txe.Exception.Error())
	}

	return txe
}

// We will abstract away the below into some library calls in Burrow - but this gives us just enough for now
func abiPackAddEvent(name, description string) []byte {
	// From AddEvents.sol
	// Pack the args
	namePacked := abiStringPack(name)
	descriptionPacked := abiStringPack(description)
	// The 'heads' are the offsets of each argument in the packed bytes, it's size is one word per argument
	headSize := uint64(64)
	descriptionOffset := headSize + uint64(len(namePacked))
	// Write function selector
	id := abi.GetFunctionID("addEvent(string,string)")
	// Write heads (where does 'name' start)
	bs := append(id[:], wordN(headSize)...)
	// Where does 'description start'
	bs = append(bs, wordN(descriptionOffset)...)
	// Write tails
	bs = append(bs, namePacked...)
	bs = append(bs, descriptionPacked...)
	return bs
}

func abiStringPack(str string) []byte {
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

func wordN(i uint64) []byte {
	word := make([]byte, 32)
	binary.PutUint64BE(word[24:32], i)
	return word
}
