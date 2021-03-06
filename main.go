package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	tpl       *template.Template
	pass      = os.Getenv("pass")
	apiKey    = ""
	cookieVal = "123"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))
}

func cookieCheck(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "session",
			Value: "0",
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		log.Println("No Cookie")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		input := r.FormValue("inputPassword")
		if input == pass {
			cookie = &http.Cookie{
				Name:  "session",
				Value: cookieVal,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		}
	}

	if cookie.Value != cookieVal {
		cookie = &http.Cookie{
			Name:  "session",
			Value: "0",
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

// IndexHandler Dashboard "/" end point
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//cookieCheck(w, r)

	dashboard, err := dashboard()
	if err != nil {
		return
	}

	var data []interface{}
	data = append(data, apiKey)
	data = append(data, dashboard)

	err = tpl.ExecuteTemplate(w, "index.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

//ContainersHandler Containers "/containers" end point
func ContainersHandler(w http.ResponseWriter, r *http.Request) {
	//cookieCheck(w, r)

	containers, err := container()

	if err != nil {
		return
	}

	logs, err := logs(containers)
	if err != nil {
		return
	}

	var data []interface{}
	data = append(data, apiKey)
	data = append(data, containers)
	data = append(data, logs)

	err = tpl.ExecuteTemplate(w, "containers.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

// StatsHandler Stats "/stats" end point
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	//cookieCheck(w, r)
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	stats, err := stats()

	if err != nil {
		return
	}

	var data []interface{}
	data = append(data, apiKey)
	data = append(data, stats)

	err = tpl.ExecuteTemplate(w, "stats.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

// ImagesHandler Images "/images" end point
func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	//cookieCheck(w, r)
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	images, err := images()

	if err != nil {
		return
	}

	var data []interface{}
	data = append(data, apiKey)
	data = append(data, images)

	err = tpl.ExecuteTemplate(w, "images.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

// VolumesHandler Volumes "/volumes" end point
func VolumesHandler(w http.ResponseWriter, r *http.Request) {
	//cookieCheck(w, r)
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	volumes, err := volumes()

	if err != nil {
		return
	}

	var data []interface{}
	data = append(data, apiKey)
	data = append(data, volumes)

	err = tpl.ExecuteTemplate(w, "volumes.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

// NetworksHandler Networks "/networks" end point
func NetworksHandler(w http.ResponseWriter, r *http.Request) {
	//cookieCheck(w, r)
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	networks, err := networks()

	if err != nil {
		return
	}

	var data []interface{}
	data = append(data, apiKey)
	data = append(data, networks)

	err = tpl.ExecuteTemplate(w, "networks.tmpl", data)
	if err != nil {
		log.Println(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		log.Println(r.Method, http.StatusNotFound, r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(r.Method, http.StatusFound, r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	value := cookie.Value

	if value == cookieVal {
		log.Println(r.Method, http.StatusFound, r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = tpl.ExecuteTemplate(w, "login.tmpl", nil)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(r.Method, http.StatusOK, r.URL.Path)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		log.Println(r.Method, http.StatusFound, r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cookie, _ := r.Cookie("session")
	cookie = &http.Cookie{
		Name:  "session",
		Value: "0",
		Path:  "/",
	}

	http.SetCookie(w, cookie)
	log.Println(r.Method, http.StatusOK, r.URL.Path)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	apiKey = GenerateAPIPassword(32)
	cookieVal = GenerateAPIPassword(140)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/containers", ContainersHandler)
	http.HandleFunc("/stats", StatsHandler)
	http.HandleFunc("/images", ImagesHandler)
	http.HandleFunc("/volumes", VolumesHandler)
	http.HandleFunc("/networks", NetworksHandler)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.HandleFunc("/api/containers", APIContainer)
	http.HandleFunc("/api/images", APIImages)
	http.HandleFunc("/api/volumes", APIVolumes)
	http.HandleFunc("/api/networks", APINetworks)
	http.HandleFunc("/api/stats", APIStats)
	http.HandleFunc("/api/logs", APILogs)
	http.HandleFunc("/api/status", APIStatus)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening http://0.0.0.0:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
