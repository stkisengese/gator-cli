package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/stkisengese/gator-cli/internal/config"
	"github.com/stkisengese/gator-cli/internal/database"
)

func main() {
	cfg := loadConfig()
	db := database.Connect(cfg.DBURL)
	defer db.Close()

	appState := &state{
		db:  database.New(db),
		cfg: cfg,
	}
	cmds := newCommands()
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

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

func loadConfig() *config.Config {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return cfg
}
