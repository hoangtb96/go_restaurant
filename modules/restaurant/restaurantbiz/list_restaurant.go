package restaurantbiz

import (
	"context"
	"example/common"
	"example/modules/restaurant/restaurantmodel"
)

type ListRestaurantStore interface {
	ListDataByCondition(ctx context.Context,
		condition map[string]interface{},
		filter *restaurantmodel.Filter,
		paging *common.Paging,
		moreKeys ...string) ([]restaurantmodel.Restaurant, error)
}

type listRestaurantBiz struct {
	store ListRestaurantStore
}

func NewListRestaurantBiz(store ListRestaurantStore) *listRestaurantBiz {
	return &listRestaurantBiz{store: store}
}

func (biz *listRestaurantBiz) ListRestaurant(ctx context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging) ([]restaurantmodel.Restaurant, error) {

	result, err := biz.store.ListDataByCondition(ctx, nil, filter, paging)
	return result, err
}
