package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/mnbbrown/flagger/pkg"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

var redisHost string
var redisClient *redis.Client
var servePort int
var serveUI bool
var filesDir string

var f flagger.Flagger

func getFlag(rw http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	environment := chi.URLParam(req, "environment")
	if name == "" || environment == "" {
		http.Error(rw, "name and environment required", http.StatusBadRequest)
		return
	}

	flag, err := f.GetFlagWithTags(name, []string{environment})
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
		if i < 0 {
			i = 0
		}
		return &flagger.Flag{InternalValue: i, Type: flagger.PERCENT}, nil
	}
	return nil, errors.New("flag type must be either bool or percent")
}

func listFlags(rw http.ResponseWriter, req *http.Request) {
	flags, err := f.ListFlags()
	if err != nil {
		log.Println(err)
		http.Error(rw, "InternalServerError", http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(flags)
	if err != nil {
		log.Println(err)
		http.Error(rw, "InternalServerError", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

func setFlag(rw http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	environment := chi.URLParam(req, "environment")
	if name == "" || environment == "" {
		http.Error(rw, "name and environment required", http.StatusBadRequest)
		return
	}
	var fr *flagRequest
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	if err := decoder.Decode(&fr); err != nil {
		log.Println(err)
		http.Error(rw, "invalid JSON", http.StatusBadRequest)
		return
	}
	flag, err := fr.Flag()
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Saving flag: %v", f)
	err = f.SaveFlag(flag)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte("OK"))
}

func uiServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// Serve runs the HTTP server
func Serve() {
	var err error
	if f, err = flagger.NewRedisFlagger(redisHost, 0); err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Get("/flags", listFlags)
	r.Get("/flags/{name}/{environment}", getFlag)
	r.Post("/flags/{name}/{environment}", setFlag)
	if serveUI {
		log.Printf("Serving ui from %s", filesDir)
		uiServer(r, "/ui", http.Dir(filesDir))
	}
	bindAddr := fmt.Sprintf(":%v", servePort)
	log.Printf("Listening on %s", bindAddr)
	http.ListenAndServe(bindAddr, cors.Default().Handler(r))
}

// ServeCommand is a cobra command for serving the API
var serveCommand = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		Serve()
	},
}

func init() {
	serveCommand.Flags().StringVarP(&redisHost, "redis", "", "localhost:6379", "redis server address")
	serveCommand.Flags().IntVarP(&servePort, "port", "p", 8082, "http api port")
	serveCommand.Flags().BoolVarP(&serveUI, "ui", "", true, "serve the UI on /ui")
	serveCommand.Flags().StringVarP(&filesDir, "uiDir", "", "/ui", "Directory from which to serve the UI")
	rootCmd.AddCommand(serveCommand)
}
