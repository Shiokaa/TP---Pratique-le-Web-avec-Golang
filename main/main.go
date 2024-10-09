package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

func main() {

	/*--------------------- DONNEE INIT ----------------------*/

	fileServer := http.FileServer(http.Dir("./assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println(fmt.Sprint("ERREUR => %v", err.Error()))
		os.Exit(02)
	}

	/*--------------------- PROMO PAGE ----------------------*/

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

	/*--------------------- CHANGE PAGE ----------------------*/

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

	/*--------------------- FORM PAGE ----------------------*/

	http.HandleFunc("/user/form", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "form", nil)
	})

	/*--------------------- TRAITEMENT PAGE ----------------------*/

	type StockageForm struct {
		CheckValue bool
		Surname    string
		Firstname  string
		Birth      string
		Gender     string
	}

	var stockageForm = StockageForm{false, "", "", "", ""}

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/erreur?code=400&message=Oups méthode incorrecte", http.StatusMovedPermanently)
			return
		}

		checkValueSurname, _ := regexp.MatchString("^[a-zA-Z]{1,32}$", r.FormValue("surname"))
		checkValueFirstname, _ := regexp.MatchString("^[a-zA-Z]{1,32}$", r.FormValue("firstname"))
		checkValueBirth, _ := regexp.MatchString("^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/([0-9]{4})$", r.FormValue("birth"))
		gender := r.FormValue("gender")

		validGender := []string{"Homme", "Femme", "Autre"}
		isValidGender := false

		for _, v := range validGender {
			if v == gender {
				isValidGender = true
				break
			}
		}

		if !isValidGender {
			http.Redirect(w, r, "/erreur?code=400&message=Oups genre invalide", http.StatusMovedPermanently)
			return
		}

		if !checkValueBirth {
			stockageForm = StockageForm{false, "", "", "", ""}
			http.Redirect(w, r, "/erreur?code=400&message=Oups la date de naissance n'est pas bonne", http.StatusMovedPermanently)
			return
		}

		if !checkValueSurname || !checkValueFirstname {
			stockageForm = StockageForm{false, "", "", "", ""}
			http.Redirect(w, r, "/erreur?code=400&message=Oups des données sont invalides", http.StatusMovedPermanently)
			return
		}

		stockageForm = StockageForm{true, r.FormValue("surname"), r.FormValue("firstname"), r.FormValue("birth"), gender}

		http.Redirect(w, r, "/user/display", http.StatusSeeOther)
	})

	/*--------------------- DISPLAY PAGE ----------------------*/

	type PageDisplay struct {
		CheckValue bool
		Surname    string
		Firstname  string
		Birth      string
		Gender     string
		IsEmpty    bool
	}

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		data := PageDisplay{stockageForm.CheckValue, stockageForm.Surname, stockageForm.Firstname, stockageForm.Birth, stockageForm.Gender, (!stockageForm.CheckValue && (stockageForm.Surname == "" || stockageForm.Firstname == "" || stockageForm.Birth == "" || stockageForm.Gender == ""))}

		temp.ExecuteTemplate(w, "formdisplay", data)
	})

	/*--------------------- ERREUR PAGE ----------------------*/

	http.HandleFunc("/erreur", func(w http.ResponseWriter, r *http.Request) {
		code, message := r.FormValue("code"), r.FormValue("message")
		if code != "" && message != "" {
			fmt.Fprintf(w, "Erreur %s - %s", code, message)
			return
		}
		fmt.Fprint(w, "Oups une erreur serveur est survenue")
	})

	/*--------------------- LANCEMENT DU SITE ----------------------*/

	http.ListenAndServe("localhost:8080", nil)

}
