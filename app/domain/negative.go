package domain

import (
	"errors"
	"time"
)

var (
	DecodeNegativeJsonError      = errors.New("decode negative json")
	EncryptCustomerDocumentError = errors.New("encrypt customer document")
	DecryptCustomerDocumentError = errors.New("decrypt customer document")
	EncryptCompanyDocumentError  = errors.New("encrypt company document")
	DecryptCompanyDocumentError  = errors.New("decrypt company document")
)

type Negative struct {
	CompanyDocument  CryptDocument `json:"companyDocument"`
	CompanyName      string        `json:"companyName"`
	CustomerDocument CryptDocument `json:"customerDocument"`
	Value            float64       `json:"value"`
	Contract         string        `json:"contract"`
	DebtDate         time.Time     `json:"debtDate"`
	InclusionDate    time.Time     `json:"inclusionDate"`
}

func (n *Negative) DatesToUTC() {
	n.DebtDate = n.DebtDate.UTC()
	n.InclusionDate = n.InclusionDate.UTC()
}

func (n *Negative) EncryptDocuments() error {
	err := n.CompanyDocument.Encrypt()
	if err != nil {
		return err
	}
	return n.CustomerDocument.Encrypt()
}

func (n *Negative) DecryptDocuments() error {
	err := n.CompanyDocument.Decrypt()
	if err != nil {
		return err
	}
	return n.CustomerDocument.Decrypt()
}

func (n *Negative) EncryptCustomerDocument() error {
	err := n.CustomerDocument.Encrypt()
	if err != nil {
		return EncryptCustomerDocumentError
	}
	return nil
}

func (n *Negative) DecryptCustomerDocument() error {
	err := n.CustomerDocument.Decrypt()
	if err != nil {
		return DecryptCustomerDocumentError
	}
	return nil
}

func (n *Negative) EncryptCompanyDocument() error {
	err := n.CompanyDocument.Encrypt()
	if err != nil {
		return EncryptCompanyDocumentError
	}
	return nil
}

func (n *Negative) DecryptCompanyDocument() error {
	err := n.CompanyDocument.Decrypt()
	if err != nil {
		return DecryptCompanyDocumentError
	}
	return nil
}
