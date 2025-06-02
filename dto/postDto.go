package dto

type CreatePostRequest struct {
	Title     string `json:"title" binding:"required,min=3,max=255"`
	Subtitle  string `json:"subtitle" binding:"max=255"`
	Image     string `json:"image" binding:"max=255"`
	Content   string `json:"content" binding:"required"`
	Status    string `json:"status" binding:"omitempty,oneof=draft published archived"`
	CreatedBy uint   `json:"createdBy" binding:"required"`
}

type UpdatePostRequest struct {
	Title    string `json:"title" binding:"omitempty,min=3,max=255"`
	Subtitle string `json:"subtitle" binding:"omitempty,max=255"`
	Image    string `json:"image" binding:"omitempty,max=255"`
	Content  string `json:"content" binding:"omitempty"`
	Status   string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type PostResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Image     string `json:"image"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	CreatedBy uint   `json:"createdBy"`
}
