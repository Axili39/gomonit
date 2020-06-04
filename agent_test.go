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
	lastAction    string
	lastService   string
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
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		r.ParseForm()
		m.lastAction = r.FormValue("action")
		m.lastService = r.FormValue("service")
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
	if monit.lastAction != "start" {
		t.Errorf("last action request mismatch %s\n", monit.lastAction)
		return
	}
	if monit.lastService != "service1" {
		t.Errorf("last service request mismatch %s\n", monit.lastService)
		return
	}
	// 	CmdService
	err = agent.CmdService("service1", "stop")
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	if monit.lastAction != "stop" {
		t.Errorf("last action request mismatch %s\n", monit.lastAction)
		return
	}
	if monit.lastService != "service1" {
		t.Errorf("last service request mismatch %s\n", monit.lastService)
		return
	}
}
func TestAgentErrors(t *testing.T) {
	agent, err := NewMonitAgent("http://127.0.0.1:2812", "admin:monit")
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	// unknown command
	err = agent.CmdService("service1", "foo")
	if err == nil {
		t.Errorf("%v\n", err)
		return
	}
	// unknown service
	err = agent.CmdService("bar", "monitor")
	if err == nil {
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
	monit = MonitStub{true, true, "", ""}
	http.HandleFunc("/_status", monit.Status)
	http.HandleFunc("/_doaction", monit.DoAction)
	go serve()
	time.Sleep(1 * time.Second)

	os.Exit(m.Run())
}
