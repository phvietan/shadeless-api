package main

import (
	"log"
	"net/http"
	"os"
	"shadeless-api/main/config"
	"time"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health check ok"))
}

// func fileServerHandler(w http.ResponseWriter, r *http.Request) {
// 	ruri := r.RequestURI
// 	if tFile.MatchString(ruri) {
// 		w.Header().Set("Content-Type", "text/plain")
// 	}
// 	fileserver.ServeHTTP(w, r)
// }

func main() {
	// Create folder for serve static files
	dir := "files"
	_ = os.Mkdir(dir, 0755)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/files/").Handler(
		http.StripPrefix("/files/", http.FileServer(http.Dir(dir))),
	)
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         config.GetInstance().GetBindAddress(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
