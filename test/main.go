package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type StockageForm struct {
	CheckValue bool
	Value      string
}

var stockageForm = StockageForm{false, ""}

func main() {

	temp, err := template.ParseGlob("./*.html")
	if err != nil {
		fmt.Println(fmt.Sprint("ERREUR => %v", err.Error()))
		os.Exit(02)
	}

	http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "form", nil)
	})

	http.HandleFunc("/form/traitement", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/erreur?code=400&message=Oups méthode incorrecte", http.StatusMovedPermanently)
			return
		}

		checkValue, _ := regexp.MatchString("^[a-zA-Z-]{1,32}$", r.FormValue("name"))

		if !checkValue {
			stockageForm = StockageForm{false, ""}
			http.Redirect(w, r, "/erreur?code=400&message=Oups les données sont invalides", http.StatusMovedPermanently)
			return
		}
		stockageForm = StockageForm{true, r.FormValue("name")}

		http.Redirect(w, r, "/formtest", http.StatusSeeOther)
	})

	type PageDisplay struct {
		CheckValue bool
		Value      string
		IsEmpty    bool
	}

	http.HandleFunc("/formtest", func(w http.ResponseWriter, r *http.Request) {
		data := PageDisplay{stockageForm.CheckValue, stockageForm.Value, (!stockageForm.CheckValue && stockageForm.Value == "")}
		temp.ExecuteTemplate(w, "formtest", data)
	})

	http.HandleFunc("/erreur", func(w http.ResponseWriter, r *http.Request) {
		code, message := r.FormValue("code"), r.FormValue("message")
		if code != "" && message != "" {
			fmt.Fprintf(w, "Erreur %s - %s", code, message)
			return
		}
		fmt.Fprint(w, "Oups une erreur serveur est survenue")
	})

	http.ListenAndServe("localhost:8080", nil)

}
