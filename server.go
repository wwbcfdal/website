package main

import (
	"log"
	"net/http"
	"html/template"
	"time"
	"io/ioutil"
	"os"
	"path/filepath"
)

const TEMPLATE_ROOT = "data/template"
var templates *template.Template

func InitTemplates() {
	templates = template.Must(template.New("").Parse(""))
	templateDir, err := filepath.Abs(TEMPLATE_ROOT)
	if err != nil {
		log.Fatal(err)
	}
	filesInfo, err := ioutil.ReadDir(templateDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range filesInfo {
		if !info.IsDir() {
			filename := info.Name()
			log.Println("ADDING: " + filename)
			bytes, err := ioutil.ReadFile(templateDir + "/" + filename)
			if err != nil {
				log.Fatal(err)
			}
			content := string(bytes)
			templates = template.Must(templates.New(filename).Parse(content))
			log.Println("ADDED TEMPLATE NAMED: "+ filename)
		}
	}
}

// exists returns whether the given file or directory exists or not
func existsAbsolute(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func existsRelative(path string) (bool, error) {
	absolute, err := filepath.Abs(path)
	log.Println(absolute)
	stat, err := os.Stat(absolute)
	if err == nil && !stat.IsDir() { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func MarkdownRenderHandler(w http.ResponseWriter, r *http.Request) {
	s := "dynamic data"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	templateName := string(path[1:])
	w.WriteHeader(http.StatusOK)
	//log.Println("Executing Template Named: "+ templateName)
	templates.ExecuteTemplate(w, templateName, nil)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Status 404"))
}

func main() {

	InitTemplates()

	// Primary servermux routes
	routes := http.NewServeMux()

	// Static server
	fs := http.FileServer(http.Dir("static"))

	// Template server
	routes.HandleFunc("/", TemplateHandler)

	// Init routes
	routes.Handle("/js/", fs)
	routes.Handle("/css/", fs)
	routes.Handle("/fonts/", fs)
	routes.Handle("/images/", fs)
	routes.HandleFunc("/markdown/", MarkdownRenderHandler)
	
	log.Println("Listening on port: 80")

	// Create Server
	server := http.Server{
		Addr: ":http",
		Handler: routes,
		ReadTimeout: time.Second*10,
		WriteTimeout: time.Second*10,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
