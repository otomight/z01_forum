package templates

import (
	"fmt"
	"forum/internal/config"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates	*template.Template

func htmlFilesWalkDirFunc(files *[]string) fs.WalkDirFunc {
	return func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dir.IsDir() && strings.HasSuffix(dir.Name(), ".html") {
			*files = append(*files, path)
		}
		return nil
	}
}

func getHtmlFiles(dirPath string) ([]string, error) {
	var files	[]string
	var err		error

	files = []string{}
	err = filepath.WalkDir(dirPath, htmlFilesWalkDirFunc(&files))
	return files, err
}

func LoadTemplates() error {
	var	funcMap	template.FuncMap
	var	files	[]string
	var	err		error

	funcMap = template.FuncMap{
		// insert here function to use into the templates
		// functions can be developped and stored into this package
		// "myTemplateFunc": myFunc,
		"doesStrMatchAny": doesStrMatchAny,
		"addToStruct": addToStruct,
	}
	files, err = getHtmlFiles(config.TemplatesDir)
	if err != nil {
		return err
	}
	templates, err = template.New("layout.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		return err
	}
	return nil
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	var	err		error

	if templates == nil {
		fmt.Println("No templates loaded")
		return
	}
	err = templates.ExecuteTemplate(w, tmpl + ".html", data)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error rendering template",
						http.StatusInternalServerError)
	}
}
