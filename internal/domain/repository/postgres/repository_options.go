package postgres

type QueryModifier interface {
	Apply(query interface{}) interface{}
}
