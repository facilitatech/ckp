package main

import (
	"os"
)

func init() {
	params := new(Params)
	params.Set(os.Args)
}
