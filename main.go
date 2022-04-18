package main

import (
	"log"
	"net/http"
	"strconv"

	"example/component"
	"example/modules/restaurant/restauranttransport/ginrestaurant"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CREATE TABLE `notes` (
// 	`id` int(11) NOT NULL AUTO_INCREMENT,
// 	`title` varchar(100) NOT NULL,
// 	`content` text,
// 	`image` json DEFAULT NULL,
// 	`has_finished` tinyint(1) DEFAULT '0',
// 	`status` int(11) DEFAULT '1',
// 	`create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
// 	`update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 	PRIMARY KEY (`id`)
//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

type Note struct {
	Id      int    `json:"id,omitempty" gorm:"column:id;"`
	Title   string `json:"title" gorm:"column:title;"`
	Content string `json:"content" gorm:"column:content;"`
}

func (Note) TableName() string {
	return "notes"
}

// CREATE TABLE `restaurants` (
// 	`id` int(11) NOT NULL AUTO_INCREMENT,
// 	`owner_id` int(11) DEFAULT NULL,
// 	`name` varchar(50) NOT NULL,
// 	`addr` varchar(255) NOT NULL,
// 	`city_id` int(11) DEFAULT NULL,
// 	`lat` double DEFAULT NULL,
// 	`lng` double DEFAULT NULL,
// 	`logo` json DEFAULT NULL,
// 	`cover` json DEFAULT NULL,
// 	`shipping_fee_per_km` double DEFAULT '0',
// 	`status` int(11) NOT NULL DEFAULT '1',
// 	`created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
// 	`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 	PRIMARY KEY (`id`)
//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	Addr string `json:"address" gorm:"column:addr"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name"`
	Addr *string `json:"address" gorm:"column:addr"`
}

func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "food_delivery:19e5a718a54a9fe0559dfbce6908@tcp(127.0.0.1:3306)/food_delivery?charset=utf8mb4&parseTime=True&loc=Local"
	// dsn := os.Getenv("DBConnectionStr")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	if err := runService(db); err != nil {
		log.Fatalln(err)
	}

	// insert new note
	// newNote := Note{Title: "Demo note", Content: "This is content for demo note"}

	// if err := db.Create(&newNote); err != nil {
	// 	fmt.Println(err)
	// }

	// query
	// var notes []Note
	// db.Where("status = ?", 1).Find(&notes)

	// fmt.Println(notes)

}

func runService(db *gorm.DB) error {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	appCtx := component.NewAppContext(db)
	// CRUD
	restaurants := r.Group("/restaurants")
	{
		restaurants.POST("", ginrestaurant.CreateRestaurant(appCtx))

		restaurants.GET("/:id", func(ctx *gin.Context) {
			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})
				return
			}

			var data Restaurant

			if err := db.Where("id = ?", id).First(&data).Error; err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})

				return
			}
			ctx.JSON(http.StatusOK, data)
		})

		restaurants.GET("", ginrestaurant.ListRestaurant(appCtx))

		restaurants.PATCH("/:id", func(ctx *gin.Context) {

			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})
				return
			}
			var data RestaurantUpdate
			if err := ctx.ShouldBind(&data); err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"ok": 1})
		})

		restaurants.DELETE("/:id", func(ctx *gin.Context) {

			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})
				return
			}
			var data Restaurant
			if err := ctx.ShouldBind(&data); err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := db.Table(Restaurant{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
				ctx.JSON(401, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, data)
		})

	}

	return r.Run()
}
