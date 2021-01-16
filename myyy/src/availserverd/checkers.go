package main

import (
	"fmt"
	"github.com/ChenLong-dev/gobase/mlog"
	"myyy/src/availserverd/alarm"
	"sync"
	"time"
)

type Checkers struct {
	sync.RWMutex
	m map[string]*Checker
}

func (checkers *Checkers) Add(checker *Checker) {
	mlog.Infof("checker=%v", checker)

	checkers.Lock()
	defer checkers.Unlock()

	if oldChecker, ok := checkers.m[checker.Region]; ok {
		oldChecker.Close()
	}
	checkers.m[checker.Region] = checker
}
func (checkers *Checkers) Del(region string) {
	mlog.Infof("del checker(%s)", region)
	//检测节点断开告警
	content := fmt.Sprintf("[检测节点异常告警] [%s] [%s] [%+v]", region, "检测节点异常或重启", time.Now().Format("2006-01-02 15:04:05"))
	alarm.SendAlarmMsg(region, content)

	checkers.Lock()
	defer checkers.Unlock()

	if oldChecker, ok := checkers.m[region]; ok {
		oldChecker.Close()
	}
	delete(checkers.m, region)
}
func (checkers *Checkers) All() (cs []*Checker) {
	checkers.RLock()
	defer checkers.RUnlock()

	for _, ch := range checkers.m {
		cs = append(cs, ch)
	}

	return cs
}

var defaultCheckers *Checkers

func InitChecks() {
	defaultCheckers = &Checkers{m: make(map[string]*Checker)}
}
