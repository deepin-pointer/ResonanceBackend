package controller

import (
	"encoding/json"
	"os"
	"rsbackend/internal/model"
	"sync"

	"github.com/gofiber/fiber/v2/log"
)

type StaticData struct {
	cache     []byte
	filePath  string
	CityList  []*model.City  `json:"city_list"`
	GoodsList []*model.Goods `json:"goods_list"`
	rwMutex   *sync.RWMutex
}

func NewStaticData(path string) *StaticData {
	sd := &StaticData{
		filePath: path,
		rwMutex:  &sync.RWMutex{},
	}
	err := sd.loadData(path)
	if err != nil {
		sd.cacheData()
	}
	return sd
}

// NewCity adds a new city to the CityList.
func (sd *StaticData) NewCity(data *model.City) error {
	if len(data.Distance) != len(sd.CityList)+1 {
		log.Error("Invalid City Data!")
		return os.ErrInvalid
	}
	sd.rwMutex.Lock()
	defer sd.rwMutex.Unlock()
	for i := range sd.CityList {
		sd.CityList[i].Distance = append(sd.CityList[i].Distance, data.Distance[i])
	}
	sd.CityList = append(sd.CityList, data)
	return sd.cacheData()
}

// ModifyCity updates an existing city in the CityList at the given index.
func (sd *StaticData) ModifyCity(index int, data *model.City) error {
	sd.rwMutex.Lock()
	defer sd.rwMutex.Unlock()
	sd.CityList[index] = data
	return sd.cacheData()
}

// NewGoods adds a new goods item to the GoodsList.
func (sd *StaticData) NewGoods(data *model.Goods) error {
	sd.rwMutex.Lock()
	defer sd.rwMutex.Unlock()
	sd.GoodsList = append(sd.GoodsList, data)
	return sd.cacheData()
}

// ModifyGoods updates an existing goods item at the given index.
func (sd *StaticData) ModifyGoods(index int, data *model.Goods) error {
	sd.rwMutex.Lock()
	defer sd.rwMutex.Unlock()
	sd.GoodsList[index] = data
	return sd.cacheData()
}

// Save static data to disk and cache for faster serving
func (sd *StaticData) cacheData() error {
	data, err := json.Marshal(sd)
	if err != nil {
		return err
	}
	sd.cache = data
	err = os.WriteFile(sd.filePath, data, 0644)
	return err
}

// Cached data to send
func (sd *StaticData) GetData() []byte {
	sd.rwMutex.RLock()
	defer sd.rwMutex.RUnlock()
	return sd.cache
}

// LoadData reads a JSON file from the specified path and deserializes it into the StaticData.
func (sd *StaticData) loadData(path string) error {
	sd.rwMutex.Lock()
	defer sd.rwMutex.Unlock()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, sd)
	if err != nil {
		return err
	}
	sd.cache = data
	return nil
}
