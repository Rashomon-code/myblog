package model

type MyPageResponse struct {
	UserProfile UserProfile      `json:"userprofile"`
	Posts       []ArticleSummary `json:"posts"`
}

type UserProfile struct {
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}
