package main
import (
	"fmt"	
	"github.com/Axili39/gomonit"
	"os"
	"flag"
)

var monitAddr string
var monitAuth string

func commonFlags(flags *flag.FlagSet) {
	flags.StringVar(&monitAddr, "addr", os.Getenv("MONITCTL_ADDR"), "Monit Server Adress")
	flags.StringVar(&monitAuth, "auth", os.Getenv("MONITCTL_AUTH"), "Monit authentication string, ex:monit:admin")
}

func show(args []string) {
	// monitctl show [service] -s <MonitServer> -p <Auth>
	var showCli = flag.NewFlagSet("show", flag.ExitOnError)
	commonFlags(showCli)
	showCli.Parse(args)

	if monitAddr == "" {
		showCli.PrintDefaults()
		return
	}

	agent, err := gomonit.NewMonitAgent(monitAddr, monitAuth)
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}
	if showCli.Arg(0) != "" {
		// show specific service
		service := agent.Status.GetService(showCli.Arg(0))
		service.Print(os.Stdout)
	} else {
		// show All
		agent.Status.Print(os.Stdout)
	}
}

func cmd(args []string, action string) {
	// monitctl start/stop/unmonitor/unmonitor <service> -s <MonitServer> -p <Auth>
	var cmdCli = flag.NewFlagSet(action, flag.ExitOnError)
	commonFlags(cmdCli)
	cmdCli.Parse(args)

	if monitAddr == "" {
		cmdCli.PrintDefaults()
		return
	}
	if cmdCli.Arg(0) == "" {
		cmdCli.PrintDefaults()
		return
	}
	agent, err := gomonit.NewMonitAgent(monitAddr, monitAuth)
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}
	fmt.Printf("exec %s on service %s\n", action, cmdCli.Arg(0))
	err = agent.CmdService(cmdCli.Arg(0),action)
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}
}

func main() {
	switch os.Args[1] {
	case "show":
		show(os.Args[2:])
	case "start":
		cmd(os.Args[2:], "start")
	case "stop":
		cmd(os.Args[2:], "stop")
	case "monitor":
		cmd(os.Args[2:], "monitor")
	case "unmonitor":
		cmd(os.Args[2:], "unmonitor")
	default:
		fmt.Printf("%s show|start|stop|unmonitor\n",os.Args[0])
	}
}