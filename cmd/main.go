package main

import (
	"flag"
	"github.com/gorilla/mux"
	"intro-rest/internal/app/handlers"
	"intro-rest/internal/app/repositories/entityrepo"
	"intro-rest/internal/app/repositories/memdatarepo"
	"intro-rest/internal/app/services/configmanager"
	"intro-rest/internal/app/services/dbclient"
	"intro-rest/internal/app/store"
	"log"
	"net/http"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.yml", "path to config file")
}

func main() {
	flag.Parse()
	config := configmanager.NewConfig()
	err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}
	db := dbclient.NewMongoDBClient(config)
	entityRepo := entityrepo.NewMongoEntityRepo(db, config)
	err = entityRepo.Init()
	if err != nil {
		log.Fatal(err)
	}
	memRepo := memdatarepo.NewMemDataRepo(entityRepo, config)
	err = memRepo.Init()
	if err != nil {
		log.Fatal(err)
	}
	mongoStore := store.NewStore(entityRepo, memRepo)

	entityCtrl := handlers.NewEntityController(mongoStore)
	memDataCtrl := handlers.NewMemDataController(mongoStore)

	router := mux.NewRouter()
	router.HandleFunc("/entity", entityCtrl.Update).Methods("PUT")
	router.HandleFunc("/entity/{id}", entityCtrl.Delete).Methods("DELETE")
	router.HandleFunc("/entities-db", entityCtrl.GetAll).Methods("GET")

	router.HandleFunc("/entities-memory", memDataCtrl.GetAll).Methods("GET")

	err = http.ListenAndServe(config.BindAddr, router)
	log.Fatal(err)
}
