package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"hunters"

	"github.com/uhhhclem/mse/src/interact"
)

var games map[string]*hunters.Game

func apiNewGame(w http.ResponseWriter, r *http.Request) {
	g := hunters.NewGame()
	go g.Run()

	if games == nil {
		games = make(map[string]*hunters.Game)
	}
	id := g.Game.ID
	games[id] = g

	resp := struct {
		ID string
	}{
		ID: id,
	}

	if b, err := json.Marshal(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

type apiGetHandler func(*hunters.Game, http.ResponseWriter, *http.Request) ([]byte, error)

func apiGetWrapper(h apiGetHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			status := http.StatusOK
			if err != nil {
				status = http.StatusInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			log.Printf("%d %s", status, r.URL)
		}()

		id := r.FormValue("ID")
		game := games[id]
		if game == nil {
			err = fmt.Errorf("Game ID %s not found.", id)
			return
		}

		var b []byte
		if b, err = h(game, w, r); err != nil {
			return
		}
		w.Write(b)
	}
}

func apiGetStatus(game *hunters.Game, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	s := <-game.NextStatus
	resp := struct {
		Status interact.Status
		End    bool
	}{
		End: s == nil,
	}
	if s != nil {
		resp.Status = *s
	}
	return json.Marshal(resp)
}

func apiGetBoard(game *hunters.Game, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	<-game.Ready
	b := struct {
		Boat   hunters.Boat
		Combat hunters.Combat
	}{
		Boat:   game.Boat,
		Combat: game.Combat,
	}
	return json.Marshal(b)
}

func apiGetPrompt(game *hunters.Game, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	p := <-game.NextPrompt
	resp := struct {
		Prompt interact.Prompt
		End    bool
	}{
		End: p == nil,
	}
	if p != nil {
		resp.Prompt = *p
	}

	return json.Marshal(resp)
}

func apiPostChoice(w http.ResponseWriter, r *http.Request) {
	var err error
	var id string
	var key string
	defer func() {
		if err == nil {
			log.Printf("%d %s id=%s key=%s", http.StatusOK, r.URL, id, key)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("%d %s %s", http.StatusInternalServerError, r.URL, err.Error())
		}
	}()

	req := struct {
		ID  string
		Key string
	}{}
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	id, key = req.ID, req.Key
	game, ok := games[id]
	if !ok {
		err = fmt.Errorf("Game %s not found.", id)
		return
	}

	err = game.MakeChoice(key)
}

func main() {
	http.HandleFunc("/api/newGame", apiNewGame)
	http.HandleFunc("/api/choice", apiPostChoice)

	handlers := []struct {
		url     string
		handler apiGetHandler
	}{
		{"/api/board", apiGetBoard},
		{"/api/status", apiGetStatus},
		{"/api/prompt", apiGetPrompt},
	}
	for _, h := range handlers {
		http.HandleFunc(h.url, apiGetWrapper(h.handler))
	}

	http.Handle("/", http.FileServer(http.Dir("./..")))

	http.ListenAndServe(":8080", nil)
}
