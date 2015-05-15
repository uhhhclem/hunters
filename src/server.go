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
	go func() {
		g.Status("Game initialized")
		for g.State != hunters.End {
			g.ChangeState()
		}
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

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func apiDumpHandler(w http.ResponseWriter, r *http.Request) {
	g.Dump(w)
}

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {
	o := <-g.StatusMessages
	if o == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	p, err := json.Marshal(o)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(p)
}

func apiPromptHandler(w http.ResponseWriter, r *http.Request) {
	o := <-g.PromptMessages
	if o == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	p, err := json.Marshal(o)
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
