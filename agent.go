package gomonit

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// MonitAgent Agent class that can :
//	- retrieve status from Monit Daemon
// 	- Invoke Action on service (start / stop / monitor / unmonitor)
type MonitAgent struct {
	Auth       bool         // Authentication needed
	AuthString string       // Authentication info ex: admin:monit
	URL        string       // Monit httpd url, can be unix:///var/run/monit.sock, httpd://<host>:<port>
	client     *http.Client // internal : http client used to perform request to Monit
	Status     *MonitStatus // Last Status Obtained from Monit
}

// NewMonitAgent Create new MonitAgent instance and automatically try to connect to daemon and
// retrieve current status.
func NewMonitAgent(url string, authString string) (*MonitAgent, error) {
	var httpc http.Client
	var baseURL string

	// Prepare Cookie JAR, cookie management is mandatory to correctly pass the SecurityToken while performing request
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	// Init http client
	fields := strings.SplitN(url, "://", 2)
	if fields[0] == "unix" {
		// Transport is UNIX
		httpc = http.Client{
			Jar: jar,
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", fields[1])
				},
			},
		}
		baseURL = "http://unix"
	} else {
		// Transport is TCP
		baseURL = url
		httpc = http.Client{Jar: jar}
	}

	// Init Agent
	agent := MonitAgent{(authString != ""), authString, baseURL, &httpc, nil}

	// Do a first request to init status
	err = agent.RequestStatus()

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// requestStatus : internal function used to build a status
func (m *MonitAgent) requestStatus() (*MonitStatus, error) {
	// Send Request to Monit HTTPd
	req, err := http.NewRequest("GET", m.URL+"/_status?format=xml", nil)
	if err != nil {
		return nil, err
	}

	if m.Auth {
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(m.AuthString)))
	}
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal XML Response
	var mr MonitStatus
	err = mr.Load(resp.Body)
	if err != nil {
		return nil, err
	}

	return &mr, nil
}

// RequestStatus Requests Monit daemon to retrieve current status
func (m *MonitAgent) RequestStatus() error {
	var err error
	m.Status, err = m.requestStatus()
	return err
}

// DoAction Perform new action on Monit Daemon.
func (m *MonitAgent) doAction(service string, action string) error {
	u, _ := url.ParseRequestURI(m.URL)
	u.Path = "/_doaction"
	urlStr := u.String()

	data := url.Values{}
	data.Set("action", action)
	data.Set("service", service)
	data.Set("format", "text")

	// retrieve security token
	var sec string
	for _, cookie := range m.client.Jar.Cookies(u) {
		if cookie.Name == "securitytoken" {
			sec = cookie.Value
		}
	}
	data.Set("securitytoken", sec)

	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(m.AuthString)))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := m.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code %d", resp.StatusCode)
	}
	//TODO do anything with resp content ? resp, error := m.client.Do(r)
	//TODO exploit status code fmt.Println(resp.Status)
	return err
}

// CmdService Perform command on service
func (m *MonitAgent) CmdService(service string, action string) error {
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

// StartService perform monit start <service> command
func (m *MonitAgent) StartService(service string) error {
	// Service Must exist
	if m.Status == nil {
		return errors.New("Invalid current Monit Status")
	}
	// Service Must exist
	if m.Status.GetService(service) == nil {
		return errors.New("Service doesn't exists")
	}

	return m.doAction(service, "start")
}
