package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/paymentdata/releaseforms/form"
)

const (
	PDFConverterEndpoint = "http://127.0.0.1:8080/convert?auth=%s&ext=html"
	weaverAuthKey        = "weaver"
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

func getPDF(data []byte) ([]byte, error) {
	pdfEndpoint := fmt.Sprintf(PDFConverterEndpoint, weaverAuthKey)
	client := &http.Client{}

	// Create buffer
	buf := new(bytes.Buffer)

	r := bytes.NewReader(data)
	w := multipart.NewWriter(buf)

	// Create file field
	fw, err := w.CreateFormFile("file", "statementsummary")
	if err != nil {
		return nil, err
	}

	// Write file field from file to upload
	_, err = io.Copy(fw, r)
	if err != nil {
		return nil, err
	}
	// Important if you do not close the multipart writer you will not have a
	// terminating boundry
	w.Close()
	req, err := http.NewRequest("POST", pdfEndpoint, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	obuf := new(bytes.Buffer)
	obuf.ReadFrom(res.Body)

	return obuf.Bytes(), err
}
