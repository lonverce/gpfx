package domain

type Repository[T, TKey any] interface {
	FindByID(id TKey) (entity *T, exist bool)
	Insert(entity *T)
	Update(entity *T)
	Delete(entity *T)
}
