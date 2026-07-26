package settings

import (
	"fmt"
	"net/smtp"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/settings"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/settings"
	"github.com/google/uuid"
)

type SettingsService interface {
	UpdateMosqueProfile(userID uuid.UUID, req settings.MosqueProfileRequest) error
	GetMosqueProfile() (*settings.MosqueProfileResponse, error)
	UpdateSMTPConfig(userID uuid.UUID, req settings.SMTPConfigRequest) error
	GetSMTPConfig() (*settings.SMTPConfigResponse, error)
	TestSMTPConnection(req settings.TestSMTPRequest) error
	GetSystemHealth() (*settings.HealthResponse, error)
	GetAuditLogs(page, limit int, userID uuid.UUID, roleID uint) (*settings.AuditLogListResponse, error)
}

type settingsService struct {
	repo repo.SettingsRepository
}

func (s *settingsService) GetSystemHealth() (*settings.HealthResponse, error) {
	dbOk := s.repo.Ping() == nil
	
	// Check SMTP health using current config
	smtpConf, _ := s.GetSMTPConfig()
	smtpOk := false
	if smtpConf.Host != "" {
		// Just a shallow check, dial the host
		addr := fmt.Sprintf("%s:%s", smtpConf.Host, smtpConf.Port)
		conn, err := smtp.Dial(addr)
		if err == nil {
			smtpOk = true
			conn.Close()
		}
	}

	return &settings.HealthResponse{
		Database: dbOk,
		SMTP:     smtpOk,
	}, nil
}

func (s *settingsService) GetAuditLogs(page, limit int, userID uuid.UUID, roleID uint) (*settings.AuditLogListResponse, error) {
	logs, total, err := s.repo.GetAuditLogs(page, limit, userID, roleID)
	if err != nil {
		return nil, err
	}

	var logResponses []settings.AuditLogResponse
	for _, l := range logs {
		logResponses = append(logResponses, settings.AuditLogResponse{
			ID:        l.ID,
			User:      l.User.Username,
			Action:    l.Action,
			Entity:    l.Entity,
			Timestamp: l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	
	if logResponses == nil {
		logResponses = []settings.AuditLogResponse{}
	}

	return &settings.AuditLogListResponse{
		Logs:       logResponses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}


func (s *settingsService) GetSMTPConfig() (*settings.SMTPConfigResponse, error) {
	res := &settings.SMTPConfigResponse{}

	val, err := s.repo.GetSetting("smtp_host")
	if err == nil { res.Host = val.Value }

	val, err = s.repo.GetSetting("smtp_port")
	if err == nil { res.Port = val.Value }

	val, err = s.repo.GetSetting("smtp_user")
	if err == nil { res.Username = val.Value }

	val, err = s.repo.GetSetting("smtp_from")
	if err == nil { res.From = val.Value }

	return res, nil
}

func (s *settingsService) TestSMTPConnection(req settings.TestSMTPRequest) error {
	auth := smtp.PlainAuth("", req.Username, req.Password, req.Host)
	addr := fmt.Sprintf("%s:%s", req.Host, req.Port)

	// Attempt to connect to the SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Try to authenticate if username is provided
	if req.Username != "" {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	return nil
}


func NewSettingsService(repo repo.SettingsRepository) SettingsService {
	return &settingsService{repo: repo}
}

func (s *settingsService) UpdateMosqueProfile(userID uuid.UUID, req settings.MosqueProfileRequest) error {
	settingsData := map[string]string{
		"mosque_name":         req.Name,
		"mosque_address":      req.Address,
		"mosque_phone":        req.Phone,
		"mosque_logo":         req.Logo,
		"legal_yayasan_id":    req.LegalYayasanID,
		"bank_name":           req.BankName,
		"bank_account_name":   req.BankAccountName,
		"bank_account_number": req.BankAccountNumber,
		"zakat_fitrah_amount": req.ZakatFitrahAmount,
		"active_hijri_year":   req.ActiveHijriYear,
	}

	for k, v := range settingsData {
		if err := s.repo.UpsertSetting(k, v, false); err != nil {
			return err
		}
	}

	// Record Audit Log
	s.repo.CreateAuditLog(domain.AuditLog{
		UserID: userID,
		Action: "UPDATE",
		Entity: "Mosque Profile",
	})

	return nil
}

func (s *settingsService) GetMosqueProfile() (*settings.MosqueProfileResponse, error) {
	res := &settings.MosqueProfileResponse{}
	
	val, err := s.repo.GetSetting("mosque_name")
	if err == nil { res.Name = val.Value }
	
	val, err = s.repo.GetSetting("mosque_address")
	if err == nil { res.Address = val.Value }
	
	val, err = s.repo.GetSetting("mosque_phone")
	if err == nil { res.Phone = val.Value }
	
	val, err = s.repo.GetSetting("mosque_logo")
	if err == nil { res.Logo = val.Value }
	
	val, err = s.repo.GetSetting("legal_yayasan_id")
	if err == nil { res.LegalYayasanID = val.Value }

	val, err = s.repo.GetSetting("bank_name")
	if err == nil { res.BankName = val.Value }

	val, err = s.repo.GetSetting("bank_account_name")
	if err == nil { res.BankAccountName = val.Value }

	val, err = s.repo.GetSetting("bank_account_number")
	if err == nil { res.BankAccountNumber = val.Value }

	val, err = s.repo.GetSetting("zakat_fitrah_amount")
	if err == nil { res.ZakatFitrahAmount = val.Value }

	val, err = s.repo.GetSetting("active_hijri_year")
	if err == nil { res.ActiveHijriYear = val.Value }

	return res, nil
}

func (s *settingsService) UpdateSMTPConfig(userID uuid.UUID, req settings.SMTPConfigRequest) error {
	settingsData := map[string]struct {
		val    string
		secret bool
	}{
		"smtp_host": {req.Host, false},
		"smtp_port": {req.Port, false},
		"smtp_user": {req.Username, false},
		"smtp_pass": {req.Password, true},
		"smtp_from": {req.From, false},
	}

	for k, v := range settingsData {
		if err := s.repo.UpsertSetting(k, v.val, v.secret); err != nil {
			return err
		}
	}

	// Record Audit Log
	s.repo.CreateAuditLog(domain.AuditLog{
		UserID: userID,
		Action: "UPDATE",
		Entity: "SMTP Configuration",
	})

	return nil
}
