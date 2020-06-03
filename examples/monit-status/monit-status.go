package main
import (
	"fmt"	
	"encoding/json"
	"github.com/Axili39/gomonit"
)

func main() {
	agent, err := gomonit.NewMonitAgent("http://192.168.1.19:2812", "admin:monit")
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}

	// Show current status Status
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