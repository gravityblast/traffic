package traffic

import (
  "os"
  "fmt"
  "log"
  "path"
  "regexp"
  "strings"
)

type ILogger interface {
  Print(...interface{})
  Printf(string, ...interface{})
}

const (
  EnvDevelopment    = "development"
  DefaultViewsPath  = "views"
  DefaultPublicPath = "public"
  DefaultConfigFile = "traffic.conf"
)

var (
  env     map[string]interface{}
  logger  ILogger
)

func init() {
  env = make(map[string]interface{})
  SetLogger(log.New(os.Stderr, "", log.LstdFlags))
}

func Logger() ILogger {
  return logger
}

func SetLogger(customLogger ILogger) {
  logger = customLogger
}

func SetVar(key string, value interface{}) {
  env[key] = value
}

func GetVar(key string) interface{} {
  if value := env[key]; value != nil {
    return value
  }

  envKey := fmt.Sprintf("TRAFFIC_%s", strings.ToUpper(key))
  if value := os.Getenv(envKey); value != "" {
    return value
  }

  return nil
}

func GetStringVar(key string) string {
  value := GetVar(key)
  if s, ok := value.(string); ok {
    return s
  }

  return ""
}

func pathToRegexpString(routePath string) string {
  re := regexp.MustCompile(":[^/#?()]+")
  regexpString := re.ReplaceAllStringFunc(routePath, func(m string) string {
    return fmt.Sprintf("(?P<%s>[^/#?]+)", m[1:len(m)])
  })

  return fmt.Sprintf("^%s$", regexpString)
}

func GetStringVarWithDefault(key, defaultValue string) string {
  value := GetStringVar(key)
  if value == "" {
    return defaultValue
  }

  return value
}

func Env() string {
  return GetStringVarWithDefault("env", EnvDevelopment)
}

func RootPath() string {
  return GetStringVarWithDefault("root", ".")
}

func ViewsPath() string {
  viewsPath := GetStringVarWithDefault("views", DefaultViewsPath)
  if path.IsAbs(viewsPath) {
    return viewsPath
  }

  return path.Join(RootPath(), viewsPath)
}

func ConfigFilePath() string {
  filePath := GetStringVarWithDefault("config_file", DefaultConfigFile)
  if path.IsAbs(filePath) {
    return filePath
  }

  return path.Join(RootPath(), filePath)
}

func PublicPath() string {
  publicPath := GetStringVarWithDefault("public", DefaultPublicPath)
  if path.IsAbs(publicPath) {
    return publicPath
  }

  return path.Join(RootPath(), publicPath)
}
