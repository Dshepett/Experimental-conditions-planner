package app

import (
	"bytes"
	"client/internal/models"
	"client/internal/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")
	tmpl, err := template.ParseFiles("web/template/index.html")
	if err != nil {
		a.logger.Errorln("Error: ", err)
	} else {
		a.logger.Infoln("Open index page")
		wrongTable := errorMsg == "true"
		tmpl.Execute(w, wrongTable)
	}
}

func (a *App) helpHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/help.html")
	if err != nil {
		a.logger.Errorln("Error: ", err)
	} else {
		a.logger.Infoln("Open help page")
		tmpl.Execute(w, nil)
	}
}

func (a *App) downloadHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Infoln("Download file with result data")
	res, _ := ioutil.ReadAll(r.Body)
	tempdata := []byte(res)
	name := fmt.Sprintf("%d.csv", rand.Intn(100))
	file, _ := os.Create(name)
	file.Write(tempdata)
	file.Close()
	Openfile, _ := os.Open(name)
	tempBuffer := make([]byte, 5*1024*1024)
	Openfile.Read(tempBuffer)
	FileContentType := http.DetectContentType(tempBuffer)
	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)
	Filename := "result.csv"
	w.Header().Set("Content-Type", FileContentType+";"+Filename)
	w.Header().Set("Content-Length", FileSize)
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)
	Openfile.Close()
	os.Remove(name)
}

func (a *App) predictHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Infoln("Predicting..")
	type Response struct {
		Status     int                 `json:"status"`
		Message    string              `json:"message"`
		Conditions []models.Conditions `json:"conditions"`
	}
	var res Response
	var data []models.Characteristics
	var err error
	if r.Header.Get("Content-type") != "application/x-www-form-urlencoded" {
		r.ParseMultipartForm(5 * 1024 * 1024)
		file, header, _ := r.FormFile("file")
		defer file.Close()
		reader := csv.NewReader(file)
		reader.Comma = ';'
		res, _ := reader.ReadAll()
		data, err = utils.ConvertToCharacteristics(res)
		filename := strings.Split(header.Filename, ".")
		if filename[len(filename)-1] != "csv" {
			http.Redirect(w, r, "/?error=true", http.StatusMovedPermanently)
			return
		}
		if err != nil {
			http.Redirect(w, r, "/?error=true", http.StatusMovedPermanently)
			return
		}
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
	type Result struct {
		Error   bool
		Data    []models.SynthesisData
		TXTdata string
		Port    string
	}
	result := &Result{}
	if err != nil {
		fmt.Print(err)
		result.Error = true
		result.Data = nil
		result.TXTdata = ""
	} else {
		json.NewDecoder(resp.Body).Decode(&res)
		synth := []models.SynthesisData{}
		for i := range data {
			synth = append(synth, models.SynthesisData{Conditions: res.Conditions[i], Characteristics: data[i]})
		}
		result.Error = false
		result.Data = synth
		result.TXTdata = utils.ConvertToCsv(result.Data)
		result.Port = a.config.Port
	}
	tmpl, err := template.ParseFiles("web/template/predict.html")
	if err != nil {
		a.logger.Errorln("Error: ", err)
	} else {
		tmpl.Execute(w, result)
	}
}
