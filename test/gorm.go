package test

import "gorm.io/gorm"

var db *gorm.DB

func test() {
	//
	db.Create(&Test1{Username: "test"})
	var tests []Test1 = []Test1{}
	db.CreateInBatches(&tests, 10)
	db.Create(&tests)
	db.Save(&Test1{Username: "test2"})

}

type Test1 struct {
	ID       int    `gorm:"primary_key"`
	Username string `gorm:"unique"`
}
