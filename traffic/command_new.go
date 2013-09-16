package main

import (
  "fmt"
  "os"
  "path/filepath"
)

type CommandNew struct {}

func (c CommandNew) Execute(args []string) {
  if len(args) > 0 {
    c.generateProject(args[0])
  } else {
    c.Usage()
  }
}

func (c CommandNew) Usage() {
  fmt.Println("Usage:")
  fmt.Println("  traffic new APP_PATH")
}

func (c CommandNew) generateProject(projectName string) {
  currentPath, err := os.Getwd()
  if err != nil {
    logger.Fatal(err)
  }

  projectPath := filepath.Join(currentPath, projectName)

  if _, err := os.Stat(projectPath); err == nil {
    logger.Fatalf("Project `%s` already exists (%s).", projectName, projectPath)
  }

  logger.Printf("Creating project folder `%s`", projectPath)
  c.createProjectFolder(projectPath)
}

func (c CommandNew) createProjectFolder(path string) {
  os.Mkdir(path, 0755)
}


