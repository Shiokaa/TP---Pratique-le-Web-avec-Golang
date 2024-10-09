package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

type StockageForm struct {
	CheckValue bool
	Value      string
}

var stockageForm = StockageForm{false, ""}

func main() {
	fileServer := http.FileServer(http.Dir("./assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println(fmt.Sprint("ERREUR => %v", err.Error()))
		os.Exit(02)
	}

	type Etudiant struct {
		Nom    string
		Prenom string
		Age    int
		Sexe   string
	}

	type PagePromo struct {
		NomDeClasse  string
		Filiere      string
		Niveau       string
		NbrEtudiant  int
		ListEtudiant []Etudiant
	}

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		etudiant1 := Etudiant{Nom: "Amaru", Prenom: "Tom", Age: 18, Sexe: "Homme"}
		etudiant2 := Etudiant{Nom: "Jean", Prenom: "Paul", Age: 21, Sexe: "Femme"}

		classe := PagePromo{NomDeClasse: "B1 Informatique", Filiere: "Informatique", Niveau: "B1", NbrEtudiant: 2,
			ListEtudiant: []Etudiant{etudiant1, etudiant2}}

		temp.ExecuteTemplate(w, "promo", classe)
	})

	type PageChange struct {
		IsPair  bool
		Counter int
	}

	var Counter int

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		Counter++
		data := PageChange{Counter: Counter, IsPair: Counter%2 == 0}

		temp.ExecuteTemplate(w, "change", data)
	})

	http.HandleFunc("/user/form", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "form", nil)
	})

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/erreur?code=400&message=Oups méthode incorrecte", http.StatusMovedPermanently)
			return
		}

		checkValue, _ := regexp.MatchString("^[a-zA-Z]{1,32}$", r.FormValue("surname"))

		if !checkValue {
			stockageForm = StockageForm{false, ""}
			http.Redirect(w, r, "/erreur?code=400&message=Oups les données sont invalides", http.StatusMovedPermanently)
			return
		}
		stockageForm = StockageForm{true, r.FormValue("surname")}

		http.Redirect(w, r, "/user/display", http.StatusSeeOther)
	})

	type PageDisplay struct {
		CheckValue bool
		Value      string
		IsEmpty    bool
	}

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		data := PageDisplay{stockageForm.CheckValue, stockageForm.Value, (!stockageForm.CheckValue && stockageForm.Value == "")}
		temp.ExecuteTemplate(w, "formdisplay", data)
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
