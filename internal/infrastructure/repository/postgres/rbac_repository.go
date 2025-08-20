package postgres

type RBACRepository struct{}

func NewRBACRepository() *RBACRepository {
	return &RBACRepository{}
}
