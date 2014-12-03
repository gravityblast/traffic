package traffic

import (
	"os"
	"testing"

	assert "github.com/pilu/miniassert"
)

func resetGlobalEnv() {
	env = make(map[string]interface{})
}

func TestSetVar(t *testing.T) {
	resetGlobalEnv()
	SetVar("foo", "bar")
	assert.Equal(t, "bar", env["foo"])
	resetGlobalEnv()
}

func TestGetVar(t *testing.T) {
	resetGlobalEnv()
	env["foo-2"] = "bar-2"
	assert.Equal(t, "bar-2", GetVar("foo-2"))

	assert.Nil(t, GetVar("os_foo"))
	os.Setenv("TRAFFIC_OS_FOO", "bar")
	assert.Equal(t, "bar", GetVar("os_foo"))

	resetGlobalEnv()
}

func TestGetStringVar(t *testing.T) {
	resetGlobalEnv()

	assert.Equal(t, "", getStringVar("foo"))
	SetVar("foo", "bar")
	assert.Equal(t, "bar", getStringVar("foo"))

	resetGlobalEnv()
}

func TestPathToRegexp(t *testing.T) {
	tests := [][]interface{}{
		{
			`/`,
			`\A/\z`,
			true,
		},
		{
			`/foo/bar`,
			`\A/foo/bar\z`,
			true,
		},
		{
			`/foo/bar`,
			`\A/foo/bar\z`,
			true,
		},
		{
			`/:foo/bar/:baz`,
			`\A/(?P<foo>[^/#?]+)/bar/(?P<baz>[^/#?]+)\z`,
			false,
		},
		{
			`(/categories/:category_id)?/posts/:id`,
			`\A(/categories/(?P<category_id>[^/#?]+))?/posts/(?P<id>[^/#?]+)\z`,
			false,
		},
	}

	for _, pair := range tests {
		path := pair[0].(string)
		expectedRegexp := pair[1].(string)
		expectedIsStatic := pair[2].(bool)

		r, isStatic := pathToRegexp(path)
		assert.Equal(t, expectedRegexp, r.String())
		assert.Equal(t, expectedIsStatic, isStatic)
	}
}

func TestEnv(t *testing.T) {
	resetGlobalEnv()
	assert.Equal(t, EnvDevelopment, Env())

	SetVar("env", "production")
	assert.Equal(t, "production", Env())
	resetGlobalEnv()
}

func TestRootPath(t *testing.T) {
	resetGlobalEnv()
	assert.Equal(t, ".", RootPath())

	SetVar("root", "foo")
	assert.Equal(t, "foo", RootPath())
	resetGlobalEnv()
}

func TestViewsPath(t *testing.T) {
	resetGlobalEnv()
	assert.Equal(t, DefaultViewsPath, ViewsPath())

	SetVar("views", "foo/views")
	assert.Equal(t, "foo/views", ViewsPath())

	SetVar("root", "/root")
	assert.Equal(t, "/root/foo/views", ViewsPath())

	SetVar("views", "/absolute/path")
	assert.Equal(t, "/absolute/path", ViewsPath())

	resetGlobalEnv()
}

func TestConfigFilePath(t *testing.T) {
	resetGlobalEnv()

	assert.Equal(t, DefaultConfigFile, ConfigFilePath())

	SetVar("config_file", "foo.conf")
	assert.Equal(t, "foo.conf", ConfigFilePath())

	SetVar("root", "/root")
	assert.Equal(t, "/root/foo.conf", ConfigFilePath())

	SetVar("config_file", "/absolute/foo.conf")
	assert.Equal(t, "/absolute/foo.conf", ConfigFilePath())

	resetGlobalEnv()
}

func TestPublicPath(t *testing.T) {
	resetGlobalEnv()

	assert.Equal(t, DefaultPublicPath, PublicPath())

	SetVar("public", "custom-public")
	assert.Equal(t, "custom-public", PublicPath())

	SetVar("root", "/root")
	assert.Equal(t, "/root/custom-public", PublicPath())

	SetVar("public", "/absolute/public")
	assert.Equal(t, "/absolute/public", PublicPath())

	resetGlobalEnv()
}

func TestPort(t *testing.T) {
	resetGlobalEnv()

	assert.Equal(t, DefaultPort, Port())

	SetVar("port", 80)
	assert.Equal(t, 80, Port())

	SetVar("port", "80")
	assert.Equal(t, 80, Port())

	SetVar("port", "")
	assert.Equal(t, DefaultPort, Port())

	resetGlobalEnv()
}

func TestSetPort(t *testing.T) {
	resetGlobalEnv()

	assert.Nil(t, GetVar("port"))

	SetPort(80)
	assert.Equal(t, 80, GetVar("port"))

	resetGlobalEnv()
}

func TestHost(t *testing.T) {
	resetGlobalEnv()

	assert.Equal(t, DefaultHost, Host())

	SetVar("host", "1.2.3.4")
	assert.Equal(t, "1.2.3.4", Host())

	resetGlobalEnv()
}

func TestSetHost(t *testing.T) {
	resetGlobalEnv()

	assert.Nil(t, GetVar("host"))

	SetHost("1.2.3.4")
	assert.Equal(t, "1.2.3.4", GetVar("host"))

	resetGlobalEnv()
}
