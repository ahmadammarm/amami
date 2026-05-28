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

type SMTPConfigResponse struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	From     string `json:"from"`
}

type TestSMTPRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type HealthResponse struct {
	Database bool `json:"database"`
	SMTP     bool `json:"smtp"`
}

type AuditLogResponse struct {
	ID        uint   `json:"id"`
	User      string `json:"user"`
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	Timestamp string `json:"timestamp"`
}

type AuditLogListResponse struct {
	Logs       []AuditLogResponse `json:"logs"`
	TotalCount int64              `json:"total_count"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}
