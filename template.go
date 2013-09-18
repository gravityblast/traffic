package traffic

import (
  "html/template"
  "net/http"
  "os"
  "path/filepath"
)

const DefaultViewsPath = "views"

type RenderFunc func(w http.ResponseWriter, template string, data interface{})

var templateManager *TemplateManager

type TemplateManager struct {
  templates *template.Template
  templatesReadError error
  templatesParseError error
  renderFunc RenderFunc
}

func (t *TemplateManager) loadTemplates() {
  paths := t.getTemplatesPaths()
  t.templates, t.templatesParseError = template.ParseFiles(paths...)
  if t.templatesParseError != nil {
    t.renderFunc = t.RenderTemplateErrors
  }
}

func (t *TemplateManager) getTemplatesPaths() []string {
  views := make([]string, 0)
  filepath.Walk(t.templatesPath(), func(path string, info os.FileInfo, err error) (e error) {
    if err != nil {
      t.templatesReadError = err
      t.renderFunc = t.RenderTemplateErrors
    } else {
      fileInfo, err := os.Stat(path)
      if err == nil && !fileInfo.IsDir() && filepath.Ext(path) == ".tmpl" {
        views = append(views, path)
      }
    }

    return
  })

  return views
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
