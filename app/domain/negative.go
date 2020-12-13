package domain

import (
	"errors"
	"time"
)

var (
	DecodeNegativeJsonError = errors.New("decode negative json")
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
	document, err := n.CustomerDocument.Encrypt()
	if err != nil {
		return err
	}
	n.CustomerDocument = document
	return nil
}

func (n *Negative) DecryptCustomerDocument() error {
	document, err := n.CustomerDocument.Decrypt()
	if err != nil {
		return err
	}
	n.CustomerDocument = document
	return nil
}

type CustomerDocument string

func (c CustomerDocument) Encrypt() (CustomerDocument, error) {
	crypt, err := encrypt(string(c))
	if err != nil {
		return "", err
	}
	return CustomerDocument(crypt), nil
}

func (c CustomerDocument) Decrypt() (CustomerDocument, error) {
	crypt, err := decrypt(string(c))
	if err != nil {
		return "", err
	}
	return CustomerDocument(crypt), nil
}
