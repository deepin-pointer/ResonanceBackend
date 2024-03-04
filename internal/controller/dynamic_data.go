package controller

import (
	"encoding/binary"
	"os"
	"rsbackend/internal/model"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type DynamicData struct {
	data    *model.PriceMatrix
	logChan chan model.PriceRecord
}

func NewDynamicData(goodsCount, cityCount int) *DynamicData {
	return &DynamicData{
		data:    model.NewPriceMatrix(goodsCount, cityCount),
		logChan: make(chan model.PriceRecord, 100),
	}
}

func (dd *DynamicData) AddCity(size int) {
	dd.data.IncreaseDimension(0, size)
}

func (dd *DynamicData) AddGoods(size int) {
	dd.data.IncreaseDimension(size, 0)
}

func (dd *DynamicData) ModifyPrice(records *[]model.PriceRecord) {
	dd.data.UpdatePrice(*records)
	for _, record := range *records {
		dd.logChan <- record
	}
}

func (dd *DynamicData) GetData() []int64 {
	return dd.data.GetData()
}

func (dd *DynamicData) LoggingWorker(path string) {
	bytes := make([]byte, 8*4)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Error("Failed to create log file:", err)
		return
	}
	for record := range dd.logChan {
		binary.LittleEndian.PutUint64(bytes, uint64(record.CityID))
		binary.LittleEndian.PutUint64(bytes[8:], uint64(record.GoodsID))
		binary.LittleEndian.PutUint64(bytes[16:], uint64(record.Price))
		binary.LittleEndian.PutUint64(bytes[24:], uint64(time.Now().Unix()))
		file.Write(bytes)
	}
}
