package common

type successRes struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}

func NewSucessResponse(data, paging, filter interface{}) *successRes {
	return &successRes{data, paging, filter}
}

func SimpleSucessResponse(data interface{}) *successRes {
	return NewSucessResponse(data, nil, nil)
}
