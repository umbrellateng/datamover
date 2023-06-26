package main

import (
	"log"
	"os/exec"
)

func main() {

	cmd := exec.Command("./scripts/from.sh")

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(output))

}
