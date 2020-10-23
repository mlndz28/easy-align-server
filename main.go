package main

import(
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"encoding/json"
	"net/http"
	"log"
	"os/exec"
	"errors"
	"fmt"

	"github.com/mlndz28/praatgo"
)

const maxBuffer = 64 // in MB

func main() {
	http.HandleFunc("/align", align)
	http.HandleFunc("/", notFound)
	fmt.Println("Ready!")
	log.Fatal(http.ListenAndServe(":7728", nil))
}


func notFound(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<div style='text-align: center;'><h1 style='font-size:100px;top: 50%;position: relative'>404</h1><a href='https://www.swagger.io'>Swagger</a></div>"))
}

func align(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	transcript, audio, header, status, err := parseAlignRequest(r)
	if err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{"error":err.Error()})
		return 
	} 
		
	err = saveAlignRequest(w, audio, transcript, header)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error":err.Error()})
		return
	}
	
	tg, err, log := callEasyAlign("/tmp/" + header.Filename)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"error":err.Error(), "log": string(log)})
		return
	}
	err = json.NewEncoder(w).Encode(tg)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"error":err.Error()})
		return
	}
}

func parseAlignRequest(r *http.Request)(transcript, audio []byte, header *multipart.FileHeader, status int, err error){
	err = r.ParseMultipartForm(maxBuffer << 20)
	if err != nil {
		status = 413	// file bigger than the buffer capacity
		return 
	}
	transcript = []byte(r.FormValue("transcript"))
	buffer, header, err := r.FormFile("audio")
	if err != nil {
		err = errors.New("No audio detected in the request")
		status = http.StatusBadRequest
		return 
		} 
	audio =  make([]byte,header.Size)
	_, err = buffer.Read(audio)
	if err != nil {	
		fmt.Println("couldnt read audio")
		status = http.StatusBadRequest
	} else {
		status = 200
	}
	return
}

func saveAlignRequest(w http.ResponseWriter, audio, transcript []byte, header *multipart.FileHeader) (err error) {
	err = ioutil.WriteFile("/tmp/" + header.Filename, audio, 0644)
	if err!= nil{
		return 
	}
	err = ioutil.WriteFile("/tmp/" + header.Filename+".txt", []byte(transcript), 0644)
	return
}

func callEasyAlign(path string) (praatgo.TextGrid, error, []byte) {
	fmt.Println("easyalign","-o",path+".tg","-l","spa","all",path,path+".txt")
	cmd := exec.Command("easyalign","-o",path+".tg","-l","spa","all",path,path+".txt")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(string(stderr.Bytes()))
        return praatgo.TextGrid{}, err, stderr.Bytes()
    }

	content, err := ioutil.ReadFile(path+".tg")
	if err != nil{
		return praatgo.TextGrid{}, err, stdout.Bytes()
	}
	tg, err := praatgo.DeserializeTextGrid(content)
	return tg, err, stdout.Bytes()
}

// common error handling function
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{"error":err.Error()})
}


