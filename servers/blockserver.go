package servers

import (
	"DistributedBlock/models"
	"DistributedBlock/pb"
	"context"
	"gorm.io/gorm"
	"strconv"
)

type blockserver struct {
	pb.UnimplementedBlockServiceServer
	db *gorm.DB
}

func NewBlockServer(db *gorm.DB) (*blockserver, error) {
	return &blockserver{db: db}, nil
}

func (s *blockserver) CreateBlock(ctx context.Context, req *pb.CreateBlockRequest) (*pb.CreateBlockResponse, error) {
	// Create a new block record
	block := &models.Block{Hash: "random_hash", Encoded: []byte(req.Data)}

	// Get the ID of the last inserted block
	var lastBlock models.Block
	result := s.db.Last(&lastBlock)
	if result.Error != nil && result.RowsAffected != 0 {
		return nil, result.Error
	}

	// Set the Parent field to the ID of the last inserted block
	if result.RowsAffected != 0 {
		block.Parent = strconv.Itoa(int(lastBlock.ID))
	} else {
		block.Parent = "0" // Set parent to 0 if no previous block exists
	}

	// Save the new block record
	result = s.db.Create(block)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.CreateBlockResponse{Block: &pb.Block{Height: int32(block.ID), Hash: block.Hash, Parent: block.Parent, Encoded: block.Encoded}}, nil
}

func (s *blockserver) UpdateBlock(ctx context.Context, req *pb.UpdateBlockRequest) (*pb.UpdateBlockResponse, error) {
	// Update the block record
	var block models.Block
	result := s.db.First(&block, req.Height)
	if result.Error != nil {
		return nil, result.Error
	}

	block.Hash = "updated_hash"
	block.Encoded = []byte(req.Data)

	result = s.db.Save(&block)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.UpdateBlockResponse{Block: &pb.Block{Height: int32(block.ID), Hash: block.Hash, Parent: block.Parent, Encoded: block.Encoded}}, nil
}

func (s *blockserver) GetAllBlocks(ctx context.Context, req *pb.GetAllBlocksRequest) (*pb.GetAllBlocksResponse, error) {
	// Retrieve all blocks from the database
	var blocks []models.Block
	result := s.db.Find(&blocks)
	if result.Error != nil {
		return nil, result.Error
	}

	var pbBlocks []*pb.Block
	for _, b := range blocks {
		pbBlocks = append(pbBlocks, &pb.Block{Height: int32(b.ID), Hash: b.Hash, Parent: b.Parent, Encoded: b.Encoded})
	}

	return &pb.GetAllBlocksResponse{Blocks: pbBlocks}, nil
}
