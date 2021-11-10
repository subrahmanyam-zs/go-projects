package log

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestEntryFromInputs(t *testing.T) {
	testcases := []struct {
		format        string
		args          []interface{}
		expectedEntry *entry
	}{
		{"", []interface{}{}, &entry{Data: map[string]interface{}{}, Message: ""}},
		{"", []interface{}{"hello logging"}, &entry{Data: map[string]interface{}{},
			Message: "hello logging"}},
		{"", []interface{}{"hello", "logging"}, &entry{Data: map[string]interface{}{},
			Message: "hello logging"}},
		{"hello %v %v", []interface{}{"logging", "yoyo"}, &entry{Data: map[string]interface{}{},
			Message: "hello logging yoyo"}},
		{"", []interface{}{"hello", "%vlogging"}, &entry{Data: map[string]interface{}{},
			Message: "hello %vlogging"}},
	}

	for i, v := range testcases {
		e, _, _ := entryFromInputs(v.format, v.args...)
		if !reflect.DeepEqual(e.Data, v.expectedEntry.Data) || !reflect.DeepEqual(e.Message, v.expectedEntry.Message) {
			t.Errorf("[TESTCASE%d]Failed.Expected Data:%v Message %v\nGot Data:%v Message %v\n", i+1,
				v.expectedEntry.Data, v.expectedEntry.Message, e.Data, e.Message)
		}
	}
}

func TestEntryFromStringForJSON(t *testing.T) {
	testcases := []struct {
		args       interface{}
		exp        entry
		expPerfLog bool
	}{
		{map[string]interface{}{"message": "hello", "correlationId": "test", "responseCode": 200},
			entry{
				Data:          map[string]interface{}{},
				CorrelationID: "test", Message: map[string]interface{}{"message": "hello", "correlationId": "test", "responseCode": 200},
			}, false,
		},
		{map[string]interface{}{"message": "hello", "correlationId": "test", "responseCode": 200, "method": "GET", "uri": "/temp", "duration": 5},
			entry{
				Data:          map[string]interface{}{"responseCode": 200.00, "method": "GET", "uri": "/temp", "duration": 5.00},
				CorrelationID: "test", Message: "hello",
			}, true,
		},
	}

	for i, tc := range testcases {
		e, _, isPerfLog := entryFromInputs("", tc.args)
		if !reflect.DeepEqual(e.Data, tc.exp.Data) {
			t.Errorf("testcase [%v] failed expected data %v\tgot%v\n", i, tc.exp.Data, e.Data)
		}

		if !reflect.DeepEqual(e.Message, tc.exp.Message) {
			t.Errorf("testcase [%v] expected message %v\tgot%v\n", i, tc.exp.Message, e.Message)
		}

		if !reflect.DeepEqual(e.CorrelationID, tc.exp.CorrelationID) {
			t.Errorf("testcase [%v] expected correlation id %v\tgot%v\n", i, tc.exp.CorrelationID, e.CorrelationID)
		}

		if isPerfLog != tc.expPerfLog {
			t.Errorf("testcase [%v] expected performanceLog: %v\tgot performanceLog: %v\n", i, tc.expPerfLog, isPerfLog)
		}
	}
}

func TestEntry_TerminalOutput(t *testing.T) {
	now := time.Now()
	formattedNow := now.Format("15:04:05")
	appInfoData := &sync.Map{}
	appInfoData.Store("a", "b")

	appData := make(map[string]interface{})

	appData["a"] = "b"

	testcases := []struct {
		input  entry
		output string
	}{
		// fatal error checking if msg and level is logged
		{entry{Level: Fatal, Message: "fatal error", Time: now},
			"FATA\u001B[0m [" + formattedNow + "]  fatal error"},
		// errorMessage
		{entry{Level: Fatal, Message: "fatal error", Time: now, Data: map[string]interface{}{"errorMessage": "error"}},
			fmt.Sprintf("\x1b[31mFATA\x1b[0m [%s]  error fatal error\n                \x1b[37m (Memory: <nil> GoRoutines: <nil>) \x1b[0m\n",
				formattedNow)},
		// DataQuery
		{entry{Level: Info, Message: "query field exists", Data: map[string]interface{}{"query": "query"}},
			"\x1b[36mINFO\x1b[0m [00:00:00]  query field exists query\n                \x1b[37m (Memory: <nil> GoRoutines: <nil>) \x1b[0m\n"},
		// correlationId
		{entry{Level: Info, CorrelationID: "test", Message: "hello"}, fmt.Sprintf(
			"INFO\u001B[0m [00:00:00]  hello\n%15s: %s", "CorrelationId", "test")},
		// data with message
		{entry{Level: Warn, Message: "hello", Data: map[string]interface{}{"name": "ZopSmart"}},
			"WARN\u001B[0m [00:00:00]  hello"},
		// statusCode
		{entry{Level: Warn, Message: "hello", Data: map[string]interface{}{"name": "ZopSmart", "responseCode": 200}},
			"WARN\u001B[0m [00:00:00]  hello"},
		// test data
		{entry{Level: Debug, Data: map[string]interface{}{"method": "get", "duration": 10000.0, "uri": "i", "datastore": "cql"}},
			fmt.Sprintf("\x1b[37mDEBU\x1b[0m [00:00:00] %s - %.2fms\n                \x1b[37m (Memory: <nil> GoRoutines: <nil>) \x1b[%vm\n",
				"cql", 10.0, 0)},
		// app data
		{entry{Level: Info, App: appInfo{Data: appData}, Message: "test"}, fmt.Sprintf(
			"INFO\u001B[0m [00:00:00]  test\n%15s: %v", "a", "b")},
	}

	for i := range testcases {
		output := testcases[i].input.TerminalOutput()
		if !strings.Contains(output, testcases[i].output) {
			t.Errorf("[TESTCASE%d]got %v\texpected %v\n", i+1, output, testcases[i].output)
		}
	}
}

func Test_AppDataWithoutPersistence(t *testing.T) {
	b := new(bytes.Buffer)
	logger := NewMockLogger(b)
	testCases := []struct {
		format          string
		values          []interface{}
		appData         map[string]interface{}
		expectedLog     string
		expectedAppData string
	}{

		{"Percent: %v %%", []interface{}{"5", map[string]interface{}{"key1": "value1"}}, nil,
			"Percent: 5 %", `"data":{"key1":"value1"}`},
		{"Hello %v", []interface{}{"gofr", map[string]interface{}{
			"key1": "value1"}, data{"Gofr", 1}}, map[string]interface{}{
			"key2": "value2"}, "Hello gofr", `"data":{"age":1,"key1":"value1","key2":"value2","name":"Gofr"}`},
		{"Hello %v", []interface{}{"gofr", map[string]interface{}{
			"key1": "value1"}, &data{"Rohan", 25}}, map[string]interface{}{
			"key2": "value2"}, "Hello gofr", `"data":{"age":25,"key1":"value1","key2":"value2","name":"Rohan"}`},
		{"Hello %v", []interface{}{"gofr", map[string]interface{}{"key1": "value1"}, nil}, map[string]interface{}{
			"key2": "value2"}, "Hello gofr", `"data":{"key1":"value1","key2":"value2"}`},
	}

	for i, v := range testCases {
		b.Reset()

		for key, val := range v.appData {
			logger.AddData(key, val)
		}

		logger.Infof(v.format, v.values...)

		if !strings.Contains(b.String(), v.expectedLog) {
			t.Errorf("TESTACASE[%v] Failed, expected: %v, got: %v", i, v.expectedLog, b.String())
		}

		if !strings.Contains(b.String(), v.expectedAppData) {
			t.Errorf("TESTACASE[%v] Failed, expected: %v, got: %v", i, v.expectedAppData, b.String())
		}
	}
}
