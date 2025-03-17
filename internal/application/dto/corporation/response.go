package corporationdto


type CorporationInfoResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Name		 string   `json:"name"`
}

type InstallationRequestResponse struct {
	ID            	uint 		`json:"id"`
	UserID        	uint 		`json:"userId"`
	Area 		  	float64   	`json:"area"`
	PowerRequested 	float64   	`json:"powerRequested"`
	MaxCost	   		float64   	`json:"maxCost"`
	Deadline       	string 		`json:"deadline"`
	BuildingType	string 		`json:"buildingType"`
	Address			string 		`json:"address"`
}

type BidsResponse struct {
	ID            			uint 		`json:"id"`
	InstallationRequestID 	uint 	`json:"installationRequestId"`
	Description 			string 		`json:"description"`
	MinCost      			float64   	`json:"minCost"`
	MaxCost      			float64   	`json:"maxCost"`
	MinDeadline  			string 		`json:"minDeadline"`
	MaxDeadline  			string 		`json:"maxDeadline"`
	InstallationTime 		string 	`json:"installationTime"`
	Status	   				string 		`json:"status"`
}