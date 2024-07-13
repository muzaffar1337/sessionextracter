package main

import (
	"crypto/tls"
	"time"

	"github.com/valyala/fasthttp"
)

func fastInfo() *fasthttp.Request {
	var req *fasthttp.Request = fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://i.instagram.com/api/v1/accounts/current_user/?edit=true")
	req.Header.Set("User-Agent", "Instagram 275.0.0.27.98 Android")
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	return req
}

func ClientHistoryP() (*fasthttp.Client, error) {
	client := &fasthttp.Client{
		TLSConfig:           &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnDuration: 5 * time.Second,
		MaxConnsPerHost:     1000,
		ReadBufferSize:      8192,
	}
	return client, nil
}
