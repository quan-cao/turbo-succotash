package entity

type UserProfile struct {
	Isid       string `json:"isid"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}
