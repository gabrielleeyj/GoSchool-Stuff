package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	gHandlers "github.com/gorilla/handlers"
	"github.com/venuebooking/app/venuebooking_v1/server/internal/config"
	"github.com/venuebooking/app/venuebooking_v1/server/internal/handlers"
)

func Start() {
	r, err := initRoutes()
	if err != nil {
		log.Fatalf("cannot initiate router: %v", err)
	}
	err = parseTemplates(config.ViewsDir())
	if err != nil {
		log.Fatalf("cannot parse templates in %q: %v", config.ViewsDir(), err)
	}

	log.Printf("server listening at %v on http...", config.HTTPAddr())
	log.Fatal(http.ListenAndServe(config.HTTPAddr(),
		gHandlers.LoggingHandler(os.Stdout, r)))
	return
}

func parseTemplates(dir string) error {
	var allFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename := file.Name()
		allFiles = append(allFiles, filepath.Join(dir, filename))
	}

	t := template.New("").Funcs(template.FuncMap{
		"set": func(renderArgs map[string]interface{}, key string, value interface{}) template.JS {
			renderArgs[key] = value
			return template.JS("")
		},
		"year": func() int {
			return time.Now().Year()
		},
		"timeFormat": func(mydate time.Time) string {
			zero := time.Time{}
			if mydate == zero {
				return ""
			}
			return mydate.Format("01/02/2006")
		},
	})

	templates, err := t.ParseFiles(allFiles...)
	if err != nil {
		return err
	}

	handlers.Templates = templates

	return nil
}
