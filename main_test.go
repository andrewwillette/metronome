package main

import (
	"fmt"
	"testing"
)

func Test_getHeader(t *testing.T) {
    if GetApplicationTitle() != "Metronome" {
        t.Error("header incorrect")
    }
}

func Test_fail(t *testing.T) {
    fmt.Println("we good")
    // t.Error("failed")
    // t.Fail()
}
