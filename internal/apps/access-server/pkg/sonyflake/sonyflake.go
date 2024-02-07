package sonyflake

import (
	"github.com/sony/sonyflake"
	"strconv"
	"sync"
)

var (
	generator *sonyflake.Sonyflake
	mutex     sync.Mutex
)

func Init(machineId uint16) {
	mutex.Lock()
	defer mutex.Unlock()

	if generator == nil {
		generator = sonyflake.NewSonyflake(sonyflake.Settings{
			MachineID: func() (uint16, error) {
				return machineId, nil
			},
		})
	}

	return
}

func NextID() uint64 {
	id, _ := generator.NextID()
	return id
}

func NextStringID() string {
	return strconv.FormatUint(NextID(), 10)
}
