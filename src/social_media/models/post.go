package models

// Post holds all fields used in the posts table in the database.
type Post struct {
	Id         string `json:"id"`
	UnixTime   int64  `json:"unixTime"`
	Username   string `json:"username"`
	Message    string `json:"message"`
	ReplyingTo string `json:"replyingTo"`
	NumReplies int    `json:"numReplies"`
	Media      string `json:"media"`
}
