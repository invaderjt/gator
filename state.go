package main

import (
	"fmt"
	"log"
	"os"

	"github.com/invaderjt/blog-aggregator/internal/config"
)

func updateState() {
	s := &state{}
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Could not read config file")
	}
	s.Cfg = &cfg
	var cmds commands
	cmds.Commands = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)

	input := os.Args
	if len(input) < 2 {
		log.Fatalln("Not enough arguments")
	}
	var cmd command
	cmd.Name = input[1]
	cmd.Args = input[2:]
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s.Cfg)

}

type state struct {
	Cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	Commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	return c.Commands[cmd.Name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		log.Fatalln("Login requires username argument")
	}
	username := cmd.Args[0]
	err := s.Cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("%s set as username.", username)
	return nil
}
