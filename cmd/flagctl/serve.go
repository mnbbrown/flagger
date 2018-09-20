package main

import (
	"encoding/json"
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
	f, ok := flag.Value()
	if !ok {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	v := strconv.FormatBool(f)
	rw.Write([]byte(v))
}

func setFlag(rw http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	environment := chi.URLParam(req, "environment")
	if name == "" || environment == "" {
		http.Error(rw, "name and environment required", http.StatusBadRequest)
		return
	}
	var f *flagger.Flag
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	if err := decoder.Decode(&f); err != nil {
		http.Error(rw, "invalid JSON", http.StatusBadRequest)
		return
	}
	err := flagger.SaveFlag(redisClient, name, environment, f)
	if err != nil {
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
