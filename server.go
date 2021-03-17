package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"market-backend/parser"
	"net/http"
	"os"
	"path/filepath"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func api(w http.ResponseWriter, r *http.Request) {
	items := parser.Search(r.URL.Query().Get("text"), r.URL.Query().Get("how"))
	if items == nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if itemsJSON, err := json.Marshal(&parser.TemplateJSON{Count: len(items), Items: items}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		_, _ = w.Write(itemsJSON)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/search", api)

	spa := spaHandler{staticPath: "public", indexPath: "index.html"}

	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}