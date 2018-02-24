package main

import (
	"os"
	"testing"
)

func init() {
	params := new(Params)
	params.Set(os.Args)
}

func TestParams_Diff(t *testing.T) {

}