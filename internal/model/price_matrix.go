package model

import (
	"sync"
	"time"
)

// Define a struct to hold the ticket price matrix and a mutex for concurrency
type PriceMatrix struct {
	sizeRow    int
	sizeColumn int
	// Holds two adjacent matrix
	// Price matrix - [size][size]int64
	// Timestamp matrix - [size][size]uint64
	data []int64
	// cache offset of second matrix
	offset  int
	rwMutex *sync.RWMutex
}

// NewPriceMatrix initializes a new PriceMatrix with the given size
func NewPriceMatrix(sizeRow, sizeColumn int) *PriceMatrix {
	matrixSize := sizeRow * sizeColumn
	data := make([]int64, matrixSize*2)
	for i := 0; i < matrixSize; i++ {
		data[i] = 100          // Initialize all prices to 100
		data[i+matrixSize] = 0 // Initialize timestamps with 0
	}
	return &PriceMatrix{
		sizeRow:    sizeRow,
		sizeColumn: sizeColumn,
		offset:     matrixSize,
		data:       data,
		rwMutex:    &sync.RWMutex{},
	}
}

// Update price record in a thread-safe manner
func (pm *PriceMatrix) UpdatePrice(row, column int, price int64) {
	if row < 0 || row >= pm.sizeRow || column < 0 || column >= pm.sizeColumn {
		// Illegal operation
		return
	}
	// Collision between price and time doesn't matter if they happen in short time
	pm.rwMutex.RLock()
	defer pm.rwMutex.RUnlock()
	pm.data[row*pm.sizeColumn+column] = price
	pm.data[pm.offset+row*pm.sizeColumn+column] = time.Now().Unix() // Record the current time of update
}

// Expand size of the price matrix
func (pm *PriceMatrix) IncreaseDimension(rowCount, columnCount int) {
	if rowCount < 0 || columnCount < 0 || rowCount+columnCount < 1 {
		// Illegal operation
		return
	}

	newRowSize := pm.sizeRow + rowCount
	newColumnSize := pm.sizeColumn + columnCount
	newMatrixSize := newRowSize * newColumnSize
	// Create a new matrix to hold the data
	newData := make([]int64, newMatrixSize*2)

	// Initialize or copy from original matrix
	pm.rwMutex.RLock()
	for i := 0; i < newRowSize; i++ {
		for j := 0; j < newColumnSize; j++ {
			if i < pm.sizeRow && j < pm.sizeColumn {
				newData[i*newColumnSize+j] = pm.data[i*pm.sizeColumn+j]                         // Copy price
				newData[newMatrixSize+i*newColumnSize+j] = pm.data[pm.offset+i*pm.sizeColumn+j] // Copy timestamp
			} else {
				newData[i*newColumnSize+j] = 100             // Initialize new prices to 100
				newData[newMatrixSize+i*newColumnSize+j] = 0 // Initialize new timestamps with 0
			}
		}
	}
	defer pm.rwMutex.RUnlock()

	// Apply the change now
	pm.rwMutex.Lock()
	defer pm.rwMutex.Unlock()
	pm.sizeRow = newRowSize
	pm.sizeColumn = newColumnSize
	pm.offset = newMatrixSize
	pm.data = newData
}

func (pm PriceMatrix) GetData() []int64 {
	return pm.data
}
