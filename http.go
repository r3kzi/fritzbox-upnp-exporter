package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	dac "github.com/xinsnake/go-http-digest-auth-client"
)

var dummy = http.Response{
	Body: ioutil.NopCloser(bytes.NewBufferString("")),
}

func newRequest(username, password, method, uri, body string) dac.DigestRequest {
	dr := dac.NewRequest(username, password, method, uri, body)
	dr.CertVal = false
	return dr
}

func do(dr dac.DigestRequest) io.ReadCloser {
	resp, err := dr.Execute()
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Warn(fmt.Sprintf("Timeout on calling URL %s", dr.URI))
			return dummy.Body
		} else {
			log.Fatalln(err)
		}
	}
	if resp.StatusCode != http.StatusOK {
		log.Warn(fmt.Sprintf("Failed to call URL %s - status code was %d", dr.URI, resp.StatusCode))
		return dummy.Body
	}
	return resp.Body
}
