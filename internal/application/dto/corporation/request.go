package corporationdto

type RegisterRequest struct {
	Name		string
	CIN 		string
	Password	string
}

type LoginRequest struct {
	CIN		string
	Password	string
}