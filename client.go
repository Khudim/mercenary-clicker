package main

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"mime/multipart"
)


type point struct {
	Confidence float32
	X          int
	Y          int
}

type MatcherResponse struct {
	TemplateId string `json:"templateId"`
}

func upload(file, url string) string {
	var strRequestURI = url + "/template/upload"

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile(file, file)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	_, _ = part.Write(b)

	_ = writer.Close()

	if buf.Len() < 70 {
		panic("No templates were found.")
	}

	req.Header.SetMethodBytes([]byte("POST"))
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.SetRequestURI(strRequestURI)
	req.SetBody(buf.Bytes())

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		panic("handle error")
	}

	if res.StatusCode() != 200 {
		panic(res.Body())
	}
	var response MatcherResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		panic(err)
	}
	log.Println(response)

	return response.TemplateId
}

func detect(image []byte, t *Template) *point {
	return detectWithConf(image, t, 0.85)
}

func detectWithConf(image []byte, t *Template, conf float32) *point {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	url := matcherUrl + "/template/detect/" + t.id
	req.SetRequestURI(url)
	req.Header.SetMethodBytes([]byte("POST"))
	req.AppendBody(image)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		log.Fatal(err)
		return nil
	}

	var response point
	if err := json.Unmarshal(res.Body(), &response); err == nil {
		if t.Debug {
			log.Printf("%+v\n", response)
			log.Printf("%+v\n", t)
		}
		if response.Confidence >= conf {
			return &response
		}
	} else {
		log.Fatal(err)
	}

	return nil
}
