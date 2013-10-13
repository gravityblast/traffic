package main

import (
  "os"
  "fmt"
  "log"
)

type Command interface {
  Execute(args []string)
}

var commands map[string]Command
var logger *log.Logger

func init() {
  logger = log.New(os.Stderr, "", 0)
  commands = map[string]Command{
    "new": &CommandNew{},
  }
}

func usage() {
  fmt.Println("Available commands:")
  for name, _ := range commands {
    fmt.Printf("  - %s\n", name)
  }
}

func main() {
  var command string
  args := os.Args

  if len(args) > 1 {
    command = args[1]
  }

  if commands[command] != nil {
    c := commands[command]
    c.Execute(args[2:len(args)])
  } else {
    fmt.Printf("Unknown command `%s`\n", command)
    usage()
  }
}
