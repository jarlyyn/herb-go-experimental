package simpleapi

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

var listener net.Listener

var server *http.Server

var runningLock = sync.Mutex{}

var runningCount = int32(0)

var running = sync.Map{}

func Start(name string, h func(w http.ResponseWriter, r *http.Request)) error {
	runningLock.Lock()
	defer runningLock.Unlock()
	s, ok := running.LoadOrStore(name, h)
	if ok == true && s != nil {
		return fmt.Errorf("simple api :\" %s\" is already running", name)
	}
	r := atomic.LoadInt32(&runningCount)
	atomic.AddInt32(&runningCount, 1)
	if r == 0 {
		return startServer()
	}
	return nil
}

func Stop(name string) error {
	runningLock.Lock()
	defer runningLock.Unlock()
	s, ok := running.Load(name)
	if ok == false || s == nil {
		return fmt.Errorf("simple api :\" %s\" is not running", name)
	}
	running.Delete(name)
	r := atomic.AddInt32(&runningCount, -1)
	if r <= 0 {
		return stopServer()
	}
	return nil
}

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	h, _ := running.Load(r.URL.Path)
	if h == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	h.(func(w http.ResponseWriter, r *http.Request))(w, r)
})

var startServer = func() error {
	c := Config()
	server = c.Server()
	server.Handler = handler
	l, err := net.Listen(c.Net, c.Addr)
	if err != nil {
		return err
	}

	listener = l

	go func() {
		if !c.TLS {
			server.Serve(l)
		} else {
			server.ServeTLS(l, c.TLSCertPath, c.TLSKeyPath)
		}
	}()
	return nil
}

var stopServer = func() error {
	err := listener.Close()
	server.Close()
	return err
}

func Reset() {
	runningLock.Lock()
	defer runningLock.Unlock()
	running = sync.Map{}
	runningCount = 0
	listener = nil
	server = nil
}
