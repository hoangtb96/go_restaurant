package ginrestaurant

import (
	"example/common"
	"example/component"
	"example/modules/restaurant/restaurantbiz"
	"example/modules/restaurant/restaurantmodel"
	"example/modules/restaurant/restaurantstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data restaurantmodel.RestaurantCreate
		if err := ctx.ShouldBind(&data); err != nil {
			ctx.JSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}
		store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantbiz.NewCreateRestaurantBiz(store)

		if err := biz.CreateRestaurant(ctx.Request.Context(), &data); err != nil {
			ctx.JSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, common.SimpleSucessResponse(data))
	}
}
