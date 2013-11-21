package traffic

import (
  "os"
  "bufio"
  "io/ioutil"
  "html/template"
)

const buildErrorFilePath = "./tmp/traffic-errors.txt"

type BuildErrorsMiddleware struct {}

func (middleware BuildErrorsMiddleware) readErrorFile() string {
  file, err := os.Open(buildErrorFilePath)
  if err != nil {
    return ""
  }

  defer file.Close()

  reader := bufio.NewReader(file)
  bytes, _ := ioutil.ReadAll(reader)

  return string(bytes)
}

func (middleware BuildErrorsMiddleware) RenderError(w ResponseWriter, r *Request) {
  data := map[string]interface{} {
    "Output": middleware.readErrorFile(),
  }

  w.Header().Set("Content-Type", "text/html")
  tpl := template.Must(template.New("ErrorPage").Parse(buildPageTpl))
  tpl.Execute(w, data)
}

func (middleware BuildErrorsMiddleware) ServeHTTP(w ResponseWriter, r *Request, next NextMiddlewareFunc) (ResponseWriter, *Request) {
  if _, err := os.Stat(buildErrorFilePath); err == nil {
    middleware.RenderError(w, r)

    return w, r
  }

  if nextMiddleware := next(); nextMiddleware != nil {
    w, r = nextMiddleware.ServeHTTP(w, r, next)
  }

  return w, r
}

const buildPageTpl string = `
  <html>
    <head>
      <title>Traffic Panic</title>
      <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
      <style>
      html, body{ padding: 0; margin: 0; }
      header { background: #C52F24; color: white; border-bottom: 2px solid #9C0606; }
      h1 { padding: 10px 0; margin: 0; }
      .container { margin: 0 20px; }
      .output { height: 300px; overflow-y: scroll; border: 1px solid #e5e5e5; padding: 10px; }
      </style>
    </head>
  <body>
    <header>
      <div class="container">
        <h1>Build Error</h1>
      </div>
    </header>

    <div class="container">
      <pre class="output">{{ .Output }}</pre>
    </div>
  </body>
  </html>
  `
