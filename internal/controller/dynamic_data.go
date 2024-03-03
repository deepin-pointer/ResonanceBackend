package controller

import "rsbackend/internal/model"

type DynamicData struct {
	data *model.PriceMatrix
}

func NewDynamicData(goodsCount, cityCount int) *DynamicData {
	return &DynamicData{
		data: model.NewPriceMatrix(goodsCount, cityCount),
	}
}

func (dd *DynamicData) AddCity(size int) {
	dd.data.IncreaseDimension(0, size)
}

func (dd *DynamicData) AddGoods(size int) {
	dd.data.IncreaseDimension(size, 0)
}

func (dd *DynamicData) ModifyPrice(record *model.PriceRecord) {
	dd.data.UpdatePrice(record.GoodsID, record.CityID, int64(record.Price))
}

func (dd *DynamicData) GetData() []int64 {
	return dd.data.GetData()
}
