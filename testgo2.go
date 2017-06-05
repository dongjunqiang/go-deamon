package main

import (
	"fmt"
	"./deamon"
)

type SumChild struct {

}

func (c *SumChild) LoadConfig()  {
	fmt.Println("conf ---")
	return
}
func (c *SumChild) Run(conf *deamon.ChildConf)  {
	var sum int
	for i :=0 ; i < 50; i++  {
		sum += i
	}
	fmt.Println("sum ---", conf.ChildId, sum)
	return
}
func main() {
	masterConf := &deamon.MasterConf{ProcessCount: 2, SleepTime: 2}

	deamon.Start(masterConf, &SumChild{})

}