package service_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hiturria/bosmarmot/vent/config"
	"github.com/hiturria/bosmarmot/vent/logger"
	"github.com/hiturria/bosmarmot/vent/service"
	"github.com/hiturria/bosmarmot/vent/test"
	"github.com/hyperledger/burrow/core"
	"github.com/hyperledger/burrow/integration"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/rpc/rpctransact"
	"github.com/stretchr/testify/require"
)

var privateAccounts = integration.MakePrivateAccounts(5) // make keys
var genesisDoc = integration.TestGenesisDoc(privateAccounts)
var inputAccount = &rpctransact.InputAccount{Address: privateAccounts[0].Address().Bytes()}
var kern *core.Kernel

func testEventLogDecoder(log *rpcevents.EventDataLog, data map[string]string) {
	data["name"] = strings.Trim(string(log.Topics[2]), "\x00")
	data["description"] = strings.Trim(string(log.Topics[3]), "\x00")
}

func TestMain(m *testing.M) {
	returnValue := integration.TestWrapper(privateAccounts, genesisDoc, func(k *core.Kernel) int {
		kern = k
		return m.Run()
	})

	time.Sleep(3 * time.Second)
	os.Exit(returnValue)
}

func TestRun(t *testing.T) {
	tCli := test.NewTransactClient(t)
	create := test.CreateContract(t, tCli, inputAccount)

	// Here is how we can generate events
	name := "TestEvent1"
	description := "Description of TestEvent1"
	test.CallAddEvent(t, tCli, inputAccount, create.CallData.Callee, name, description)

	name = "TestEvent2"
	description = "Description of TestEvent2"
	test.CallAddEvent(t, tCli, inputAccount, create.CallData.Callee, name, description)

	name = "TestEvent3"
	description = "Description of TestEvent3"
	test.CallAddEvent(t, tCli, inputAccount, create.CallData.Callee, name, description)

	name = "TestEvent4"
	description = "Description of TestEvent4"
	test.CallAddEvent(t, tCli, inputAccount, create.CallData.Callee, name, description)

	// This is a workaround for off-by-one on latest bound fixed in burrow
	time.Sleep(time.Second * 2)

	// create test db
	db, closeDB := test.NewTestDB(t)
	defer closeDB()

	// Run consumer to listen to events
	cfg := config.DefaultFlags()
	cfg.DBSchema = db.Schema
	cfg.CfgFile = os.Getenv("GOPATH") + "/src/github.com/hiturria/bosmarmot/vent/test/sqlsol_example.json"

	log := logger.NewLogger(cfg.LogLevel)

	consumer := service.NewConsumer(cfg, log)
	consumer.AddEventLogDecoder("TEST_EVENTS", testEventLogDecoder)

	err := consumer.Run()
	require.NoError(t, err)

	// test data from two different block ids
	blockID := "2"
	eventName := "EventTest"
	eventData, err := db.GetBlock(blockID)
	require.NoError(t, err)
	require.Equal(t, "2", eventData.Block)
	require.Equal(t, 2, len(eventData.Tables))

	tblData := eventData.Tables[strings.ToLower(eventName)]
	require.Equal(t, 1, len(tblData))
	require.Equal(t, "0", tblData[0]["index"])
	require.Equal(t, "2", tblData[0]["height"])
	require.Equal(t, "\263+\307\246\336\224C\335\205i\011<\220\217\006\312\272\342\311\347", tblData[0]["txhash"])
	require.Equal(t, "LogEvent", tblData[0]["eventtype"])
	require.Equal(t, "TEST_EVENTS", tblData[0]["eventname"])
	require.Equal(t, "TestEvent1", tblData[0]["testname"])
	require.Equal(t, "Description of TestEvent1", tblData[0]["testdescription"])

	blockID = "5"
	eventName = "EventTest"
	eventData, err = db.GetBlock(blockID)
	require.NoError(t, err)
	require.Equal(t, "5", eventData.Block)
	require.Equal(t, 2, len(eventData.Tables))

	tblData = eventData.Tables[strings.ToLower(eventName)]
	require.Equal(t, 1, len(tblData))
	require.Equal(t, "0", tblData[0]["index"])
	require.Equal(t, "5", tblData[0]["height"])
	require.Equal(t, "5", tblData[0]["height"])
	require.Equal(t, "\213\366n\264)'\264\332ybb\3704\015\177\031\214\372\340\004", tblData[0]["txhash"])
	require.Equal(t, "LogEvent", tblData[0]["eventtype"])
	require.Equal(t, "TEST_EVENTS", tblData[0]["eventname"])
	require.Equal(t, "TestEvent4", tblData[0]["testname"])
	require.Equal(t, "Description of TestEvent4", tblData[0]["testdescription"])
}
