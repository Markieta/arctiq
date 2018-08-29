package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const version string = "1.0.0"

func main() {
	showversion := flag.Bool("version", false, "display version")
	frontend := flag.Bool("frontend", false, "run in frontend mode")
	port := flag.Int("port", 8080, "port to bind")
	backend := flag.String("backend-service", "http://127.0.0.1:8081", "hostname of backend server")
	flag.Parse()

	if *showversion {
		fmt.Printf("Version %s\n", version)
		return
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	if *frontend {
		frontendMode(*port, *backend)
	} else {
		backendMode(*port)
	}

}

func backendMode(port int) {
	log.Println("Operating in backend mode...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "Hello World...")
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}

func frontendMode(port int, backendURL string) {
	log.Println("Operating in frontend mode...")
	tpl := template.Must(template.New("out").Parse(html))

	transport := http.Transport{DisableKeepAlives: false}
	client := &http.Client{Transport: &transport}
	req, _ := http.NewRequest(
		"GET",
		backendURL,
		nil,
	)
	req.Close = false

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, "x")
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Backend could not be connected to: %s", err.Error())
			return
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
