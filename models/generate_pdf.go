package models

import (
	"bytes"
	"gade/srv-goldcard/logger"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// PDF RequestPdf struct
type RequestPdf struct {
	body string
}

// New request to pdf function
func (reqPdf *RequestPdf) NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

// ParseTemplate function to parse data to HTML file
func (reqPdf *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	template, err := template.ParseFiles(templateFileName)
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return err
	}

	buf := new(bytes.Buffer)
	if err = template.Execute(buf, data); err != nil {
		logger.Make(nil, nil).Debug(err)
		return err
	}

	reqPdf.body = buf.String()

	return nil
}

// GeneratePDF Function to generate pdf file
func (reqPdf *RequestPdf) GeneratePDF() ([]byte, error) {
	time := time.Now().Unix()
	uniqueNumber := strconv.FormatInt(int64(time), 10)
	// Write Temporary HTML File
	err := ioutil.WriteFile("template/"+uniqueNumber+".html", []byte(reqPdf.body), 0644)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return nil, err
	}

	// Open File HTML
	fileHtml, err := os.Open("template/" + uniqueNumber + ".html")
	if fileHtml != nil {
		defer fileHtml.Close()
	}

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return nil, err
	}

	// Generate New PDF File
	pdfg, err := wkhtmltopdf.NewPDFGenerator()

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return nil, err
	}

	// Add HTML File to PDF
	pdfg.AddPage(wkhtmltopdf.NewPageReader(fileHtml))

	// Set Page Size
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	// Set DPI
	pdfg.Dpi.Set(300)

	// Set Output Type to Buffer
	outBuf := new(bytes.Buffer)
	pdfg.SetOutput(outBuf)

	// Create PDF Buffer File
	err = pdfg.Create()
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return nil, err
	}

	bufPDF := outBuf.Bytes()

	// Delete Temporary HTML File
	err = os.Remove("template/" + uniqueNumber + ".html")
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return nil, err
	}

	return bufPDF, nil
}
