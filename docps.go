package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Docker struct {
	Name       string `json:"name"`
	Image      string `json:"image"`
	Size       string `json:"size"`
	RunningFor string `json:"runningFor"`
	Status     string `json:"status"`
}

//IndexHandler Execute the docker ps -a command and reading the stdout
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var container []Docker
	// Docker ps -a argument with \t for splitting string later.
	cmdArgs := []string{"ps", "-a", "--format", "{{.Names}}\t{{.Image}}\t{{.Size}}\t{{.RunningFor}}\t{{.Status}}"}

	cmd := exec.Command("docker", cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err.Error())
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			outPut := scanner.Text()

			s := strings.Split(outPut, "\t")

			container = append(container,
				Docker{Name: s[0],
					Image:      s[1],
					Size:       s[2],
					RunningFor: s[3],
					Status:     s[4][:1],
				})
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err.Error())
		return
	}

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, container)

}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", IndexHandler)
	return r
}

func main() {
	log.Println("Listening:8080..")

	err := http.ListenAndServe(":8080", handler())
	if err != nil {
		log.Fatal(err)
	}

}