package models

// Profile holds all fields used in the users table in the database, and a boolean for ownership of the profile page.
type Profile struct {
	Username string  `json:"username"`
	JoinDate float64 `json:"joinDate"`
	Bio      string  `json:"bio"`
	Media    string  `json:"media"`
	Owner    bool    `json:"owner"`
	Admin    bool    `json:"admin"`
}
