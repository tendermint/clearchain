package main

import (
	"fmt"
	"html/template"
	"net/http"
	"github.com/tendermint/clearchain/types"
)

var templates = template.Must(template.ParseFiles("view.html"))

func main() {
	http.HandleFunc("/view/", viewHandler)

	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unimplemented request for %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var tx *types.Tx = nil
	//TODO: query to tendermint for balances.... 
	renderTemplate(w, "view", tx)
}

func renderTemplate(w http.ResponseWriter, tmpl string, tx *types.Tx) {
	err := templates.ExecuteTemplate(w, tmpl+".html", tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
