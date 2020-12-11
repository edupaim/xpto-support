package domain

import (
	"time"
)

type Negative struct {
	CompanyDocument  string    `json:"companyDocument"`
	CompanyName      string    `json:"companyName"`
	CustomerDocument string    `json:"customerDocument"`
	Value            float64   `json:"value"`
	Contract         string    `json:"contract"`
	DebtDate         time.Time `json:"debtDate"`
	InclusionDate    time.Time `json:"inclusionDate"`
}

func (n *Negative) DatesToUTC() {
	n.DebtDate = n.DebtDate.UTC()
	n.InclusionDate = n.InclusionDate.UTC()
}

func (n *Negative) EncryptCustomerDocument() error {

	return nil
}

func (n *Negative) DecryptCustomerDocument() error {

	return nil
}
