package model

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Alias from MYSQL
type Alias struct {
	Model
	SourceUsername      string `json:"source_username" validate:"catchall" gorm:"unique:unq_alias"`
	SourceDomainID      int32  `json:"source_domain_id" gorm:"required,unique:unq_alias" validate:"required"`
	SourceDomain        Domain `json:"source_domain"`
	DestinationUsername string `json:"destination_username" gorm:"required,unique:unq_alias" validate:"required"`
	DestinationDomain   string `json:"destination_domain" gorm:"required,unique:unq_alias" validate:"required,fqdn"`
	Enabled             bool   `json:"enabled"`
}

// Domain from MYSQL
type Domain struct {
	Model
	Name         string `json:"name" gorm:"uniqueIndex" validate:"required,fqdn"`
	DkimSelector string `json:"dkim-selector"`
	Checks       Checks
}

type Checks struct {
	Model

	DomainID    uint
	MXRecord    bool
	SPFRecord   bool
	DMARCRecord bool
	DKIMRecord  bool
}

// Account from MYSQL
type Account struct {
	Model
	ID       int     `json:"id"`
	Username string  `json:"username" validate:"required"`
	DomainID uint    `json:"domain_id" validate:"required_if=update false"`
	Domain   *Domain `json:"domain"`
	Password string  `json:"password,omitempty" validate:"required_if=update false"`
	Quota    int32   `json:"quota"`
	Enabled  bool    `json:"enabled"`
	SendOnly bool    `json:"sendonly"`
}

type TLSPolicy struct {
	Model
	DomainID int `json:"domain_id" gorm:"unique"`
	Domain   *Domain
	Policy   string  `json:"policy" validate:"required,oneof=none may encrypt dane dane-only fingerprint verify secure"`
	Params   *string `json:"params"`
}
