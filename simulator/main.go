package main

import (
	"flag"
	"sync"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	var wg sync.WaitGroup
	numberOfSwitch := 1
	//wg.Add(numberOfSwitch)
	glog.Infof("Start switch registration (%d switches)\n", numberOfSwitch)
	regist := make([]bool, numberOfSwitch)
	for i := 0; i < numberOfSwitch; i++ {
		wg.Add(1)
		regist[i] = singleSwitchRegistration(i, &wg)
		if regist[i] {
			glog.Infof("Registration succeded for Switch %d\n", i)
		} else {
			glog.Infof("Registration failed for Switch %d\n", i)
		}
	}
	wg.Wait()
	glog.Infof("Registration procedure all done\n")
	glog.Infof("Start sending UDP packets\n")
	for {

	}
}
