package plugins

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockLogger struct {
	logs []loggedMessage
}

type loggedMessage struct {
	level   logging.LogLevel
	message []any
}

func (m *MockLogger) Trace(v ...any) {
	m.logs = append(m.logs, loggedMessage{level: logging.Trace, message: v})
}

func (m *MockLogger) Debug(v ...any) {
	m.logs = append(m.logs, loggedMessage{level: logging.Debug, message: v})
}

func (m *MockLogger) Info(v ...any) {
	m.logs = append(m.logs, loggedMessage{level: logging.Info, message: v})
}

func (m *MockLogger) Warn(v ...any) {
	m.logs = append(m.logs, loggedMessage{level: logging.Warn, message: v})
}

func (m *MockLogger) Error(v ...any) {
	m.logs = append(m.logs, loggedMessage{level: logging.Error, message: v})
}

func (m *MockLogger) Fatal(v ...any) {
	m.logs = append(m.logs, loggedMessage{level: logging.Fatal, message: v})
}

func (m *MockLogger) AssertLogged(t *testing.T, level logging.LogLevel, message ...any) {
	assert.Equal(t, m.logs[len(m.logs)-1], loggedMessage{level, message})
}

func TestLogWrapper_Arguments(t *testing.T) {
	mock := &MockLogger{}
	logWrapper := &LogWrapper{logger: mock}

	// Without arguments

	logWrapper.Log(hclog.Info, "This is a default test!")
	mock.AssertLogged(t, logging.Info, "This is a default test!")

	// With string argument

	logWrapper.Log(hclog.Info, "This is a test!", "Argument-Key", "Argument-Value")
	mock.AssertLogged(t, logging.Info, "This is a test! (Argument-Key = \"Argument-Value\")")

	// With multiple string arguments

	logWrapper.Log(hclog.Info, "This is a test with multiple arguments!", "Argument-Key1", "Argument-Value1", "Argument-Key2", "Argument-Value2")
	mock.AssertLogged(t, logging.Info, "This is a test with multiple arguments! (Argument-Key1 = \"Argument-Value1\", Argument-Key2 = \"Argument-Value2\")")

	// With nil argument key

	logWrapper.Log(hclog.Info, "This is a nil test!", nil, "Some Value")
	mock.AssertLogged(t, logging.Info, "This is a nil test!")

	// With int argument key

	logWrapper.Log(hclog.Info, "This is an int test!", 1, "Some Value")
	mock.AssertLogged(t, logging.Info, "This is an int test!")

	// With int argument value

	logWrapper.Log(hclog.Info, "This is an int value test!", "Some key", 1)
	mock.AssertLogged(t, logging.Info, "This is an int value test! (Some key = 1)")

	// With []string argument value

	logWrapper.Log(hclog.Info, "This is a slice value test!", "Some key", []string{"First", "Second", "Third"})
	mock.AssertLogged(t, logging.Info, "This is a slice value test! (Some key = [\"First\", \"Second\", \"Third\"])")

	// Persistent argument

	logWrapper.With("This is persistent", "This is too").Log(hclog.Info, "This is a persistent argument test!")
	mock.AssertLogged(t, logging.Info, "This is a persistent argument test! (This is persistent = \"This is too\")")
	logWrapper.Log(hclog.Info, "The original logger stays the same!")
	mock.AssertLogged(t, logging.Info, "The original logger stays the same!")
}

func TestLogWrapper_Levels(t *testing.T) {
	mock := &MockLogger{}
	logWrapper := &LogWrapper{logger: mock}

	// Trace
	logWrapper.Trace("Trace")
	mock.AssertLogged(t, logging.Trace, "Trace")
	logWrapper.Log(hclog.Trace, "Trace")
	mock.AssertLogged(t, logging.Trace, "Trace")

	// Debug
	logWrapper.Debug("Debug")
	mock.AssertLogged(t, logging.Debug, "Debug")
	logWrapper.Log(hclog.Debug, "Debug")
	mock.AssertLogged(t, logging.Debug, "Debug")

	// Info
	logWrapper.Info("Info")
	mock.AssertLogged(t, logging.Info, "Info")
	logWrapper.Log(hclog.Info, "Info")
	mock.AssertLogged(t, logging.Info, "Info")

	// Warn
	logWrapper.Warn("Warn")
	mock.AssertLogged(t, logging.Warn, "Warn")
	logWrapper.Log(hclog.Warn, "Warn")
	mock.AssertLogged(t, logging.Warn, "Warn")

	// Error
	logWrapper.Error("Error")
	mock.AssertLogged(t, logging.Error, "Error")
	logWrapper.Log(hclog.Error, "Error")
	mock.AssertLogged(t, logging.Error, "Error")
}

func TestLogWrapper_Name(t *testing.T) {
	mock := &MockLogger{}
	logWrapper := &LogWrapper{logger: mock}

	logWrapper.Named("TestName").Log(hclog.Info, "Message with name")
	mock.AssertLogged(t, logging.Info, "[ TestName ] - Message with name")

	logWrapper.ResetNamed("TestName").Log(hclog.Info, "Message with name")
	mock.AssertLogged(t, logging.Info, "[ TestName ] - Message with name")

	logWrapper.Named("PrimaryName").Named("SecondaryName").Log(hclog.Info, "Message with name")
	mock.AssertLogged(t, logging.Info, "[ PrimaryName > SecondaryName ] - Message with name")

	logWrapper.Named("TestName").ResetNamed("NewTestName").Log(hclog.Info, "Message with reset name")
	mock.AssertLogged(t, logging.Info, "[ NewTestName ] - Message with reset name")
}

func TestLogWrapper_ImpliedArgs(t *testing.T) {
	logWrapper := &LogWrapper{}

	args := []string{"Test1", "Test2", "Test3", "Test4"}
	argumentLogWrapper := logWrapper.With("Test1", "Test2", "Test3", "Test4")

	assert.Equal(t, len(args), len(argumentLogWrapper.ImpliedArgs()))

	for i, arg := range args {
		assert.Equal(t, arg, argumentLogWrapper.ImpliedArgs()[i])
	}
}
