package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/paymentdata/releaseforms/form"
)

func main() {
	http.HandleFunc("/releaseForm", HandleReleaseEvent)

	e := http.ListenAndServe("0.0.0.0:3331", nil)
	if e != nil {
		fmt.Println(e.Error())
	}
}

func HandleReleaseEvent(w http.ResponseWriter, r *http.Request) {
	var (
		e   error
		rtd form.ReleaseTemplateData
	)

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&rtd)
	if err != nil {
		panic(err)
	}

	log.Printf("Decoded: %+v", rtd)
	t, e := template.New("releaseForm").Parse(form.ReleaseTemplate)
	if e != nil {
		fmt.Printf("Err %v", e)
	}

	e = t.Execute(w, rtd)

	if e != nil {
		log.Printf("TEMPLATE ERR: %v\n", e)
	}

}
