package main

import (
	"core.bank/datamover/cmd"
	"core.bank/datamover/log"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("something wrong, received from panic: %v", r)
		}
	}()

	cmd.Execute()

}
