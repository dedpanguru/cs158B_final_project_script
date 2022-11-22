package main

import (
	"fmt"  // needed for printing to STDERR
	"os"   // needed for properly stopping a program
	"sync" // needed for managing goroutines

	"github.com/reiver/go-telnet" // 3rd party library that makes dealing with TELNET easier, referenced as 'telnet' in the code
)

// these can be configured based on the network configuration
var (
	HOST_IP     = "192.168.1.61"   // host's IPv4 address
	NUM_ROUTERS = 10               // 10 routers total
	TELNET_IP   = "192.168.56.101" // routers' telnet address
	// set of commands as a string array
	// this script will direct all system logs to the host ip
	COMMANDS = []string{
		"config t",
		"logging " + HOST_IP,
		"end",
		"wr mem",
	}
)

func main() {
	// manage concurrency and errors
	var (
		wg   sync.WaitGroup             // tracks goroutine (thread) activity via a counter
		errs []error        = []error{} // collects errors
	)

	for i := 1; i <= NUM_ROUTERS; i++ {
		// generate telnet address
		targetAddress := fmt.Sprintf(TELNET_IP+":%d", 5000+i)
		// fire off a goroutine
		wg.Add(1)
		go func(address string) {
			// pass commands, collect any errors
			if err := PassCommands(address, COMMANDS...); err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}(targetAddress)
	}
	// wait for all goroutines to stop
	wg.Wait()
	// output any errors
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// PassCommands - establishes a telnet connection and passes commands to it
func PassCommands(address string, commands ...string) error {
	// make connection to telnet address
	conn, err := telnet.DialTo(address)
	// close connection at end of function
	defer func() {
		_ = conn.Close()
	}()
	// return any error in connecting
	if err != nil {
		return err
	}
	// print successful connection confirmation
	fmt.Println("Connected to", address)
	// loop through commands
	for _, command := range commands {
		// pass in each command and emulate an enter key hit
		_, err := conn.Write([]byte(command + "\r\n"))
		// handle any error in passing the command
		if err != nil {
			return err
		}
	}
	// print successful message transmission confirmation
	fmt.Println("Successfully passed commands to ", address)
	return nil
}
