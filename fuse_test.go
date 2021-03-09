package fuse

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestFusing(t *testing.T) {

	var fuseRequestNum int32
	var failRequestNum int32
	// 20秒内，允许失败5次，达到上限时，熔断10秒
	schme := NewFuse(5, 10, 20, 1)

	f := func(i int, wg *sync.WaitGroup) {
		defer wg.Done()
		var key = "/user/get-user-info/v2/"
		time.Sleep(time.Duration(time.Now().UnixNano()%1000) * time.Millisecond)
		if !schme.FuseOk(key) {
			if i == -1 {
				panic("-1的请求不应该熔断")
				t.Fail()
				return
			}
			fmt.Printf("操作【%d】因为【熔断】直接返回\n", i)
			atomic.AddInt32(&fuseRequestNum, 1)
			return
		}

		fmt.Printf("操作【%d】调用获取用户信息\n", i)

		schme.Fail(key)
		atomic.AddInt32(&failRequestNum, 1)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i ++ {
		wg.Add(1)
		go f(i, &wg)
	}

	wg.Add(1)
	go func() {
		time.Sleep(time.Duration(schme.last+1) * time.Second)
		f(-1, &wg)
	}()

	wg.Wait()

	fmt.Println("熔断失败个数:", fuseRequestNum)
	fmt.Println("请求失败个数", failRequestNum)
}
