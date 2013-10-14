package main

import (
  "io"
  "os"
  "fmt"
  "errors"
  "os/exec"
  "io/ioutil"
  "github.com/howeyc/fsnotify"
)

type CommandRun struct {}

func (c CommandRun) Execute(args []string) {
  buildPath := "tmp/traffic-build"

  if message, ok := c.build(buildPath); !ok {
    logger.Fatal(message)
  }

  runDone, message, ok := c.run(buildPath)

  if !ok {
    logger.Fatal(message)
  }

  doneWatching := make(chan bool)

  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    logger.Fatal(err)
  }

  go func() {
    for {
      select {
      case err := <-runDone:
        if err !=nil {
          logger.Print(err)
        }
      case ev := <-watcher.Event:
        logger.Println("event:", ev)
        runDone <- errors.New("file updated")
      case err := <-watcher.Error:
        logger.Println("error:", err)
      }
    }
  }()

  err = watcher.Watch(".")
  if err != nil {
    logger.Fatal(err)
  }

  <-doneWatching

  logger.Print("Command stopped")
}

func (c CommandRun) run(executablePath string) (chan error, string, bool) {
  done := make(chan error)

  cmd := exec.Command(executablePath)
  stderr, err := cmd.StderrPipe()
  if err != nil {
    logger.Fatal(err)
  }

  stdout, err := cmd.StdoutPipe()
  if err != nil {
    logger.Fatal(err)
  }

  if err := cmd.Start(); err != nil {
    logger.Fatal(err)
  }

  go func () {
    io.Copy(os.Stderr, stderr)
    io.Copy(os.Stdout, stdout)
  }()
  go func() {
    done <- cmd.Wait()
  }()

  return done, "", true
}


func (c CommandRun) build(outputPath string) (string, bool) {
  cmd := exec.Command("go", "build", "-o", outputPath)
  stderr, err := cmd.StderrPipe()
  if err != nil {
    logger.Fatal(err)
  }

  stdout, err := cmd.StdoutPipe()
  if err != nil {
    logger.Fatal(err)
  }

  err = cmd.Start();
  if err != nil {
    logger.Fatal(err)
  }

  io.Copy(os.Stdout, stdout)
  errBuf, _ := ioutil.ReadAll(stderr)

  err = cmd.Wait()
  if err != nil {
    return string(errBuf), false
  }

  return "", true
}

func (c CommandRun) Usage() {
  fmt.Println("Usage:")
  fmt.Println("  traffic run")
}
