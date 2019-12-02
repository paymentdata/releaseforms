package util

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	//PDFConverterEndpoint points to the url which ought to route to the /convert endpoint for an instance of the AthenaPDF microservice.
	PDFConverterEndpoint = "http://athenapdf:8080/convert?auth=%s&ext=html"
	weaverAuthKey        = "arachnys-weaver"
)

//GetPDF is intended to receive a rendered HTML template payload, which is intended to be submitted to the athenapdf service.
//The func responds with the PDF response as []byte, and/or an error if there is one.
func GetPDF(data []byte) ([]byte, error) {
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
