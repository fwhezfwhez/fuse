package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestGinFuse(t *testing.T) {
	go func() {
		r := gin.Default()
		// 加入熔断保障
		r.Use(GinHTTPFuse)
		r.GET("/", func(c *gin.Context) {
			c.JSON(500, gin.H{"message": "pretend hung up"})
		})
		r.Run(":8080")
	}()

	time.Sleep(3 * time.Second)

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(time.Now().UnixNano()%20) * time.Millisecond)
			defer wg.Done()
			rsp, e := http.Get("http://localhost:8080/")
			if e != nil {
				panic(e)
			}

			bdb, e := ioutil.ReadAll(rsp.Body)
			if e != nil {
				panic(e)
			}

			fmt.Println(rsp.StatusCode, string(bdb))
		}()
	}

	// after 10s, will recover recv 500
	time.Sleep(15 * time.Second)
	rsp, e := http.Get("http://localhost:8080/")
	if e != nil {
		panic(e)
	}

	bdb, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		panic(e)
	}

	fmt.Println(rsp.StatusCode, string(bdb))
	wg.Wait()

}
