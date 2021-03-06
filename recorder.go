package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	// "github.com/kr/pretty"
)

type Recoder struct {
	Name     string `json:"-"`
	Request  Inbound
	Response Outbound
	req      *http.Request  `json:"-"`
	resp     *http.Response `json:"-"`
}

type Inbound struct {
	URI      string `json:"-"`
	Host     string
	Path     string
	Method   string
	Body     []byte
	BodyText string
}

type Outbound struct {
	Status     string
	StatusCode int
	Body       []byte
	BodyText   string
}

func newRecoder(req *http.Request) Recoder {
	iBody, err := httputil.DumpRequest(req, true)
	fatal(err)

	unproxyURL(req)
	return Recoder{
		req: req,
		Request: Inbound{
			Host:     req.Host,
			Path:     req.URL.Path,
			Method:   req.Method,
			Body:     iBody,
			BodyText: byteToStr(iBody),
		},
	}
}

func (r Recoder) callService(t *Transport) (*http.Response, error) {
	if foundIncludeList(r) && isReplayMode() {
		cache := data.FindInCache(r)
		if cache != nil {
			fmt.Printf("CACHE_HIT  : current cache %v record, url=%v\n", len(data.List), r.req.RequestURI)
			return cache, nil
		}
		fmt.Printf("CACHE_MISS : current cache %v record, call http url=%v\n", len(data.List), r.req.RequestURI)
	} else {
		fmt.Printf("BY_PASS    : current cache %v record, call http url=%v\n", len(data.List), r.req.RequestURI)
	}

	resp, err := t.RoundTripper.RoundTrip(r.req)
	if err == nil {
		addToCache(r, resp)
	}
	return resp, err
}

func addToCache(row Recoder, resp *http.Response) {
	if foundIncludeList(row) && isRecordMode() {
		oBody, err := httputil.DumpResponse(resp, true)
		fatal(err)
		row.resp = resp
		row.Response.Status = resp.Status
		row.Response.StatusCode = resp.StatusCode
		row.Response.Body = oBody
		row.Response.BodyText = byteToStr(oBody)
		row.Name = generateKey(row.Request)

		data.List[row.Name] = row
		fmt.Printf("CACHE_ADDED: current cache %v record, cache_key=%v\n", len(data.List), row.Name)
	}
}
