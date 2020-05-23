package gomonit

import (
	"encoding/xml"
	"encoding/json"
	"golang.org/x/net/html/charset"
	"io"	
	"os"
	"fmt"	
)

// MonitStatus This structure embedded all data returned by monit daemon
type MonitStatus struct {
	XMLName xml.Name `xml:"monit" json:"-"`
	Server 			Server 		`xml:"server" json:"server"`
	Platform 		Platform 	`xml:"platform" json:"platform"`
	Services 		[]Service 	`xml:"service" json:"service"`
}

// Server "server" part of MonitStatus
type Server struct {
	Hostname 		string 	`xml:"localhostname" json:"localhostname"`
	ControlFile 	string 	`xml:"controlfile" json:"controlfile"`
	HTTPd struct {
		Address 		string 	`xml:"address" json:"address"`
		Port 			int 	`xml:"port" json:"port"`
		SSL 			int 	`xml:"ssl" json:"ssl"`
	} 						`xml:"httpd" json:"httpd"`
	Incarnation 	int64	`xml:"incarnation" json:"incarnation"`
	Version 		string 	`xml:"version" json:"version"`
	Uptime 			int 	`xml:"uptime" json:"uptime"`
	ID 				string 	`xml:"id" json:"id"`
	Poll 			int 	`xml:"poll" json:"poll"`
	StartDelay 		int 	`xml:"startdelay" json:"startdelay"`
}

// Platform "platform" part of MonitStatus
type Platform struct {
	Memory 			int64 `xml:"memory" json:"memory"`
	Swap			int64	`xml:"swap" json:"swap"`
	Name			string	`xml:"name" json:"name"`
	Release			string	`xml:"release" json:"release"`
	Version			string	`xml:"version" json:"version"`
	Machine			string	`xml:"machine" json:"machine"`
	CPU				int		`xml:"cpu" json:"cpu"`
	}

// Service service entry which contains service status 
type Service struct {
	Type 			int 	`xml:"type,attr" json:"type"` 
	Name 			string 	`xml:"name" json:"name"` 
	CollectedUsec 	int 	`xml:"collected_usec" json:"collected_usec"`
	Status 			int 	`xml:"status" json:"status"`
	MonitorMode 	int		`xml:"monitormode" json:"monitormode"`
	Pid 			int		`xml:"pid" json:"pid"`
	PPid			*int	`xml:"ppid" json:"ppid,omitempty"`
	UID 			*int	`xml:"uid" json:"uid,omitempty"`
	EUID			*int	`xml:"euid" json:"euid,omitempty"`
	GID				*int	`xml:"gid,omitempty" json:"gid,omitempty"`
	Uptime 			int 	`xml:"uptime" json:"uptime"`
	Threads			int		`xml:"threads" json:"threads"`
	Children		int		`xml:"children" json:"children"`
	StatusHint		int		`xml:"status_hint" json:"status_hint"`
	OnReboot		int		`xml:"onreboot" json:"onreboot"`
	CollectedSec	int		`xml:"collected_sec" json:"collected_sec"`
	Monitor			int		`xml:"monitor" json:"monitor"`
	PendingAction	int		`xml:"pendingaction" json:"pendingaction"`
	CPU *struct {
		Percent			float32	`xml:"percent" json:"percent"`
		PercentTotal 	float32	`xml:"percenttotal" json:"percenttotal"`
	} 						`xml:"cpu,omitempty" json:"cpu,omitempty"`
	//	TODO		<read></read>			<write></write>
	Memory	*struct {
		Percent			float32	`xml:"percent" json:"percent"`
		PercentTotal 	float32	`xml:"percenttotal" json:"percenttotal"`
		KiloBytes		int64	`xml:"kilobyte" json:"kilobyte"`
		KiloBytesTotal	int64	`xml:"kilobytetotal" json:"kilobytetotal"`
	}						`xml:"memory" json:"memory,omitempty"`
	System 	*struct {
		Swap struct {
			Percent 	float32	`xml:"percent" json:"percent"`
			KiloBytes	int64	`xml:"kilobyte" json:"kilobyte"`
		}			`xml:"swap" json:"swap"`
		Load struct {
			Avg01		float32	`xml:"avg01" json:"avg01"`
			Avg05		float32	`xml:"avg05" json:"avg05"`
			Avg15		float32	`xml:"avg15" json:"avg15"`
		} `xml:"load" json:"load"`
		CPU struct {
			User		float32	`xml:"user" json:"user"`
			System		float32	`xml:"system" json:"system"`
			Wait		float32	`xml:"wait" json:"wait"`
		} `xml:"cpu" json:"cpu"`
		Memory struct {
			Percent 	float32	`xml:"percent" json:"percent"`
			KiloBytes	int64	`xml:"kilobyte" json:"kilobyte"`
		}	`xml:"memory" json:"memory"`
	} 						`xml:"system" json:"system,omitempty"`
	Link *struct {
		State 		int			`xml:"state" json:"state"`
		Speed		int64		`xml:"speed" json:"speed"`
		Duplex		int			`xml:"duplex" json:"duplex"`
		Download	*NetStats	`xml:"download" json:"download"`
		Upload		*NetStats	`xml:"upload" json:"upload"`
	}	`xml:"link" json:"link,omitempty"`
}

// NetStats Link network statistics
type NetStats struct {
	Packet	*NetStatElem `xml:"packets" json:"packets"`
	Bytes	*NetStatElem `xml:"bytes" json:"bytes"`	
	Errors	*NetStatElem `xml:"errors" json:"errors"`	
}
// NetStatElem parts of NetStats				
type NetStatElem struct {
	Now		int64	`xml:"now" json:"now"`
	Total	int64	`xml:"total" json:"total"`
}

// Load : Load MonitStatus from io.Reader
func (r *MonitStatus)Load(file io.Reader) error {
	decoder := xml.NewDecoder(file)
    decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(r)
	return err
}

// 
func (r *MonitStatus)GetService(name string) *Service {
	for _,s := range r.Services {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func (r *MonitStatus)Print(w *os.File) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err == nil {
			fmt.Fprintln(w, string(b))
	}
}
func (s *Service)Print(w *os.File) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err == nil {
			fmt.Fprintln(w, string(b))
	}
}
