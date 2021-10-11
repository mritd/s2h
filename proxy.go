package main

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"time"
)

func handleHTTP(w http.ResponseWriter, req *http.Request, dialer proxy.Dialer) {
	tp := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			select {
			case <-ctx.Done():
				return nil, errors.New("context canceled")
			default:
				return dialer.Dial(network, addr)
			}
		},
	}
	resp, err := tp.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func handleTunnel(w http.ResponseWriter, req *http.Request, dialer proxy.Dialer) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	srcConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	dstConn, err := dialer.Dial("tcp", req.Host)
	if err != nil {
		_ = srcConn.Close()
		return
	}

	_, _ = srcConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go transfer(dstConn, srcConn)
	go transfer(srcConn, dstConn)
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer func() {
		_ = dst.Close()
		_ = src.Close()
	}()

	_, _ = io.Copy(dst, src)
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("New Conn Accepted: %s => %s => [%s] %s", r.RemoteAddr, socket5Addr, r.Method, r.Host)
	d := &net.Dialer{Timeout: 10 * time.Second}
	dialer, err := proxy.SOCKS5("tcp", socket5Addr, nil, d)
	if err != nil {
		logrus.Errorf("Proxy Connect Failed: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if r.Method == "CONNECT" {
		handleTunnel(w, r, dialer)
	} else {
		handleHTTP(w, r, dialer)
	}
}
