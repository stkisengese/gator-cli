package main

import (
	"log"
	"os"

	"github.com/stkisengese/gator-cli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	appState := &state{cfg}
	cmds := newCommands()
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(appState, cmd); err != nil {
		log.Fatalf("Error running command: %v", err)
	}

}
