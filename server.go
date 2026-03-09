package main

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Batch struct {
	ID         uint32 `json:"id" binding:"required"`
	ClickCount uint32 `json:"clickcount"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// ENDPOINT 1: Add to the total
	r.POST("/process-batch", func(c *gin.Context) {
		var newBatch Batch
		if err := c.BindJSON(&newBatch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// if a single server sends more than 500 clicks in a second then something is clearly wrong
		if newBatch.ClickCount == 0 || newBatch.ClickCount > 500 {
			c.JSON(http.StatusLoopDetected, gin.H{"error": "someone may be cheating"})
		}
		// 1. add to queue
		q.Enqueue(newBatch.ClickCount)

		// 2. spin up background processing (Non-blocking)
		go processQueue()

		c.IndentedJSON(http.StatusCreated, gin.H{
			"received_id":  newBatch.ID,
			"status":       "batch enqueued",
			"global_total": atomic.LoadUint32(&globalNumber),
		})
	})

	// ENDPOINT 2: Just read the total
	r.GET("/global-number", func(c *gin.Context) {
		currentValue := atomic.LoadUint32(&globalNumber)

		c.IndentedJSON(http.StatusOK, gin.H{
			"global_total": currentValue,
		})
	})

	r.Run(":8080")
}
