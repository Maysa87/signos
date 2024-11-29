package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

type Signo struct {
	DataInicio string `xml:"dataInicio"`
	DataFim    string `xml:"dataFim"`
	Nome       string `xml:"nome"`
	Descricao  string `xml:"descricao"`
}

type Signos struct {
	Signos []Signo `xml:"signo"`
}

func loadSignos() (Signos, error) {
	var signos Signos
	file, err := os.ReadFile("data/signos.xml")
	if err != nil {
		return signos, err
	}
	err = xml.Unmarshal(file, &signos)
	return signos, err
}

func getSigno(signos Signos, dataNascimento string) (Signo, error) {
	date, err := time.Parse("02/01", dataNascimento)
	if err != nil {
		return Signo{}, err
	}

	for _, signo := range signos.Signos {
		inicio, _ := time.Parse("02/01", signo.DataInicio)
		fim, _ := time.Parse("02/01", signo.DataFim)

		if (date.After(inicio) || date.Equal(inicio)) && (date.Before(fim) || date.Equal(fim)) {
			return signo, nil
		}
	}

	return Signo{}, fmt.Errorf("signo não encontrado")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		dataNascimento := strings.TrimSpace(r.FormValue("dataNascimento"))

		signos, err := loadSignos()
		if err != nil {
			http.Error(w, "Erro ao carregar os signos", http.StatusInternalServerError)
			return
		}

		signo, err := getSigno(signos, dataNascimento)
		if err != nil {
			http.Error(w, "Erro ao determinar o signo", http.StatusBadRequest)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/result.html"))
		tmpl.Execute(w, signo)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/result", resultHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		return
	}
	fmt.Println("Servidor iniciado em http://localhost:3000")
}
