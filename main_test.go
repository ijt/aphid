package main

import (
	"testing"
)

func Test_parseConfig_returnsErrorOnNotFound(t *testing.T) {
	msg := `
:error: Not Found
`
	_, err := parseConfig([]byte(msg))
	if err == nil {
		t.Error("No error returned for a not-found page.")
	}
}

