// https://busynose.com/posts/pessimistic_lock/
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Order struct {
	Id          int64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	GoodsID     int64 `gorm:"column:goods_id;default:0;NOT NULL"` // 商品ID
	GoodsNumber int64 `gorm:"column:goods_number;NOT NULL"`       // 商品库存
	Version     int64 `gorm:"column:version;default:0"`           // 版本
}

func init() {
	var err error
	dsn := "root:123@/test?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = db.Migrator().DropTable(&Order{})
	_ = db.AutoMigrate(&Order{})
	// 初始化100件库存
	db.Create(&Order{GoodsID: 1, GoodsNumber: 100})
}

func main() {
	if err := StartServer(); err != nil {
		panic(err)
	}
	select {}
}

func StartServer() error {
	r := gin.Default()
	r.POST("/order", order)
	return r.Run()
}

func order(c *gin.Context) {
	type request struct {
		ID int64 `json:"id"`
	}
	var body request
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var row Order
	db.First(&row, Order{Id: 1})
	if row.GoodsNumber <= 0 {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	row.GoodsNumber = row.GoodsNumber - 1
	column := db.Model(&row).Where("id", row.Id).Where("version", row.Version).
		UpdateColumn("goods_number", row.GoodsNumber).
		UpdateColumn("version", gorm.Expr("version+1"))
	if column.Error != nil || column.RowsAffected == 0 {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, nil)
}
