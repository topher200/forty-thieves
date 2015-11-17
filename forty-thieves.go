package main

import (
	"fmt"
	"net/http"
)

func handleResources(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Providing", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func showHttp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Root handling request:", r.URL.Path[1:])
	http.ServeFile(w, r, "res/cards.html")
}

func main() {
	fmt.Println(NewGame())
	http.HandleFunc("/res/", handleResources)
	http.HandleFunc("/", showHttp)
	// http.ListenAndServe(":8080", nil)
}
