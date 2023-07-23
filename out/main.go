package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/files", getFileData).Methods("GET")
	http.ListenAndServe(":8080", r)
}

type pData struct {
	Callsign string      `json:"callsign"`
	Data     [][]float64 `json:"route"`
}

func getFileData(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("../track/f")
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	DataList := make([]pData, 0)

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()

			content, err := ioutil.ReadFile(filepath.Join("../track/f", filename))
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", filename, err)
				continue
			}
			//去除字符串中的.txt
			callsign := strings.TrimSuffix(filename, filepath.Ext(filename))
			posData := strings.Split(string(content), "\n")
			posData_all := make([][]float64, 0)
			//以,分割每一个posData，将其放入到posData中
			for _, line := range posData {

				parts := strings.Split(line, ",")

				posData_n := make([]float64, 0)
				for _, part := range parts {
					part, _ := strconv.ParseFloat(part, 64)
					posData_n = append(posData_n, part)
				}
				//判断如果posData_n内的元素个数大于等于2，则将其放入到posData_all中
				if len(posData_n) >= 2 {
					posData_all = append(posData_all, posData_n)
				}
			}
			//合并posData_all的相同元素
			for i := 0; i < len(posData_all)-1; i++ {
				if posData_all[i][0] == posData_all[i+1][0] && posData_all[i][1] == posData_all[i+1][1] {
					posData_all = append(posData_all[:i], posData_all[i+1:]...)
					i--
				}
			}

			Data := pData{
				Callsign: callsign,
				Data:     posData_all,
			}
			DataList = append(DataList, Data)
		}
	}
	jsonData, err := json.Marshal(DataList)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
