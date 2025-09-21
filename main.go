package main

import (
	"os"

	"github.com/jonassiebler/chatmate/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
