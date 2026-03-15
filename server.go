package main

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Batch struct {
	ID         uint32 `json:"id" binding:"required"`
	ClickCount int32  `json:"clickcount"`
}

var magicNumbers = map[int32]bool{
	676767:  true, // reset
	676766:  true, // multiply 2
	767676:  true, // multiply 5
	6767678: true, // divide 2
	7676768: true, // divide 5
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

		// Allow magic numbers through before cheat detection
		if !magicNumbers[newBatch.ClickCount] {
			if newBatch.ClickCount > 500 || newBatch.ClickCount == 0 {
				c.JSON(http.StatusAccepted, gin.H{
					"received_id":  newBatch.ID,
					"status":       "didnt enqueue due to possible cheating or 0 clicks",
					"global_total": atomic.LoadUint32(&globalNumber),
				})
				return
			}
		}

		// Add to queue
		q.Enqueue(newBatch.ClickCount)

		// Spin up background processing (non-blocking)
		go processQueue()

		c.IndentedJSON(http.StatusCreated, gin.H{
			"received_id":  newBatch.ID,
			"status":       "batch enqueued",
			"global_total": atomic.LoadUint32(&globalNumber),
		})
		return
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

// package main

// import (
// 	"net/http"
// 	"sync/atomic"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// )

// type Batch struct {
// 	ID         uint32 `json:"id" binding:"required"`
// 	ClickCount int32  `json:"clickcount"`
// }

// func main() {
// 	r := gin.Default()
// 	r.Use(cors.Default())

// 	// ENDPOINT 1: Add to the total
// 	r.POST("/process-batch", func(c *gin.Context) {
// 		var newBatch Batch
// 		if err := c.BindJSON(&newBatch); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// if a single server sends more than 500 clicks in a second then something is clearly wrong, we'll still send status 200 though
// 		if newBatch.ClickCount > 500 || newBatch.ClickCount == 0 {
// 			c.JSON(http.StatusAccepted, gin.H{
// 				"received_id":  newBatch.ID,
// 				"status":       "didnt enqueue due to possible cheating or 0 clicks",
// 				"global_total": atomic.LoadUint32(&globalNumber),
// 			})
// 			return
// 		}

// 		// 1. add to queue
// 		q.Enqueue(newBatch.ClickCount)

// 		// 2. spin up background processing (Non-blocking)
// 		go processQueue()

// 		c.IndentedJSON(http.StatusCreated, gin.H{
// 			"received_id":  newBatch.ID,
// 			"status":       "batch enqueued",
// 			"global_total": atomic.LoadUint32(&globalNumber),
// 		})
// 		return
// 	})

// 	// ENDPOINT 2: Just read the total
// 	r.GET("/global-number", func(c *gin.Context) {
// 		currentValue := atomic.LoadUint32(&globalNumber)

// 		c.IndentedJSON(http.StatusOK, gin.H{
// 			"global_total": currentValue,
// 		})
// 	})

// 	r.Run(":8080")
// }
