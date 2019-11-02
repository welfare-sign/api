package wsgin

// BasePagingResponse 统一返回对象
type BasePagingResponse struct {
	BaseResponse
	Total   int `json:"total,omitempty"`   //总数量
	Current int `json:"current,omitempty"` //当前页页码
}
