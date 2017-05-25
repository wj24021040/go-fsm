package main

import (
	"fmt"
	"time"

	. "github.com/wj24021040/go-fsm"
)

type testInstance struct {
	curstate State
	name     string
}

var test1 = &testInstance{
	curstate: "idle",
	name:     "test",
}

func main() {
	f := NewFsm()
	idle2init := &Trans{
		FromState: "idle",
		ToState:   "init",
		TriggEven: "start",
		Action:    Start,
	}

	init2run := &Trans{
		FromState: "init",
		ToState:   "run",
		TriggEven: "complete",
		Action:    complete,
	}

	run2idle := &Trans{
		FromState: "run",
		ToState:   "idle",
		TriggEven: "stop",
		Action:    stop,
	}
	f.RegisterAction(idle2init)
	f.RegisterAction(init2run)
	f.RegisterAction(run2idle)

	f.Trigger(test1.curstate, "start", test1)
	time.Sleep(5 * time.Second)
	f.Trigger(test1.curstate, "complete", test1)
	time.Sleep(5 * time.Second)
	f.Trigger(test1.curstate, "stop", test1)

	f.Export("./statemechine.png")
}

func Start(s State, e Even, args ...interface{}) {
	argslist := []interface{}(args)
	if ins, ok := argslist[0].(*testInstance); ok {
		fmt.Println(ins.name, "into stateMechine: ", s, " by ", e, " to ", "init")
		ins.curstate = "init"
	}
}

func complete(s State, e Even, args ...interface{}) {
	argslist := []interface{}(args)
	if ins, ok := argslist[0].(*testInstance); ok {
		fmt.Println(ins.name, "into stateMechine: ", s, " by ", e, " to ", "run")
		ins.curstate = "idle"
	}
}

func stop(s State, e Even, args ...interface{}) {
	argslist := []interface{}(args)
	if ins, ok := argslist[0].(*testInstance); ok {
		fmt.Println(ins.name, "into stateMechine: ", s, " by ", e, " to ", "idle")
		ins.curstate = "idle"
	}
}
