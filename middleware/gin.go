package middleware

import (
	"fmt"
	"github.com/fwhezfwhez/fuse"
	"github.com/gin-gonic/gin"
	"time"
)

var fm = fuse.NewFuse(20, 10, 5, 128)

func ResetFm(fuseTimes int, last int, pern int, slotNum int) {
	fm = fuse.NewFuse(fuseTimes, last, pern, slotNum)
}

func GinHTTPFuse(c *gin.Context) {
	if ok := fm.FuseOk(c.FullPath()); !ok {
		c.AbortWithStatusJSON(400, gin.H{
			"tip": fmt.Sprintf("http api '%s' has be fused for setting {%d times/%ds} and will lasting for %d second to retry", c.FullPath(), fm.FuseTimes(), fm.Perns(), fm.Last()),
		})
		return
	}

	t0 := time.Now()

	c.Next()

	if time.Now().Sub(t0).Milliseconds() > 7000 {
		fm.Fail(c.FullPath())
		return
	}

	if c.Writer.Status() > 410 {
		fm.Fail(c.FullPath())
		return
	}
}
