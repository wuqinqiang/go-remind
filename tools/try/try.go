package try

import (
	"errors"
	"time"
)

// 最大允许重试次数
var MaxRetries = 3

var errMaxRetriesReached = errors.New("exceeded retry limit")

type Func func(attempt int) (retry bool, err error)

func Do(fn Func) error {
	var err error
	var cont bool
	attempt := 1
	for {
		if attempt > 1 {
			time.Sleep(2 * time.Second)
		}
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > MaxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}
