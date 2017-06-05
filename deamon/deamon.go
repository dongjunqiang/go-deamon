package deamon

import (
	"fmt"
	"time"
)

type QuitMsg struct{
	Pid int
	IsNormal bool
}

type MasterConf struct {
	ProcessCount int
	Msg chan QuitMsg
	SleepTime int
}

type ChildConf struct {
	ChildId int
	Parent *MasterConf
}

//接口，共child去重载实现
type ChildProcess interface {
	LoadConfig()
	Run(conf *ChildConf)

}


func Start(masterConf *MasterConf, childProcess ChildProcess)  {
	fmt.Println("master start")

	masterConf.Msg = make(chan QuitMsg)
	//开启了chan，一定不能忘记close掉
	defer func() { close(masterConf.Msg) }()
	//启动线程
	for i := 0; i < masterConf.ProcessCount; i++ {
		//需要把childProcess转成child来进行使用
		childConf := &ChildConf{ChildId: i, Parent: masterConf}
		go masterConf.CommonHandle(childConf, childProcess)
	}

	//需要判断子进程如果失败了，需要重新启动，即对子进程的监督，所以需要死for循环
	for {
		quitMsg := <-masterConf.Msg
		if quitMsg.IsNormal {
			fmt.Println("child exit normal", quitMsg.Pid)
		}else {
			fmt.Println("child exit unnormal", quitMsg.Pid)
		}

		fmt.Println("重新启动restart", quitMsg.Pid)
		childConf := &ChildConf{ChildId: quitMsg.Pid, Parent: masterConf}
		go masterConf.CommonHandle(childConf, childProcess)

	}



}

func (master *MasterConf) CommonHandle(childConf *ChildConf, childProcess ChildProcess)()  {
	//子进程是否安全结束,发送msg信号
	//子进程也是需要常驻内存的，需要死for循环，这个需要

	defer func() {
		retNo := recover()

		if retNo != nil {
			childConf.Parent.Msg <- QuitMsg{childConf.ChildId, false}
		} else {
			childConf.Parent.Msg <- QuitMsg{childConf.ChildId, true}
		}
	}()
	for {
		childProcess.Run(childConf)
		time.Sleep(time.Second * time.Duration(childConf.Parent.SleepTime))
	}

}

