package main

import (
	"archive/zip"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

var WORK_DIR = "."

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "woop")
}

func unzip(filename string) (dirname string, err error) {
	dirname, err = ioutil.TempDir(WORK_DIR, "uploaded")
	if err != nil {
		return
	}

	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			// Ordinary file.
			outpath := path.Join(dirname, f.Name)

			// Create enclosing directories.
			err = os.MkdirAll(path.Dir(outpath), os.ModeDir|os.ModePerm)
			if err != nil {
				return
			}

			// Open the file from the zip archive.
			rc, err := f.Open()
			if err != nil {
				return dirname, err
			}

			// Create the file on disk.
			out, err := os.Create(outpath)
			if err != nil {
				return dirname, err
			}

			// Copy the data to disk.
			_, err = io.Copy(out, rc)
			if err != nil {
				return dirname, err
			}

			rc.Close()
			out.Close()
		}
	}

	return
}

func receiveFile(r *http.Request, name string) (filename string, err error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	out, err := ioutil.TempFile(WORK_DIR, "uploaded")
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return out.Name(), err
	}

	return out.Name(), nil
}

func handleApp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		filename, err := receiveFile(r, "file")
		if err != nil {
			fmt.Fprintln(w, "could not get file from form")
			return
		}

		// Clean up the zip file after we've unzipped it.
		defer os.Remove(filename)

		dirname, err := unzip(filename)
		if err != nil {
			fmt.Fprintln(w, "could not unzip archive")
		}

		vars := mux.Vars(r)
		fmt.Println("got archive for", vars["name"], "at", dirname)

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
