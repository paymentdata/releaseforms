package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/paymentdata/releaseforms/util"

	"github.com/paymentdata/releaseforms/listener"

	"github.com/paymentdata/releaseforms/form"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	//Listen for WebHooks
	listener.Listen()

	//Prepare for Form generation requests
	http.HandleFunc("/releaseForm", RenderReleaseForm)

	e := http.ListenAndServe(os.Getenv("FormHost"), nil)
	if e != nil {
		fmt.Println(e.Error())
	}
}

func RenderReleaseForm(w http.ResponseWriter, r *http.Request) {
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
	util.GetPDF(Render(rtd))

}

func Render(rtd form.ReleaseTemplateData) []byte {
	var (
		tpl bytes.Buffer
		e   error
	)
	t, e := template.New("thing").Parse(form.ReleaseTemplate)
	if e != nil {
		fmt.Printf("Err %v", e)
	}
	e = t.Execute(&tpl, rtd)
	if e != nil {
		log.Printf("TEMPLATE ERR: %v\n", e)
	}
	return tpl.Bytes()
}
