package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/klokare/evo/x/web/bolt"
)

var (
	port = flag.Int("port", 2016, "port upon which the web server listens")
	path = flag.String("db-path", "evoweb.db", "path to the data store file")
)

var (
	client *bolt.Client
)

func main() {

	// Load the command line flags
	flag.Parse()

	// Open the data store
	client = &bolt.Client{Path: *path}
	if err := client.Open(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Set up the routes
	r := mux.NewRouter()
	r.HandleFunc("/studies/add", addStudy).Methods("POST", "GET")
	r.HandleFunc("/studies/edit", setStudy).Methods("POST")
	r.HandleFunc("/experiments/edit", setExperiment).Methods("POST")

	r.HandleFunc("/", getStudies).Methods("GET")
	r.HandleFunc("/studies", getStudies).Methods("GET")
	r.HandleFunc("/studies/{sid}", getExperiments).Methods("GET")
	r.HandleFunc("/studies/{sid}/experiments/{eid}", getTrials).Methods("GET")
	r.HandleFunc("/studies/{sid}/experiments/{eid}/trials/{tid}", getIterations).Methods("GET")

	r.HandleFunc("/api/studies/add", addStudyApi).Methods("POST")
	r.HandleFunc("/api/experiments/{eid}", getExperimentApi).Methods("GET")
	r.HandleFunc("/api/experiments/add", addExperimentApi).Methods("POST")
	r.HandleFunc("/api/experiments/{eid}/trials/add", addTrialApi).Methods("POST")
	r.HandleFunc("/api/trials/{tid}/iterations/add", addIterationApi).Methods("POST")

	// Run the server
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
