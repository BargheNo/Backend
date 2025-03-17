package corporationdto


type CorporationInfoResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Name		 string   `json:"name"`
}

type InstallationRequestResponse struct {
	ID            	uint 		`json:"id"`
	UserID        	string 		`json:"userId"`
	Area 		  	float64   	`json:"area"`
	PowerRequested 	float64   	`json:"powerRequested"`
	MaxCost	   		float64   	`json:"maxCost"`
	Deadline       	string 		`json:"deadline"`
	BuildingType	string 		`json:"buildingType"`
}