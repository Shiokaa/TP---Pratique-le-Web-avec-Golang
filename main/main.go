package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

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

	})
	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {

	})

	http.ListenAndServe("localhost:8080", nil)

}
