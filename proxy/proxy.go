package proxy

import (
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

var client = &http.Client{}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

type HttpHeader struct {
	Id    string
	Value string
}

type Proxy struct {
	Headers     []HttpHeader
	Destination *url.URL
}

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Println(p.Destination, " ", req.Method, " ", req.URL)

	req.RequestURI = ""
	req.URL.Scheme = p.Destination.Scheme
	req.URL.Host = p.Destination.Host

	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	log.Println(p.Destination, " ", resp.Status)

	copyHeader(wr.Header(), resp.Header)
	for _, header := range p.Headers {
		wr.Header().Add(header.Id, header.Value)
	}
	log.Println(p.Destination, " ", resp.Status)

	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
}
