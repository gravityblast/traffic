package traffic

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

const TemplateExtension = ".tpl"

type RenderFunc func(w ResponseWriter, template string, data interface{})

var templateManager *TemplateManager

type TemplateManager struct {
	viewsBasePath       string
	templates           *template.Template
	renderFunc          RenderFunc
	templatesReadError  error
	templatesParseError error
}

var templateFuncs = make(map[string]interface{})

func TemplateFuncs(funcs map[string]interface{}) {
	for name, fn := range funcs {
		TemplateFunc(name, fn)
	}
}

func TemplateFunc(name string, fn interface{}) {
	templateFuncs[name] = fn
}

func (t *TemplateManager) loadTemplates() {
	t.templates = template.New("templates")
	t.templates.Funcs(templateFuncs)

	if t.viewsBasePath == "" {
		panic("views base path is blank")
	}
	filepath.Walk(t.viewsBasePath, t.WalkFunc)
}

func (t *TemplateManager) WalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		t.templatesReadError = err
		t.renderFunc = t.RenderTemplateErrors

		return err
	}

	if extension := filepath.Ext(path); !info.IsDir() && extension == TemplateExtension {
		relativePath, err := filepath.Rel(t.viewsBasePath, path)
		if err != nil {
			t.templatesReadError = err
			t.renderFunc = t.RenderTemplateErrors

			return err
		}

		templateName := relativePath[0:(len(relativePath) - len(extension))]
		t.addTemplate(templateName, path)
	}

	return nil
}

func (t *TemplateManager) addTemplate(name, path string) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		t.templatesReadError = err
		t.renderFunc = t.RenderTemplateErrors

		return
	}

	templateContent := fmt.Sprintf(`{{ define "%s" }}%s{{ end }}`, name, fileContent)
	tmpl, err := t.templates.Parse(templateContent)
	if tmpl == nil && err != nil {
		t.templatesParseError = err
		t.renderFunc = t.RenderTemplateErrors
	}
}

func (t *TemplateManager) templatesPath() string {
	if path := GetVar("views"); path != nil {
		return path.(string)
	}

	return DefaultViewsPath
}

func (t *TemplateManager) Render(w ResponseWriter, template string, data interface{}) {
	err := t.templates.ExecuteTemplate(w, template, data)
	if err != nil {
		panic(err)
	}
}

func (t *TemplateManager) RenderTemplateErrors(w ResponseWriter, template string, data interface{}) {
	if t.templatesReadError != nil {
		panic(t.templatesReadError)
	}

	if t.templatesParseError != nil {
		panic(t.templatesParseError)
	}
}

func newTemplateManager() *TemplateManager {
	t := &TemplateManager{}
	t.viewsBasePath = ViewsPath()
	t.renderFunc = t.Render
	t.loadTemplates()

	return t
}

func initTemplateManager() {
	templateManager = newTemplateManager()
}

func RenderTemplate(w ResponseWriter, templateName string, data ...interface{}) {
	if len(data) == 0 {
		templateManager.renderFunc(w, templateName, nil)
	} else {
		templateManager.renderFunc(w, templateName, data[0])
	}
}
