package mssqlRunner_test

import (
	"mssql-runner"
	. "runner"

	"net"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

)

var _ = Describe("TestRunner", func() {

	Describe("MSSQLRunner", func() {
		router := gin.New()
		mssqlRunner := mssqlRunner.MssqlRunner{}
		router.POST("/mssql_test", mssqlRunner.Handler)

		Describe("Metadata tests", func() {
			Context("MSSQL Runner Metadata", func() {
				It("Has the right parameters", func() {
					metadata := mssqlRunner.Metadata()
					Expect(ParameterNames(metadata.Params)).To(ContainElement("host"))
					Expect(ParameterNames(metadata.Params)).To(ContainElement("username"))
					Expect(ParameterNames(metadata.Params)).To(ContainElement("password"))
					Expect(ParameterNames(metadata.Params)).To(ContainElement("port"))
					Expect(ParameterNames(metadata.Params)).To(ContainElement("databaseName"))
					Expect(FindParameterWithName(metadata.Params, "host").Secret).To(BeFalse())
					Expect(FindParameterWithName(metadata.Params, "username").Secret).To(BeFalse())
					Expect(FindParameterWithName(metadata.Params, "password").Secret).To(BeTrue())
					Expect(FindParameterWithName(metadata.Params, "port").Secret).To(BeFalse())
					Expect(FindParameterWithName(metadata.Params, "databaseName").Secret).To(BeFalse())
					Expect(metadata.UpsCompatible).To(BeTrue())
				})
			})
		})

		Describe("Integration tests", func() {
			Context("when it can reach the host", func() {
				It("returns true", func() {
					var responseBody = DoPost(router, "/mssql_test", map[string]string{
						"host":         "localhost",
						"port":         "1433",
						"username":     "demouser1",
						"password":     "root",
						"databaseName": "master",
					})
					Expect(responseBody).To(Equal("{\"canConnect\":true,\"connectionError\":\"\"}"))
				})
			})

			Context("when it can not reach the host", func() {
				It("returns false", func() {
					var responseBody = DoPost(router, "/mssql_test", map[string]string{
						"host":         "unknownhost",
						"port":         "1433",
						"username":     "demouser1",
						"password":     "root",
						"databaseName": "master",
					})
					Expect(responseBody).To(ContainSubstring("\"canConnect\":false"))
					Expect(responseBody).NotTo(ContainSubstring("\"connectionError\":\"\""))
				})
			})

			Context("when it tries to connect with an incorrect username", func() {
				It("returns false", func() {
					var responseBody = DoPost(router, "/mssql_test", map[string]string{
						"host":         "localhost",
						"port":         "1433",
						"username":     "blah",
						"password":     "root",
						"databaseName": "master",
					})
					Expect(responseBody).To(ContainSubstring("\"canConnect\":false"))
					Expect(responseBody).To(ContainSubstring("Access denied for user"))
					Expect(responseBody).To(ContainSubstring("using password: YES"))
				})
			})

			Context("when it tries to connect with an incorrect password", func() {
				It("returns false", func() {
					var responseBody = DoPost(router, "/mssql_test", map[string]string{
						"host":         "localhost",
						"port":         "1433",
						"username":     "demouser1",
						"password":     "blah",
						"databaseName": "master",
					})
					Expect(responseBody).To(ContainSubstring("\"canConnect\":false"))
					Expect(responseBody).To(ContainSubstring("Access denied for user"))
					Expect(responseBody).To(ContainSubstring("using password: YES"))
				})
			})

			Context("when it tries to connect to a nonexistent database", func() {
				It("returns false", func() {
					var responseBody = DoPost(router, "/mssql_test", map[string]string{
						"host":         "localhost",
						"port":         "1433",
						"username":     "demouser1",
						"password":     "root",
						"databaseName": "blah",
					})
					Expect(responseBody).To(ContainSubstring("\"canConnect\":false"))
					Expect(responseBody).To(ContainSubstring("Unknown database"))
				})
			})

			Context("when it tries to connect to a non-MySQL process", func() {
				var ln net.Listener
				BeforeEach(func() {
					ln, _ = net.Listen("tcp", "127.0.0.1:3333")
				})
				AfterEach(func() {
					ln.Close()
				})

				It("returns false", func() {
					var responseBody = DoPost(router, "/mssql_test", map[string]string{
						"host":         "localhost",
						"port":         "3333",
						"username":     "demouser1",
						"password":     "root",
						"databaseName": "space_agent_db",
					})
					Expect(responseBody).To(ContainSubstring("\"canConnect\":false"))
					Expect(responseBody).To(ContainSubstring("driver: bad connection"))
				})
			})
		})
	})
})
