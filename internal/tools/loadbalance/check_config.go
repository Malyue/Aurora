package loadbalance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
	DefaultWeight         = 50
)

type CheckConf struct {
	observers    []Observer
	confIpWeight map[string]string
	activeList   []string
	format       string
}

func NewCheckConf(format string, conf map[string]string) (*CheckConf, error) {
	list := make([]string, 0)
	for item, _ := range conf {
		list = append(list, item)
	}

	checkConf := &CheckConf{
		format:       format,
		activeList:   list,
		confIpWeight: conf,
	}

	return checkConf, nil
}

func (c *CheckConf) Attach(o Observer) {
	c.observers = append(c.observers, o)
}

func (c *CheckConf) NotifyAllObservers() {
	for _, obs := range c.observers {
		obs.Update()
	}
}

func (c *CheckConf) GetConf() []string {
	confList := make([]string, 0)
	for _, ip := range c.activeList {
		weight, ok := c.confIpWeight[ip]
		if !ok {
			weight = string(rune(DefaultWeight))
		}
		confList = append(confList, fmt.Sprintf(c.format, ip)+","+weight)
	}
	return confList
}

func (c *CheckConf) WatchConf() {
	go func() {
		confIpErrNum := map[string]int{}
		for {
			changedList := make([]string, 0)
			for item, _ := range c.confIpWeight {
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCheckTimeout)*time.Second)
				if err == nil {
					conn.Close()
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] = 0
					}
				}
				if err != nil {
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] += 1
					} else {
						confIpErrNum[item] = 1
					}
				}
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			sort.Strings(changedList)
			sort.Strings(c.activeList)
			if !reflect.DeepEqual(changedList, c.activeList) {
				c.UpdateConf(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

func (c *CheckConf) UpdateConf(conf []string) {
	c.activeList = conf
	for _, obs := range c.observers {
		obs.Update()
	}
}
