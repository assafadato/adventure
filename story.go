package adventure

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

var defaultHandlerTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>`

var defaultTemplate *template.Template

func init()  {
	defaultTemplate = template.Must(template.New("default").Parse(defaultHandlerTemplate))
}

func JsonStory(reader io.Reader) (Story, error) {
	d := json.NewDecoder(reader)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func defaultRouteFn(r *http.Request) string {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}

	// remove preceding slash
	return path[1:]
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathParserFn (ppfn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathParserFn = ppfn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, defaultTemplate, defaultRouteFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s            Story
	t            *template.Template
	pathParserFn func(r *http.Request) string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathParserFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			fmt.Printf("%v\n", err)
			http.Error(w, "Expect the unexpected", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, fmt.Sprintf("'%s' chapter does not exist", path), http.StatusNotFound)
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Chapter  string `json:"arc"`
}
