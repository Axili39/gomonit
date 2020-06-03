package gomonit

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var monit MonitStub

type MonitStub struct {
	replyStatus   bool
	replydoAction bool
}

func (m *MonitStub) Status(w http.ResponseWriter, r *http.Request) {
	if m.replyStatus {
		file, _ := ioutil.ReadFile("tests/status-test.xml")
		_, err := w.Write(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MonitStub) DoAction(w http.ResponseWriter, r *http.Request) {
	if m.replydoAction {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
func TestAgent(t *testing.T) {
	agent, err := NewMonitAgent("http://127.0.0.1:2812", "admin:monit")
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	// Request Status
	err = agent.RequestStatus()
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	// Do an action on service
	err = agent.StartService("service1")
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	// Do an action on unknown service
	err = agent.StartService("unknown")
	if err == nil {
		t.Errorf("%v\n", err)
		return
	}
	// Request Status error
	monit.replyStatus = false
	err = agent.RequestStatus()
	if err == nil {
		t.Errorf("%v\n", err)
		return
	}

	// Do an action on service error
	monit.replydoAction = false
	err = agent.StartService("service1")
	if err == nil {
		t.Errorf("%v\n", err)
		return
	}
}

func serve() {
	err := http.ListenAndServe("0.0.0.0:2812", nil)
	if err != nil {
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	monit = MonitStub{true, true}
	http.HandleFunc("/_status", monit.Status)
	http.HandleFunc("/_doaction", monit.DoAction)
	go serve()
	time.Sleep(1 * time.Second)

	os.Exit(m.Run())
}
