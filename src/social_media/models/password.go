package models

// PasswordChange holds data needed for changing passwords.
type PasswordChange struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
	SessionId       string `json:"sessionId"`
}
