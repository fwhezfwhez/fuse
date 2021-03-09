package fuse

import (
	"fmt"
	"github.com/fwhezfwhez/cmap"
	"time"
)

type Fuse struct {
	m *cmap.MapV2

	fuseTimes int
	last      int // second
	perns     int // second
}

func NewFuse(fuseTimes int, last int, perns int, slotNum int) Fuse {
	return Fuse{
		m:         cmap.NewMapV2(nil, slotNum, 30*time.Minute),
		fuseTimes: fuseTimes,
		last:      last,
		perns:     perns,
	}
}

func (f *Fuse) FuseTimes() int {
	return f.fuseTimes
}
func (f *Fuse) Last() int {
	return f.last
}
func (f *Fuse) Perns() int {
	return f.perns
}

// true, 未熔断，放行
// false, 熔断态，禁止通行
func (f *Fuse) FuseOk(key string) bool {
	fuseKey := fmt.Sprintf("is_fused:%s", key)
	v, exist := f.m.Get(fuseKey)

	if !exist {
		return true
	}

	vs, ok := v.(string)
	if exist && ok && vs == "fused" {
		return false
	}
	return false
}

// 某一次请求失败了，则需要调用Fail()
// 当fail次数达到阈值时，将会使得f.FuseOK(conn ,key) 返回false，调用方借此来熔断操作
func (f *Fuse) Fail(key string) {

	multi := time.Now().Unix() / int64(f.perns)

	timeskey := fmt.Sprintf("%s:%d", key, multi)

	rs := f.m.IncrByEx(timeskey, 1, f.perns)

	var ok bool
	ok = rs <= int64(f.fuseTimes)

	// 未达到配置的熔断阈值，fail无操作
	if ok {
		return
	}

	// 达到了熔断点
	fuseKey := fmt.Sprintf("is_fused:%s", key)
	f.m.SetEx(fuseKey, "fused", f.last)
}
