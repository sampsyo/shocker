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

func unzip(filename string) (err error) {
	dirname, err := ioutil.TempDir(WORK_DIR, "uploaded")
	if err != nil {
		fmt.Println("could not create directory to receive")
		return err
	}

	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		fmt.Println("expanding", f.Name)
		if f.FileInfo().IsDir() {
			// Directory.
			fmt.Println("is directory")

		} else {
			// Ordinary file.
			outpath := path.Join(dirname, f.Name)
			fmt.Println("into", outpath)

			// Create enclosing directories.
			err = os.MkdirAll(path.Dir(outpath), os.ModeDir|os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			// Open the file from the zip archive.
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			// Create the file on disk.
			out, err := os.Create(outpath)
			if err != nil {
				log.Fatal(err)
			}

			// Copy the data to disk.
			_, err = io.Copy(out, rc)
			if err != nil {
				log.Fatal(err)
			}

			rc.Close()
			out.Close()
		}
	}

	return nil
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
		}

		vars := mux.Vars(r)
		fmt.Println("got file for", vars["name"])

		unzip(filename)

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
