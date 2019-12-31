package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var dummy = http.Response{
	Body: ioutil.NopCloser(bytes.NewBufferString("")),
}

func newRequest(method, uri, body string) *http.Request {
	request, _ := http.NewRequest(method, uri, strings.NewReader(body))
	return request
}

func do(dr *http.Request) io.ReadCloser {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.DefaultClient.Do(dr)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Warn(fmt.Sprintf("Timeout on calling URL %s", dr.URL))
			return dummy.Body
		} else {
			log.Fatalln(err)
		}
	}
	if resp.StatusCode != http.StatusOK {
		log.Warn(fmt.Sprintf("Failed to call URL %s - status code was %d", dr.URL, resp.StatusCode))
		return dummy.Body
	}
	return resp.Body
}
