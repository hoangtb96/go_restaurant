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

func ListRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var filter restaurantmodel.Filter
		if err := ctx.ShouldBind(&filter); err != nil {
			ctx.JSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}
		var paging common.Paging
		if err := ctx.ShouldBind(&paging); err != nil {
			ctx.JSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}
		paging.Fulfill()
		store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantbiz.NewListRestaurantBiz(store)

		result, err := biz.ListRestaurant(ctx.Request.Context(), &filter, &paging)
		if err != nil {
			ctx.JSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, common.NewSucessResponse(result, paging, filter))
	}
}
