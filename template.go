package traffic

import (
  "os"
  "fmt"
  "html/template"
  "net/http"
  "path/filepath"
  "io/ioutil"
)

const DefaultViewsPath = "views"

type RenderFunc func(w http.ResponseWriter, template string, data interface{})

var templateManager *TemplateManager

type TemplateManager struct {
  viewsBasePath       string
  templates           *template.Template
  renderFunc          RenderFunc
  templatesReadError  error
  templatesParseError error
}

func (t *TemplateManager) loadTemplates() {
  t.templates = template.New("templates")
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

  if extension := filepath.Ext(path); !info.IsDir() && extension == ".tmpl" {
    relativePath, err  := filepath.Rel(t.viewsBasePath, path)
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

func (t *TemplateManager) Render(w http.ResponseWriter, template string, data interface{}) {
  err := t.templates.ExecuteTemplate(w, template, data)
  if err != nil {
    panic(err)
  }
}

func (t *TemplateManager) RenderTemplateErrors(w http.ResponseWriter, template string, data interface{}) {
  if t.templatesReadError != nil {
    panic(t.templatesReadError)
  }

  if t.templatesParseError != nil {
    panic(t.templatesParseError)
  }
}

func newTemplateManager() *TemplateManager {
  t := &TemplateManager{}
  if path := GetVar("views"); path != nil {
    t.viewsBasePath = path.(string)
  } else {
    t.viewsBasePath = DefaultViewsPath
  }

  t.renderFunc = t.Render
  t.loadTemplates()

  return t
}

func initTemplateManager() {
  templateManager = newTemplateManager()
}

func Render(w http.ResponseWriter, template string, data interface{}) {
  templateManager.renderFunc(w, template, data)
}
