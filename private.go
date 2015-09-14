package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"net/http"
	"os"
	"time"
)

func (s *Service) servePrivate(w http.ResponseWriter, r *http.Request) {
	if s.validatePrivateServer(r) {
		s.Handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}

func (s *Service) validatePrivateServer(r *http.Request) bool {
	key := r.Header.Get("PrivateServiceKey")
	if key != s.privateServiceKey && key != s.privateServiceKeyOld {
		return false
	}
	return true
}

func (s *Service) registerPrivateServer() {
	s.privateServiceKey = newPrivateServerKey()
	s.privateServiceKeyOld = newPrivateServerKey()
	s.follow()
	s.lead()
}

func (s *Service) follow() {
	watchPrivateKey(s.privateServiceKeyPath(), s.updatePrivateServerKeys)
}

func (s *Service) lead() {
	client := etcd.NewClient(Env.Machines())
	_, err := client.Create(s.privateServiceKeyPath(), newPrivateServerKey(), Frequency)
	defer client.Close()
	if err != nil { //TODO: check if err is key already set
		return
	}
	go func() {
		for _ = range time.Tick(time.Duration(Frequency-1) * time.Second) {
			_, err := client.Set(s.privateServiceKeyPath(), newPrivateServerKey(), Frequency)
			if err != nil {
				//TODO: handle error
				fmt.Fprint(os.Stderr, err)
			}
		}
	}()
}

func (s *Service) updatePrivateServerKeys(r *etcd.Response) {
	if r.Action == "expire" || r.Action == "delete" {
		s.lead()
		return
	}
	if r.PrevNode != nil {
		s.privateServiceKeyOld = r.PrevNode.Value
	}
	s.privateServiceKey = r.Node.Value
}

func (s *Service) watchPrivateClientKeys() {
	for _, r := range s.Requires {
		if r.Private {
			key := fmt.Sprintf("%v.%v", r.Title, r.Version)
			watchPrivateKey(key, s.updatePrivateClientKey)
		}
	}
}

func (s *Service) updatePrivateClientKey(r *etcd.Response) {
	s.privateClientKeys[r.Node.Key] = r.Node.Value
}

func watchPrivateKey(key string, set func(*etcd.Response)) {
	client := etcd.NewClient(Env.Machines())
	defer client.Close()
	wc := make(chan *etcd.Response)
	go client.Watch(key, 0, false, wc, nil)
	go func() {
		for r := range wc {
			set(r)
		}
	}()
}

func newPrivateServerKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		//TODO output err
		return newPrivateServerKey()
	}
	return base64.URLEncoding.EncodeToString(key)
}
