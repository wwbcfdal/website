package main

import (
	"log"
	"net/http"
	"html/template"
	"time"
	"io/ioutil"
	"path/filepath"
	"bytes"
)

const PORT = "80"
const TEMPLATE_ROOT = "data/template"
var templates *template.Template

type Model struct {
	Title string
	Subtitle string
	Path string
	TemplateMenu template.HTML
	TemplateBody template.HTML
	Data interface{}
}

type Link struct {
	Text string
	Href string
}

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
			bytes, err := ioutil.ReadFile(templateDir + "/" + filename)
			if err != nil {
				log.Fatal(err)
			}
			content := string(bytes)
			templates = template.Must(templates.New("/"+filename).Parse(content))
			log.Println("ADDED TEMPLATE NAMED: "+ filename)
		}
	}
}

func GetModel(name string) (Model, error) {
	model := Model{
		Path: name,
		Title: "Westwood Baptist Church",
		Subtitle: "Forestdale Alabama",
		TemplateMenu: "",
		TemplateBody: "",
		Data: struct{}{},
	}
	return model, nil
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	if path == "/" || path == "/index.htm" || path == "/index" || path == "/home" || path == "/home.html" || path == "/home.htm" {
		path = "/index.html"
	}
	templateName := string(path)
	modelName := templateName
	model, _ := GetModel(modelName)
	menuBuffer := bytes.NewBufferString("")
	bodyBuffer := bytes.NewBufferString("")
	pageBuffer := bytes.NewBufferString("")
	err1 := templates.ExecuteTemplate(menuBuffer, "/menu.html", model)
	if err1 != nil {
		log.Println(err1)
		NotFoundHandler(w, r)
		return
	}
	err2 := templates.ExecuteTemplate(bodyBuffer, templateName, model)
	if err2 != nil {
		log.Println(err2)
		NotFoundHandler(w, r)
		return
	}
	menu := menuBuffer.String()
	body := bodyBuffer.String()
	model.TemplateMenu = template.HTML(menu);
	model.TemplateBody = template.HTML(body);
	err3 := templates.ExecuteTemplate(pageBuffer, "/frame.html", model)
	if err3 != nil {
		log.Println(err3)
		NotFoundHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(pageBuffer.Bytes())
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	templateName := "/404.html"
	modelName := templateName
	model, _ := GetModel(modelName)
	menuBuffer := bytes.NewBufferString("")
	bodyBuffer := bytes.NewBufferString("")
	pageBuffer := bytes.NewBufferString("")
	err1 := templates.ExecuteTemplate(menuBuffer, "/menu.html", model)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 File Not Found"))
		return
	}
	err2 := templates.ExecuteTemplate(bodyBuffer, templateName, model)
	if err2 != nil {
		log.Println(err2)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 File Not Found"))
		return
	}
	menu := menuBuffer.String()
	body := bodyBuffer.String()
	model.TemplateMenu = template.HTML(menu);
	model.TemplateBody = template.HTML(body);
	err3 := templates.ExecuteTemplate(pageBuffer, "/frame.html", model)
	if err3 != nil {
		log.Println(err3)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 File Not Found"))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(pageBuffer.Bytes())
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
	
	log.Println("Listening on port: 80")

	// Create Server
	server := http.Server{
		Addr: ":" + PORT,
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
