package traffic

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
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
	DefaultPort       = 3000
	DefaultHost       = "127.0.0.1"
)

var (
	env    map[string]interface{}
	logger ILogger
)

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

func getStringVar(key string) string {
	value := GetVar(key)
	if s, ok := value.(string); ok {
		return s
	}

	return ""
}

func pathToRegexp(routePath string) (*regexp.Regexp, bool) {
	var (
		re       *regexp.Regexp
		isStatic bool
	)

	regexpString := routePath

	isStaticRegexp := regexp.MustCompile(`[\(\)\?\<\>:]`)
	if !isStaticRegexp.MatchString(routePath) {
		isStatic = true
	}

	// Dots
	re = regexp.MustCompile(`([^\\])\.`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf(`%s\.`, string(m[0]))
	})

	// Wildcard names
	re = regexp.MustCompile(`:[^/#?()\.\\]+\*`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf("(?P<%s>.+)", m[1:len(m)-1])
	})

	re = regexp.MustCompile(`:[^/#?()\.\\]+`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)])
	})

	s := fmt.Sprintf(`\A%s\z`, regexpString)

	return regexp.MustCompile(s), isStatic
}

func getStringVarWithDefault(key, defaultValue string) string {
	value := getStringVar(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func Env() string {
	return getStringVarWithDefault("env", EnvDevelopment)
}

func RootPath() string {
	return getStringVarWithDefault("root", ".")
}

func ViewsPath() string {
	viewsPath := getStringVarWithDefault("views", DefaultViewsPath)
	if path.IsAbs(viewsPath) {
		return viewsPath
	}

	return path.Join(RootPath(), viewsPath)
}

func ConfigFilePath() string {
	filePath := getStringVarWithDefault("config_file", DefaultConfigFile)
	if path.IsAbs(filePath) {
		return filePath
	}

	return path.Join(RootPath(), filePath)
}

func PublicPath() string {
	publicPath := getStringVarWithDefault("public", DefaultPublicPath)
	if path.IsAbs(publicPath) {
		return publicPath
	}

	return path.Join(RootPath(), publicPath)
}

func SetPort(port int) {
	SetVar("port", port)
}

func Port() int {
	port := GetVar("port")

	if port == nil {
		return DefaultPort
	}

	if i, ok := port.(int); ok {
		return i
	}

	if s, ok := port.(string); ok {
		i, err := strconv.Atoi(s)
		if err != nil {
			return DefaultPort
		}

		return i
	}

	return DefaultPort
}

func Host() string {
	host := getStringVarWithDefault("host", DefaultHost)

	return host
}

func SetHost(host string) {
	SetVar("host", host)
}
