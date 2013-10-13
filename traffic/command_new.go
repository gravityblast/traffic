package main

import (
  "os"
  "fmt"
  "go/build"
  "io/ioutil"
  "path/filepath"
)

const NewProjectTemplateFolder = "new_project_templates"

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
  c.createFiles(projectPath)
}

func (c CommandNew) createProjectFolder(path string) {
  os.Mkdir(path, 0755)
}

func (c CommandNew) createFiles(targetPath string) {
  pkg, err := build.Import("github.com/pilu/traffic/traffic", "", build.FindOnly)
  if err != nil {
    logger.Fatal(err)
  }

  sourcePath := filepath.Join(pkg.Dir, NewProjectTemplateFolder)
  filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      logger.Fatal(err)
    }

    relPath, err := filepath.Rel(sourcePath, path)
    if relPath == "." {
      return nil
    }

    if err != nil {
      logger.Fatal(err)
    }

    newPath := filepath.Join(targetPath, relPath)
    c.createFile(path, newPath)

    return nil
  })
}

func (c CommandNew) createFile(sourcePath, targetPath string) {
  sourceStat, err := os.Stat(sourcePath)
  if err != nil {
    logger.Fatal(err)
  }

  targetStat, err := os.Stat(targetPath)
  if err == nil {
    fileType := "File"
    if targetStat.IsDir() {
      fileType = "Directory"
    }
    logger.Printf("%s already exists `%s`", fileType, targetPath)
    return
  }

  if sourceStat.IsDir() {
    logger.Printf("Creating folder `%s`", targetPath)
    os.Mkdir(targetPath, 0755)
    return
  }

  content, err := ioutil.ReadFile(sourcePath)
  if err != nil {
    logger.Fatal(err)
  }

  err = ioutil.WriteFile(targetPath, content, 0644)
  if err != nil {
    logger.Fatal(err)
  }

  logger.Printf("File created `%s`", targetPath)
}
