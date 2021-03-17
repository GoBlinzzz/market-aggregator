package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"market-backend/cart"
	"market-backend/parser"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

func search(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

func addToCart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newItem parser.Item
	_ = decoder.Decode(&newItem)
	var key string

	for _, c := range r.Cookies() {
		if c.Name == "cart_id" {
			key = c.Value
		}
	}
	if key == "" {
		key = fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprint(time.Now().Unix()))))
		http.SetCookie(w, &http.Cookie{
			Name: "cart_id",
			Value: key,
		})
	}
	cart.AddToCart(key,&newItem)
}

func cartApi(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var key string
	for _, c := range r.Cookies() {
		if c.Name == "cart_id" {
			key = c.Value
		}
	}
	if r.Method == "GET" {
		_, _ = w.Write(cart.GetCart(key))
	} else if r.Method == "DELETE" {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		cart.DeleteFromCart(key, id)
	} else {
		return
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/search", search).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/add-to-cart", addToCart).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cart", cartApi).Methods("GET", "DELETE", "OPTIONS")

	spa := spaHandler{staticPath: "public", indexPath: "index.html"}

	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}