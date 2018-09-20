package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

var flags map[string]map[string]string
var globalDefault = "on"

func loadFromFile() error {
	flags = make(map[string]map[string]string)
	f, err := ioutil.ReadFile("./flags.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, &flags)
	if err != nil {
		return err
	}
	log.Printf("Loaded %v flags", len(flags))
	return nil
}

func listFlags(rw http.ResponseWriter, req *http.Request) {
	b, _ := json.Marshal(flags)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

func getDefault(flag map[string]string) string {
	if flagDefault, ok := flag["default"]; ok {
		return flagDefault
	}
	return globalDefault
}

func getFlag(rw http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	environment := chi.URLParam(req, "environment")
	if name == "" || environment == "" {
		http.Error(rw, "name and environment required", http.StatusBadRequest)
		return
	}
	state := ""
	if flag, ok := flags[name]; ok {
		if state, ok = flag[environment]; !ok {
			state = getDefault(flag)
		}
	}
	if state == "" {
		http.Error(rw, "Flag not found", http.StatusNotFound)
		return
	}
	rw.Write([]byte(state))
}

func main() {
	if err := loadFromFile(); err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Get("/flags", listFlags)
	r.Get("/flags/{name}/{environment}", getFlag)
	http.ListenAndServe(":8082", r)
}
