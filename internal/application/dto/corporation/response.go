package corporationdto

type CorporationInfoResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Name		 string   `json:"name"`
}