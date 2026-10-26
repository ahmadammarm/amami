package settings

import (
	"fmt"
	"net/smtp"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/dto/settings"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/settings"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func NewSettingsService(repo repo.SettingsRepository) SettingsService {
	return &settingsService{repo: repo}
}

func (s *settingsService) UpdateMosqueProfile(userID uuid.UUID, req settings.MosqueProfileRequest) error {
	err := s.repo.UpsertSettings([]domain.SystemSetting{
		{Key: "mosque_name", Value: req.Name},
		{Key: "mosque_address", Value: req.Address},
		{Key: "mosque_phone", Value: req.Phone},
		{Key: "mosque_logo", Value: req.Logo},
		{Key: "legal_yayasan_id", Value: req.LegalYayasanID},
		{Key: "bank_name", Value: req.BankName},
		{Key: "bank_account_name", Value: req.BankAccountName},
		{Key: "bank_account_number", Value: req.BankAccountNumber},
		{Key: "zakat_fitrah_amount", Value: req.ZakatFitrahAmount},
		{Key: "active_hijri_year", Value: req.ActiveHijriYear},
	})
	if err != nil {
		return apperror.Internal("failed to update mosque profile", err)
	}

	s.recordAudit(userID, "Mosque Profile")
	return nil
}

func (s *settingsService) GetMosqueProfile() (*settings.MosqueProfileResponse, error) {
	values, err := s.loadSettings()
	if err != nil {
		return nil, err
	}

	return &settings.MosqueProfileResponse{
		Name:              values["mosque_name"],
		Address:           values["mosque_address"],
		Phone:             values["mosque_phone"],
		Logo:              values["mosque_logo"],
		LegalYayasanID:    values["legal_yayasan_id"],
		BankName:          values["bank_name"],
		BankAccountName:   values["bank_account_name"],
		BankAccountNumber: values["bank_account_number"],
		ZakatFitrahAmount: values["zakat_fitrah_amount"],
		ActiveHijriYear:   values["active_hijri_year"],
	}, nil
}

func (s *settingsService) UpdateSMTPConfig(userID uuid.UUID, req settings.SMTPConfigRequest) error {
	err := s.repo.UpsertSettings([]domain.SystemSetting{
		{Key: "smtp_host", Value: req.Host},
		{Key: "smtp_port", Value: req.Port},
		{Key: "smtp_user", Value: req.Username},
		{Key: "smtp_pass", Value: req.Password, IsSecret: true},
		{Key: "smtp_from", Value: req.From},
	})
	if err != nil {
		return apperror.Internal("failed to update SMTP configuration", err)
	}

	s.recordAudit(userID, "SMTP Configuration")
	return nil
}

// GetSMTPConfig never returns the password.
func (s *settingsService) GetSMTPConfig() (*settings.SMTPConfigResponse, error) {
	values, err := s.loadSettings()
	if err != nil {
		return nil, err
	}

	return &settings.SMTPConfigResponse{
		Host:     values["smtp_host"],
		Port:     values["smtp_port"],
		Username: values["smtp_user"],
		From:     values["smtp_from"],
	}, nil
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

func (s *settingsService) GetSystemHealth() (*settings.HealthResponse, error) {
	dbOk := s.repo.Ping() == nil

	// Shallow SMTP check: dial the configured host.
	smtpOk := false
	if smtpConf, err := s.GetSMTPConfig(); err == nil && smtpConf.Host != "" {
		addr := fmt.Sprintf("%s:%s", smtpConf.Host, smtpConf.Port)
		if conn, err := smtp.Dial(addr); err == nil {
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

	logResponses := make([]settings.AuditLogResponse, 0, len(logs))
	for _, l := range logs {
		logResponses = append(logResponses, settings.AuditLogResponse{
			ID:        l.ID,
			User:      l.User.Username,
			Action:    l.Action,
			Entity:    l.Entity,
			Timestamp: l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &settings.AuditLogListResponse{
		Logs:       logResponses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

// loadSettings reads every setting in one query, keyed by name. A key that was
// never saved reads as "".
func (s *settingsService) loadSettings() (map[string]string, error) {
	rows, err := s.repo.GetAllSettings()
	if err != nil {
		return nil, apperror.Internal("failed to fetch settings", err)
	}

	values := make(map[string]string, len(rows))
	for _, row := range rows {
		values[row.Key] = row.Value
	}
	return values, nil
}

// recordAudit is best effort: the settings are already saved, so a failed
// audit write is logged rather than reported to the caller.
func (s *settingsService) recordAudit(userID uuid.UUID, entity string) {
	if err := s.repo.CreateAuditLog(domain.AuditLog{UserID: userID, Action: "UPDATE", Entity: entity}); err != nil {
		zap.L().Warn("failed to write audit log", zap.String("entity", entity), zap.Error(err))
	}
}
