package contract

type PagedDto[T any] struct {
	Total int
	Items []T
}

type PagingRequestDto struct {
	Skip int `form:"skip"`
	Take int `form:"take"`
}
