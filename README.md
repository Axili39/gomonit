![Build & Test](https://github.com/Axili39/gomonit/workflows/Build%20&%20Test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/Axili39/gomonit/badge.svg?branch=master)](https://coveralls.io/github/Axili39/gomonit?branch=master)

GoMonit Golang API for Monit
============================

*Monit is a free open source utility for managing and monitoring, processes, programs, files, directories and filesystems on a UNIX system. Monit conducts automatic maintenance and repair and can execute meaningful causal actions in error situations.*

**Monit Source Code** https://bitbucket.org/tildeslash/monit/src/master/

**Monit Tildeslash official website** https://mmonit.com/monit/

Simple golang agent API for Monit.
Can be used to perform some simple action on Monit daemon :
- retrieve current status,
- start a service,
- stop a service
- activate / deactivate monitoring on a service

*usage*
```
go get http://github.com/Axili39/gomonit
```

*prerequisites*

Activate http in /etc/monit/monitrc config file :
- specify a bind address,
- specify admin/password credential

*see /etc/monit/monitrc config file instructions, for more informations*

*Example*
```
set httpd port 2812 and
     use address 127.0.0.1
     allow admin:monit
```
*see examples* 
```go
package main
import (
	"fmt"	
	"encoding/json"
	"github.com/Axili39/gomonit"
)

func main() {
	agent, err := gomonit.NewMonitAgent("http://127.0.0.1:2812", "admin:monit")
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}

	// Show current status
	out,_ := json.Marshal(&agent.Status)
	fmt.Printf("%s\n",string(out))

	// Do an action on service
	err = agent.StartService("foo")
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}

	// Request Status
	err = agent.RequestStatus()
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}

	out,_ = json.Marshal(&agent.Status)
	fmt.Printf("%s\n",string(out))

}
```
