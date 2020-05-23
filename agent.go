package gomonit

import (
	"encoding/base64"
	"net/http"
	"net/http/cookiejar"
	"golang.org/x/net/publicsuffix"
	"strings"
	"net/url"
	"strconv"
	"errors"
)

// MonitAgent Agent class that can :
//	- retrieve status from Monit Daemon
// 	- Invoke Action on service (start / stop / monitor / unmonitor) 
type MonitAgent struct {
	Auth 		bool
	AuthString  string
	URL			string
	client 		*http.Client
	Status		*MonitStatus
}

// NewMonitAgent Create new MonitAgent instance and automatically try to connect to daemon and 
// retrieve current status.
func NewMonitAgent(URL string, AuthString string) (*MonitAgent, error) {
	// Prepare Cookie JAR
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	agent := MonitAgent{(AuthString != ""), AuthString, URL, &http.Client{Jar: jar}, nil}

	// Do a first request to init status
	err = agent.RequestStatus()

	return &agent, err
}

// requestStatus : internal function used to build a status
func (m *MonitAgent)requestStatus() (*MonitStatus, error) {
	// Send Request to Monit HTTPd
	req, err := http.NewRequest("GET", m.URL + "/_status?format=xml", nil)
	if m.Auth {
		req.Header.Add("Authorization","Basic " + base64.StdEncoding.EncodeToString([]byte(m.AuthString))) 
	}
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Unmarshal XML Response
	var mr MonitStatus
	err = mr.Load(resp.Body)
	if err != nil {
		return nil,err
	}

	return &mr,nil
}

// RequestStatus Requests Monit daemon to retrieve current status
func (m *MonitAgent)RequestStatus() error{
	var err error
	m.Status,err = m.requestStatus()
	return err
}

// DoAction Perform new action on Monit Daemon.
func (m *MonitAgent)doAction(service string, action string) error {
	u, _ := url.ParseRequestURI(m.URL)
    u.Path = "/_doaction"
    urlStr := u.String()

	data := url.Values{}
	data.Set("action", action)
	data.Set("service", service)
	data.Set("format", "text")

	// retreive security token
	var sec string 
	for _, cookie := range m.client.Jar.Cookies(u) {
		if cookie.Name == "securitytoken" {
			sec = cookie.Value
		}
	}
	data.Set("securitytoken", sec)

	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization","Basic " + base64.StdEncoding.EncodeToString([]byte(m.AuthString))) 
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//TODO do anything with resp content ? resp, error := m.client.Do(r)
	_, error := m.client.Do(r)
	
	return error
	//TODO exploit status code fmt.Println(resp.Status)
}

func (m *MonitAgent)CmdService(service string, action string) error {
	// Commond must exists
	if action != "start" &&
		action != "stop" &&
		action != "monitor" &&
		action != "unmonitor" {
			return errors.New("Unsupported action")
	}
	// Service Must exist
	if m.Status.GetService(service) == nil {
		return errors.New("Service doesn't exists")
	}

	return m.doAction(service, action)
}

func (m *MonitAgent)StartService(service string) error {
	// Service Must exist
	if m.Status.GetService(service) == nil {
		return errors.New("Service doesn't exists")
	}

	return m.doAction(service, "start")
}

