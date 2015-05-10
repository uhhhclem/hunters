package main

import (
    "log"
	"math/rand"
    "net/http"
	"time"
	
	"hunters"

    
)

var g *hunters.Game

func main() {
    
	rand.Seed(time.Now().UnixNano())

	g = hunters.NewGame()
	g.LoadTestData()
	g.State = hunters.CombatStart
    g.Status <- "Game initialized."
    
    runServer()
}

func runServer() {
    http.Handle("/", http.FileServer(http.Dir("../client")))
    http.Handle("/js", http.FileServer(http.Dir("../client/js")))
    http.Handle("/templates", http.FileServer(http.Dir("../client/templates")))
    
    http.HandleFunc("/api/dump", apiDumpHandler)
    http.HandleFunc("/api/status", apiStatusHandler)
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func apiDumpHandler(w http.ResponseWriter, r *http.Request) {
    g.Dump(w)
    g.Status <- "State dumped."
}

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(<- g.Status))
}
