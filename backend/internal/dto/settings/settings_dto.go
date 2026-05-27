package settings

type MosqueProfileRequest struct {
	Name            string `json:"name" binding:"required"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	Logo            string `json:"logo"`
	LegalYayasanID  string `json:"legal_yayasan_id"`
}

type MosqueProfileResponse struct {
	Name            string `json:"name"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	Logo            string `json:"logo"`
	LegalYayasanID  string `json:"legal_yayasan_id"`
}

type SMTPConfigRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	From     string `json:"from" binding:"required"`
}
