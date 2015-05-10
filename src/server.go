package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"hunters"
)

var g *hunters.Game

func main() {

	rand.Seed(time.Now().UnixNano())
	restartGame()
	runServer()
}

func restartGame() {
	g = hunters.NewGame()
	g.LoadTestData()
	g.State = hunters.CombatStart
	g.Status <- "Game initialized."
	go func() {
		for g.State != hunters.End {
			g.ChangeState()
		}
		g.Status <- "Game over."
	}()

}

func runServer() {
	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.Handle("/js", http.FileServer(http.Dir("../client/js")))
	http.Handle("/templates", http.FileServer(http.Dir("../client/templates")))

	http.HandleFunc("/api/dump", apiDumpHandler)
	http.HandleFunc("/api/status", apiStatusHandler)
	http.HandleFunc("/api/prompt", apiPromptHandler)
	http.HandleFunc("/api/restart", apiRestartHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func apiDumpHandler(w http.ResponseWriter, r *http.Request) {
	g.Dump(w)
	g.Status <- "State dumped."
}

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(<-g.Status))
}

func apiPromptHandler(w http.ResponseWriter, r *http.Request) {
	p, err := json.Marshal(<-g.Output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(p)
}

func apiRestartHandler(w http.ResponseWriter, r *http.Request) {
	restartGame()
	w.WriteHeader(http.StatusOK)
}
