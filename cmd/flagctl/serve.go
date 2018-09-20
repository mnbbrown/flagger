package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	flagger "github.com/mnbbrown/flagger/pkg"
	"github.com/spf13/cobra"
)

var redisClient *redis.Client

func getFlag(rw http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	environment := chi.URLParam(req, "environment")
	if name == "" || environment == "" {
		http.Error(rw, "name and environment required", http.StatusBadRequest)
		return
	}

	flag, err := flagger.GetFlag(redisClient, name, environment)
	if err != nil {
		switch err {
		case flagger.ErrFlagNotFound:
			http.Error(rw, "Not Found", http.StatusNotFound)
			return
		default:
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	v := strconv.FormatBool(flag.Value())
	rw.Write([]byte(v))
}

type flagRequest struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

func (f *flagRequest) Flag() (*flagger.Flag, error) {
	if f.Value == "true" {
		f.Value = "1"
	}
	if f.Value == "false" {
		f.Value = "0"
	}
	i, err := strconv.Atoi(f.Value)
	if err != nil {
		return nil, err
	}

	switch f.Type {
	case "bool", "BOOL":
		if i > 0 {
			i = 1
		}
		if i < 0 {
			i = -1
		}
		return &flagger.Flag{InternalValue: i, Type: flagger.BOOL}, nil
	case "percent", "PERCENT":
		if i > 100 {
			i = 100
		}
		if i < 100 {
			i = 0
		}
		return &flagger.Flag{InternalValue: i, Type: flagger.PERCENT}, nil
	}
	return nil, errors.New("flag type must be either bool or percent")
}

func setFlag(rw http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	environment := chi.URLParam(req, "environment")
	if name == "" || environment == "" {
		http.Error(rw, "name and environment required", http.StatusBadRequest)
		return
	}
	var f *flagRequest
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	if err := decoder.Decode(&f); err != nil {
		log.Println(err)
		http.Error(rw, "invalid JSON", http.StatusBadRequest)
		return
	}
	flag, err := f.Flag()
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Saving flag: %v", f)
	err = flagger.SaveFlag(redisClient, name, environment, flag)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte("OK"))
}

// Serve runs the HTTP server
func Serve() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Get("/flags/{name}/{environment}", getFlag)
	r.Post("/flags/{name}/{environment}", setFlag)
	log.Println("Listening on :8082")
	http.ListenAndServe(":8082", r)
}

// ServeCommand is a cobra command for serving the API
var serveCommand = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCommand)
}
