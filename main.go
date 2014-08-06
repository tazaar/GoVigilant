package main

import "fmt"
import "flag"
import "time"

// Try connecting to host, report reachability and time
func main() {
	// Define and parse flags
	hostAddr := flag.String("hostAddr", nil, "Host address")
	flag.Parse()
	// try connecting
	if hostAddr != nil {
		fmt.Println("Attempting to reach site")
		t0 := time.Now()
		c, err := net.Dial("tcp", *hostAddr)
	    if err != nil {
	    	t1 := time.Now()
	        fmt.Println("Host unreachable, time: %v", t1.Sub(t0))
	        return
	    }
	    else {
	    	t1 := time.Now()
	    	fmt.Println("Host reachable, time: %v", t1.Sub(t0))
	    	c.Close()
	    }
	}
	else {
		fmt.Println("Host address empty")
	}
}