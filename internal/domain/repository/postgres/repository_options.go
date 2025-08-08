package postgres

type QueryModifier2 interface {
	Apply(query interface{}) interface{}
}
