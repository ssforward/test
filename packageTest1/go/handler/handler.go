package handler

import (
	"net/http"
	"html/template"
)

func TopHandler(w http.ResponseWriter, r *http.Request){
	var templatefile = template.Must(template.ParseFiles("../packageTest1/layout/html/top.html"))
	templatefile.Execute(w, "top.html")
}
