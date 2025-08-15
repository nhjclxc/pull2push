package request

type CreateUserRequest struct {
	ID       int    `json:"id" binding:"required"`
	LiveURL  string `json:"live_url" binding:"required"`
	LiveName string `json:"live_name" binding:"required"`
}
