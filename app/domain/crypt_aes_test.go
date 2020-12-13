package domain

import (
	"github.com/onsi/gomega"
	"testing"
)

func Test_encrypt(t *testing.T) {
	g := gomega.NewWithT(t)
	t.Run("same data should should be same encrypted", func(t *testing.T) {
		encrypted, err := encrypt("ola")
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		encrypted2, err := encrypt("ola")
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(encrypted).Should(gomega.Equal(encrypted2))
	})
	t.Run("crypt and decrypt should be same data", func(t *testing.T) {
		encrypted, err := encrypt("ola")
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		decrypted, err := decrypt(encrypted)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(decrypted).Should(gomega.Equal("ola"))
	})
}
