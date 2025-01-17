package assets

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"

	"github.com/rs/zerolog/log"
)

type assetsfs struct {
	embedded  *embed.FS
	templates map[string]*template.Template
}

var FS assetsfs

func (afs assetsfs) Open(name string) (fs.File, error) {
	if afs.embedded != nil {
		return afs.embedded.Open(name)
	}
	return os.DirFS("src/assets").Open(name)
}

func Template(path string) *template.Template {
	var tmpl *template.Template
	tmpl, found := FS.templates[path]
	if !found {
		tmpl = template.Must(template.New(path).Delims("{%", "%}").Funcs(template.FuncMap{
			"inline": func(svg string) template.HTML {
				svgfile, err := FS.Open("graphicarts/" + svg)
				// should never happen
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				defer svgfile.Close()

				content, err := io.ReadAll(svgfile)
				// should never happen
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				return template.HTML(content)
			},
		}).ParseFS(FS, path))
		if FS.embedded != nil {
			FS.templates[path] = tmpl
		}
	}
	return tmpl
}

func Render(path string, writer io.Writer, data interface{}) {
	tmpl := Template(path)
	tmpl.Execute(writer, data)
}

func init() {
	FS.templates = make(map[string]*template.Template)
}
