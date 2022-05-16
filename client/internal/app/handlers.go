package app

import (
	"bytes"
	"client/internal/models"
	"client/internal/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/index.html")
	if err != nil {
		log.Print(err)
	} else {
		tmpl.Execute(w, nil)
	}
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/help.html")
	if err != nil {
		log.Print(err)
	} else {
		tmpl.Execute(w, nil)
	}
}

var filename string = "data.csv"

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	// data := r.PostForm.Get("data")
	// tempBuffer := []byte(data)
	// FileContentType := http.DetectContentType(tempBuffer)
	// Filename := "table.csv"
	// w.Header().Set("Content-Type", FileContentType+";"+Filename)
	// io.Copy(w, bytes.NewReader(tempBuffer))
	Openfile, err := os.Open(filename)
	defer Openfile.Close() //Close after function return

	if err != nil {
		log.Print(err)
		http.Error(w, "File not found.", 404) //return 404 if file is not found
		return
	}

	tempBuffer := make([]byte, 512)                       //Create a byte array to read the file later
	Openfile.Read(tempBuffer)                             //Read the file into  byte
	FileContentType := http.DetectContentType(tempBuffer) //Get file header

	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	Filename := "result.csv"

	//Set the headers
	w.Header().Set("Content-Type", FileContentType+";"+Filename)
	w.Header().Set("Content-Length", FileSize)
	Openfile.Seek(0, 0)  //We read 512 bytes from the file already so we reset the offset back to 0
	io.Copy(w, Openfile) //'Copy' the file to the client
}

func (a *App) predictHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status     int                 `json:"status"`
		Message    string              `json:"message"`
		Conditions []models.Conditions `json:"conditions"`
	}
	var res Response
	var data []models.Characteristics
	if r.Header.Get("Content-type") != "application/x-www-form-urlencoded" {
		r.ParseMultipartForm(5 * 1024 * 1024)
		file, _, _ := r.FormFile("file")
		defer file.Close()
		reader := csv.NewReader(file)
		reader.Comma = ';'
		res, _ := reader.ReadAll()
		data, _ = utils.ConvertToCharacteristics(res)
	} else {
		r.ParseForm()
		size_txt := r.PostForm.Get("size")
		consistance_txt := r.PostForm.Get("consistance")
		stability_txt := r.PostForm.Get("stability")
		size, _ := strconv.ParseFloat(size_txt, 64)
		consistance, _ := strconv.ParseFloat(consistance_txt, 64)
		stability, _ := strconv.ParseFloat(stability_txt, 64)
		data = append(data, models.Characteristics{Size: size, Consistence: consistance, Stability: stability})
	}
	response, _ := json.Marshal(map[string][]models.Characteristics{"characteristics": data})
	addr := fmt.Sprintf("%s/calculate", a.config.APIAddres)
	resp, err := http.Post(addr, "application/json", bytes.NewReader(response))
	if err != nil {
		log.Print(err)
	} else {
		json.NewDecoder(resp.Body).Decode(&res)
		synth := []models.SynthesisData{}
		for i := range data {
			synth = append(synth, models.SynthesisData{Conditions: res.Conditions[i], Characteristics: data[i]})
		}
		tmpl, err := template.ParseFiles("web/template/predict.html")
		if err != nil {
			log.Print(err)
		} else {
			log.Printf("Status:%d, ,essage: %s", res.Status, res.Message)
			tmpl.Execute(w, synth)
		}
	}
}
