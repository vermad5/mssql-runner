package mssqlRunner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMssqlRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mssqlRunner Suite")
}
