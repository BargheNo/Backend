package repository

type QueryModifier interface {
	Apply(query interface{}) interface{}
}
