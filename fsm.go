package fsm

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Even string
type State string
type actionFun func(s State, e Even, args ...interface{})

type Trans struct {
	FromState State
	ToState   State
	TriggEven Even
	Action    actionFun
}

type Processer interface {
	Processer(in chan interface{})
}

type defaultProcesser interface {
	Processer(in chan interface{})
}

type Fsm struct {
	Transfers map[State]map[Even]*Trans
}

var ErrIsExist = errors.New("the instance is exist")
var ErrIsNotExist = errors.New("the instance is'not exist")

func NewFsm() *Fsm {
	return &Fsm{
		Transfers: make(map[State]map[Even]*Trans),
	}
}

func (f *Fsm) RegisterAction(t *Trans) error {
	if _, ok := f.Transfers[t.FromState]; !ok {
		f.Transfers[t.FromState] = make(map[Even]*Trans)
	}

	stateTransfer := f.Transfers[t.FromState]
	if _, ok := stateTransfer[t.TriggEven]; ok {
		return ErrIsExist
	}

	stateTransfer[t.TriggEven] = t
	return nil
}

func (f *Fsm) getTransfer(s State, e Even) (*Trans, error) {
	if st, ok := f.Transfers[s]; ok {
		if trans, ok := st[e]; ok {
			return trans, nil
		}
	}

	return nil, ErrIsNotExist
}

func (f *Fsm) Trigger(s State, e Even, args ...interface{}) error {
	trans, err := f.getTransfer(s, e)
	if err != nil {
		return err
	}
	trans.Action(s, e, args...)
	return nil
}

//生成图表
func (f *Fsm) Export(outfile string) error {
	dot := `digraph StateMachine {

	rankdir=LR
	node[width=1 fixedsize=true shape=circle style=filled fillcolor="darkorchid1" ]
	
	`

	for _, v := range f.Transfers {
		for _, trans := range v {
			link := fmt.Sprintf(`%s -> %s [label="%s"]`, trans.FromState, trans.ToState, trans.TriggEven)
			dot = dot + "\r\n" + link
		}
	}

	dot = dot + "\r\n}"
	cmd := fmt.Sprintf("dot -o%s -T%s -K%s -s%s %s", outfile, "png", "dot", "72", "-Gsize=10,5 -Gdpi=200")
	return system(cmd, dot)
}

func system(c string, dot string) error {
	cmd := exec.Command(`/bin/sh`, `-c`, c)
	cmd.Stdin = strings.NewReader(dot)
	return cmd.Run()
}
