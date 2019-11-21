package main

import (
	"encoding/json"
	"fmt"
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
	go listener.Listen()

	//Prepare for Form generation requests
	http.HandleFunc("/releaseForm", RenderReleaseForm)

	e := http.ListenAndServe(os.Getenv("FormHost"), nil)
	if e != nil {
		fmt.Println(e.Error())
	}
}

func RenderReleaseForm(w http.ResponseWriter, r *http.Request) {
	var (
		rtd form.ReleaseTemplateData
	)

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&rtd)
	if err != nil {
		panic(err)
	}
	pdfResponse, err := util.GetPDF(rtd.Render())
	if err != nil {
		panic(err)
	}
	log.Printf("Decoded: %+v", rtd)

	n, e := w.Write(pdfResponse)
	if e != nil {
		panic(e)
	}
	fmt.Printf("Wrote %d bytes!", n)
}
