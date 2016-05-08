package static

import (
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/mephux/kolide/config"
	"github.com/mephux/kolide/controller/v1"
)

// Load templates
func Load(configuration *config.Config) *template.Template {
	tmpl := template.New("_")
	tmpl.Funcs(v1.HelperFunctions)
	tmpl.Delims("<%", "%>")

	var dir []string

	if configuration.Server.Production {
		dir, _ = AssetDir("views")
	} else {
		viewFiles, _ := ioutil.ReadDir("./static/ui/views")

		for _, vf := range viewFiles {
			dir = append(dir, vf.Name())
		}
	}

	for _, name := range dir {
		if filepath.Ext(name) != ".html" {
			continue
		}

		var src []byte
		var err error

		if configuration.Server.Production {
			src = MustAsset(filepath.Join("views", name))
		} else {
			src, err = ioutil.ReadFile(filepath.Join("./static/ui/views", name))

			if err != nil {
				panic(err)
			}
		}

		tmpl = template.Must(tmpl.New(name).Parse(string(src)))
	}

	return tmpl
}
