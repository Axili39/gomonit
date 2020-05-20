GoMonit Golang API for Monit
============================

Simple golang agent API for Monit.
Can be used to perform some simple action on Monit daemon :
- retreive current status,
- start a service,
- stop a service
- activate / deactivate monitoring on a service

*usage*
```
go get http://github.com/Axili39/gomonit
```

*see examples* 
```
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

	// Show current status Status
	out,_ := json.Marshal(&agent.Status)
	fmt.Printf("%s\n",string(out))

	// Do an action on service
	err = agent.DoAction("foo", "start")
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