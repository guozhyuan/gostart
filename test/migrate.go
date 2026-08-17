package test

import (
	"encoding/json"
	"fmt"
	"gostart/internal/config"
	"gostart/internal/pkg"
	"io"
	"net/http"
)

type Platform struct {
	ID   uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func (Platform) TableName() string {
	return "platform"
}

type Result struct {
	Success string
	Data    []Platform
	Total   int
}

// 根据model生成表

func insertPlatform() {
	config.ReadConfigByViper()
	pkg.ConnectDB()
	// 生成表
	// DB.AutoMigrate(&model.UserDO{})
	// service.DB.AutoMigrate(&Platform{})

	pkg.ReadConfig()
	pkg.ConnectDB()
	var ret Result
	// var platforms []Platform = make([]Platform, 100)
	resp, err := http.Get("https://live.suancaihu.eu.org/api/platforms")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bytes, &ret)
	fmt.Println(ret)
	fmt.Println("-----------------------------------------------")
	err2 := pkg.DB.Create(ret.Data).Error
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	fmt.Println(ret)
}

type RootEntity struct {
	Success    bool             `json:"success"`
	Data       []DataEntity     `json:"data"`
	Pagination PaginationEntity `json:"pagination"`
}

type PaginationEntity struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}
type DataEntity struct {
	ID         uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      string `json:"title"`
	Img        string `json:"img"`
	Address    string `json:"address"`
	Platform   string `json:"platform"`
	PlatformID uint64 `json:"platform_id"`
}

func (DataEntity) TableName() string {
	return "streamer_online"
}

func insertPlatformStreamer() {
	config.ReadConfigByViper()
	pkg.ConnectDB()

	var ret RootEntity = RootEntity{}
	resp, err := http.Get("https://live.suancaihu.eu.org/api/streamers/mihu?page=10&pageSize=20")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bytes, &ret)

	var platforms []Platform = make([]Platform, 0, 20)
	if err := pkg.DB.Find(&platforms).Error; err != nil {
		fmt.Println(err.Error())
	}

	for i := range ret.Data {
		for _, p := range platforms {
			if ret.Data[i].Platform == p.Name {
				ret.Data[i].PlatformID = p.ID
			}
		}
	}

	dbError := pkg.DB.Create(&ret.Data).Error
	if dbError != nil {
		fmt.Println(dbError)
	}

}

func migrate() {
	// config.ReadConfigByViper()
	// service.ConnectDB()
	// 生成表
	// DB.AutoMigrate(&model.UserDO{})
	// service.DB.AutoMigrate(&Platform{})
}
