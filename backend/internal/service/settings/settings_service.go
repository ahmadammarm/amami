package settings

import (
	"github.com/ahmadammarm/amami/backend/internal/dto/settings"
	repo "github.com/ahmadammarm/amami/backend/internal/repository/settings"
)

type SettingsService interface {
	UpdateMosqueProfile(req settings.MosqueProfileRequest) error
	GetMosqueProfile() (*settings.MosqueProfileResponse, error)
	UpdateSMTPConfig(req settings.SMTPConfigRequest) error
}

type settingsService struct {
	repo repo.SettingsRepository
}

func NewSettingsService(repo repo.SettingsRepository) SettingsService {
	return &settingsService{repo: repo}
}

func (s *settingsService) UpdateMosqueProfile(req settings.MosqueProfileRequest) error {
	settingsData := map[string]string{
		"mosque_name":       req.Name,
		"mosque_address":    req.Address,
		"mosque_phone":      req.Phone,
		"mosque_logo":       req.Logo,
		"legal_yayasan_id":  req.LegalYayasanID,
	}

	for k, v := range settingsData {
		if err := s.repo.UpsertSetting(k, v, false); err != nil {
			return err
		}
	}
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

	return res, nil
}

func (s *settingsService) UpdateSMTPConfig(req settings.SMTPConfigRequest) error {
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
	return nil
}
