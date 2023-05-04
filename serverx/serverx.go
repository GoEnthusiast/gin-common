package serverx

import (
	"fmt"
	"net/http"
	"time"
)

func NewServer(host string, port, readTimeout, writeTimeout int) *http.Server {
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", host, port),
		ReadTimeout:    time.Duration(readTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(writeTimeout) * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}
	return s
}
