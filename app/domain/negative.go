package domain

import (
	"errors"
	"time"
)

var (
	DecodeNegativeJsonError      = errors.New("decode negative json")
	EncryptCustomerDocumentError = errors.New("encrypt customer document")
	DecryptCustomerDocumentError = errors.New("decrypt customer document")
)

type Negative struct {
	CompanyDocument  string           `json:"companyDocument"`
	CompanyName      string           `json:"companyName"`
	CustomerDocument CustomerDocument `json:"customerDocument"`
	Value            float64          `json:"value"`
	Contract         string           `json:"contract"`
	DebtDate         time.Time        `json:"debtDate"`
	InclusionDate    time.Time        `json:"inclusionDate"`
}

func (n *Negative) DatesToUTC() {
	n.DebtDate = n.DebtDate.UTC()
	n.InclusionDate = n.InclusionDate.UTC()
}

func (n *Negative) EncryptCustomerDocument() error {
	return n.CustomerDocument.Encrypt()
}

func (n *Negative) DecryptCustomerDocument() error {
	return n.CustomerDocument.Decrypt()
}

type CustomerDocument string

func (c *CustomerDocument) Encrypt() error {
	crypt, err := encrypt(string(*c))
	if err != nil {
		return err
	}
	*c = CustomerDocument(crypt)
	return nil
}

func (c *CustomerDocument) Decrypt() error {
	crypt, err := decrypt(string(*c))
	if err != nil {
		return err
	}
	*c = CustomerDocument(crypt)
	return nil
}
