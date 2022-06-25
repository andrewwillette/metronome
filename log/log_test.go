package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configureLog(t *testing.T) {
	const logFile = "testfile.log"
	defer os.Remove(logFile)
	_, err := os.Stat(logFile)
	if !errors.Is(err, os.ErrNotExist) {
		t.Log(fmt.Sprintf("%v logfile should not exist", logFile))
		t.Fail()
	}
	ConfigureLog(logFile, true)
	_, err = os.Stat(logFile)
	if err != nil {
		t.Fail()
	}
	const toLog = "hello log"
	Lg(toLog)
	logFileBytes, _ := ioutil.ReadFile(logFile)
	assert.Contains(t, string(logFileBytes), toLog)
}
