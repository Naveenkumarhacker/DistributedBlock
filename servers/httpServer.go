package servers

import (
	"DistributedBlock/constants"
	"DistributedBlock/dao"
	"DistributedBlock/models"
	"DistributedBlock/pb"
	"DistributedBlock/pkg/node"
	models2 "DistributedBlock/pkg/node/models"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func InitHttpServer(client pb.BlockServiceClient, nodeServer *node.NodeUtils) {
	r := gin.Default()

	// Create a new block
	r.POST("/blocks", func(c *gin.Context) {
		var input struct {
			Data string `json:"data"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		// Call gRPC service to create block
		res, err := client.CreateBlock(context.Background(), &pb.CreateBlockRequest{Data: input.Data})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		block := models.Block{
			ID:      uint(res.Block.Height),
			Parent:  res.Block.Parent,
			Encoded: res.Block.Encoded,
			Hash:    res.Block.Hash,
		}

		message := models2.NewNodeMessage(constants.BlockTopic, constants.Insert, block, true)
		nodeServer.BroadcastMessage(message)

		// Return response to user
		c.JSON(http.StatusOK, res.Block)
	})

	// Update an existing block
	r.PUT("/blocks/:height", func(c *gin.Context) {
		// Parse height from URL parameter
		heightStr := c.Param("height")
		height, err := strconv.ParseUint(heightStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid height"})
			return
		}

		// Parse input from request body
		var input struct {
			Data string `json:"data"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		// Call gRPC service to update block
		res, err := client.UpdateBlock(context.Background(), &pb.UpdateBlockRequest{Height: int32(height), Data: input.Data})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		block := models.Block{
			ID:      uint(res.Block.Height),
			Parent:  res.Block.Parent,
			Encoded: res.Block.Encoded,
			Hash:    res.Block.Hash,
		}

		message := models2.NewNodeMessage(constants.BlockTopic, constants.Update, block, true)
		nodeServer.BroadcastMessage(message)

		// Return response to user
		c.JSON(http.StatusOK, block)
	})

	// Get a block by height
	r.GET("/blocks/:height", func(c *gin.Context) {
		height, err := strconv.ParseUint(c.Param("height"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid height"})
			return
		}

		blockDao, err := dao.GetBlockByHeight(constants.DbMap, uint(height))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
			return
		}

		c.JSON(http.StatusOK, blockDao)
	})

	// Get blocks by minHeight
	r.GET("/blocks", func(c *gin.Context) {
		minHeight, err := strconv.ParseUint(c.Query("minHeight"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minHeight"})
			return
		}

		blocks, err := dao.GetBlocksByHeight(constants.DbMap, uint(minHeight))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, blocks)
	})

	r.Run(":" + constants.Port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", constants.Port),
		Handler: r,
	}

	// Start Gin router in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting Gin router: %v", err)
		}
	}()
}
