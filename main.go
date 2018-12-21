// idworker_srv project main.go
package main

import (
	"fmt"
	"time"

	// "fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/zheng-ji/goSnowFlake"
)

var idWorkerMap = make(map[int]*goSnowFlake.IdWorker)

const (
	SERVER_STATUS_OK    = 10000
	SERVER_STATUS_ERROR = 10001
)

func Test() {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
}

func main() {

	// Test()

	engine := gin.Default()

	engine.GET("/newid/:serverid", func(c *gin.Context) {
		serverid := c.Param("serverid")
		if len(serverid) == 0 {
			c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_ERROR, "m": "error", "d": "serverid error"})
			return
		}

		if m, _ := regexp.MatchString("^[0-9]+$", serverid); !m {
			c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_ERROR, "m": "error", "d": "serverid error"})
			return
		}

		id, err := strconv.Atoi(serverid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_ERROR, "m": "error", "d": "serverid error"})
			return
		}

		if id < 0 || id > 1023 {
			c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_ERROR, "m": "error", "d": "serverid error. 0 <= serverid < 1024"})
			return
		}

		idmap, ok := idWorkerMap[id]
		if ok {
			nid, err := idmap.NextId()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_ERROR, "m": "error", "d": "server error"})
				return
			}
			cur := time.Now()
			c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_OK, "m": "ok", "d": gin.H{
				"id": nid,
				"orderid": fmt.Sprintf("%d%d%d%d%d%d%d",
					cur.Year(), cur.Month(), cur.Day(), cur.Hour(), cur.Minute(), cur.Second(), nid),
			}})
		} else {
			idmap, err := goSnowFlake.NewIdWorker(int64(id))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_ERROR, "m": "error", "d": "server error"})
				return
			}
			idWorkerMap[id] = idmap
			nid, err := idmap.NextId()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_ERROR, "m": "error", "d": "server error"})
				return
			}
			cur := time.Now()
			c.JSON(http.StatusOK, gin.H{"c": SERVER_STATUS_OK, "m": "ok", "d": gin.H{
				"id": nid,
				"orderid": fmt.Sprintf("%d%d%d%d%d%d%d",
					cur.Year(), cur.Month(), cur.Day(), cur.Hour(), cur.Minute(), cur.Second(), nid),
			}})
		}
	})

	engine.Run(":9092")
}
