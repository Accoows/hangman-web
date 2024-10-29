package main

import (
	"fmt"
	"hangman_web/Hangmanclassic"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var word string
var revealed []bool
var attemptsLeft int
var lettersRevealed map[rune]bool

func initializeGame() {
	// Initialisation du jeu
	words, err := Hangmanclassic.ReadFileLines("words.txt", false)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}
	word, revealed = Hangmanclassic.FindWord(words)
	attemptsLeft = Hangmanclassic.Maxtentative
	lettersRevealed = make(map[rune]bool)
}

func renderTemplate(w http.ResponseWriter, data interface{}) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Hangman Game</title>
	</head>
	<body>
		<h1>Jeu du Pendu</h1>
		<p>Mot à deviner : {{.WordDisplay}}</p>
		<p>Tentatives restantes : {{.AttemptsLeft}}</p>
		<form action="/hangman" method="POST">
			<label for="letter">Entrez une lettre :</label>
			<input type="text" name="letter" maxlength="1" required>
			<button type="submit">Vérifier</button>
		</form>
	</body>
	</html>
	`
	t, err := template.New("hangman").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	wordDisplay := Hangmanclassic.DisplayWord(word, revealed, ' ')
	data := struct {
		WordDisplay  string
		AttemptsLeft int
	}{
		WordDisplay:  wordDisplay,
		AttemptsLeft: attemptsLeft,
	}
	renderTemplate(w, data)
}

func Hangman(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	letter := r.FormValue("letter")
	if letter == "" || len(letter) != 1 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	revealedRune := rune(strings.ToUpper(letter)[0])
	if lettersRevealed[revealedRune] {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	lettersRevealed[revealedRune] = true

	found := false
	for i, ltr := range word {
		if ltr == revealedRune {
			revealed[i] = true
			found = true
		}
	}

	if !found {
		attemptsLeft--
	}

	if attemptsLeft <= 0 || Hangmanclassic.AllRevealed(revealed) {
		initializeGame() // Réinitialiser le jeu si toutes les lettres sont trouvées ou si les tentatives sont épuisées
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	initializeGame()

	http.HandleFunc("/", Home)
	http.HandleFunc("/hangman", Hangman)

	port := "8080"
	fmt.Printf("Serveur démarré sur http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Échec du démarrage du serveur : %v", err)
	}
}
