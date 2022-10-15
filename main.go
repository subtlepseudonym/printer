package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	templatePath = "filament.tmpl"
	dataPath     = "filament.json"
	historyPath  = "history"
)

var page = template.Must(template.ParseFiles(templatePath))

type Data struct {
	Materials map[string]Material `json:"materials"`
}

type Material struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"` // in grams
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	dataFile, err := os.OpenFile(dataPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("ERR: parse data flat file: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dataFile.Close()

	var data Data
	if err = json.NewDecoder(dataFile).Decode(&data); err != nil {
		fmt.Printf("ERR: parse json: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		if err = r.ParseForm(); err != nil {
			fmt.Printf("ERR: parse form %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		history, err := os.OpenFile(historyPath, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
		if err != nil {
			fmt.Printf("ERR: open history %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		now := time.Now()
		for key, values := range r.PostForm {
			if len(values) < 1 || values[0] == "" {
				continue
			}

			if _, ok := data.Materials[key]; !ok {
				fmt.Printf("ERR: unknown material ID: %q\n", key)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			consumed, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				fmt.Printf("ERR: parse consumed amount: %q: %s", values[0], err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			material := data.Materials[key]
			material.Amount -= consumed
			data.Materials[key] = material

			out := fmt.Sprintf("%s %s %f\n", now.Format(time.RFC3339), key, consumed)
			_, err = history.WriteString(out)
			if err != nil {
				fmt.Printf("ERR: write history: %s: %s", out, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if err = dataFile.Truncate(0); err != nil {
			fmt.Printf("ERR: truncate data file: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = dataFile.Seek(0, 0); err != nil {
			fmt.Printf("ERR: data file seek: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		encoder := json.NewEncoder(dataFile)
		encoder.SetIndent("", "  ")
		if err = encoder.Encode(&data); err != nil {
			fmt.Printf("ERR: write json: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fallthrough
	case http.MethodGet:
		if err = page.Execute(w, &data); err != nil {
			fmt.Printf("ERR: execute template %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	srv := &http.Server{
		Addr:    ":9000",
		Handler: http.HandlerFunc(pageHandler),
	}
	srv.ListenAndServe()
}
