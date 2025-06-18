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
	SourceUsername      string `json:"source_username" validate:"catchall" gorm:"uniqueIndex:unq_alias"`
	SourceDomainID      uint   `json:"source_domain_id" gorm:"index:unq_alias,unique" validate:"required"`
	SourceDomain        Domain `json:"source_domain"`
	DestinationUsername string `json:"destination_username" gorm:"index:unq_alias,unique" validate:"required"`
	DestinationDomain   string `json:"destination_domain" gorm:"index:unq_alias,unique" validate:"required,fqdn"`
	Enabled             bool   `json:"enabled"`
}

// Domain from MYSQL
type Domain struct {
	Model
	Name string `json:"name" gorm:"uniqueIndex" validate:"required,fqdn"`
}

// Account from MYSQL
type Account struct {
	Model
	Username string  `json:"username" validate:"required" gorm:"index:unq_account,unique"`
	DomainID uint    `json:"domain_id" gorm:"index:unq_account,unique"`
	Domain   *Domain `json:"domain"`
	Password string  `json:"password,omitempty"`
	Quota    int32   `json:"quota"`
	Enabled  bool    `json:"enabled"`
	SendOnly bool    `json:"sendonly"`
}

type TLSPolicy struct {
	Model
	DomainID uint `json:"domain_id" gorm:"unique"`
	Domain   *Domain
	Policy   string  `json:"policy" validate:"required,oneof=none may encrypt dane dane-only fingerprint verify secure"`
	Params   *string `json:"params"`
}
