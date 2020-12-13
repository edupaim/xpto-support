package domain

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"testing"
	"time"
)

func TestNegative_DatesToUTC(t *testing.T) {
	g := gomega.NewWithT(t)
	location := time.FixedZone("GMT -3", -3*3600)
	negative := Negative{
		CompanyDocument:  "59291534000167",
		CompanyName:      "ABC S.A.",
		CustomerDocument: "51537476467",
		Value:            1235.23,
		Contract:         "bc063153-fb9e-4334-9a6c-0d069a42065b",
		DebtDate:         time.Date(2015, 11, 13, 20, 32, 51, 0, location),
		InclusionDate:    time.Date(2020, 11, 13, 20, 32, 51, 0, location),
	}
	negative.DatesToUTC()
	g.Expect(negative).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"DebtDate":      gomega.BeIdenticalTo(time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC)),
		"InclusionDate": gomega.BeIdenticalTo(time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC)),
	}))
}

func TestNegative_EncryptCustomerDocument(t *testing.T) {
	g := gomega.NewWithT(t)
	customerDocument := CryptDocument("51537476467")
	companyDocument := CryptDocument("59291534000167")
	negative := Negative{
		CompanyDocument:  companyDocument,
		CompanyName:      "ABC S.A.",
		CustomerDocument: customerDocument,
		Value:            1235.23,
		Contract:         "bc063153-fb9e-4334-9a6c-0d069a42065b",
		DebtDate:         time.Date(2015, 11, 13, 20, 32, 51, 0, time.UTC),
		InclusionDate:    time.Date(2020, 11, 13, 20, 32, 51, 0, time.UTC),
	}
	err := negative.EncryptDocuments()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(negative.CustomerDocument).ShouldNot(gomega.Equal(customerDocument))
	g.Expect(negative.CompanyDocument).ShouldNot(gomega.Equal(companyDocument))
	err = negative.DecryptDocuments()
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(negative.CustomerDocument).Should(gomega.Equal(customerDocument))
	g.Expect(negative.CompanyDocument).Should(gomega.Equal(companyDocument))
}
