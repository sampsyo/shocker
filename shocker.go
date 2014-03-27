package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

var WORK_DIR = "."

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "woop")
}

func handleApp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintln(w, "could not get file from form")
			return
		}
		defer file.Close()

		vars := mux.Vars(r)
		fmt.Println("got file", header.Filename, "for", vars["name"])

		out, err := ioutil.TempFile(WORK_DIR, "uploaded")
		if err != nil {
			fmt.Fprintln(w, "could not open file to save")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Println(w, "could not copy file")
			return
		}

		fmt.Fprintln(w, "success")
	} else {
		fmt.Fprintf(w, "nothing to see here")
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/apps/{name}", handleApp)
	r.HandleFunc("/", handleHome)

	fmt.Println("http://0.0.0.0:8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
