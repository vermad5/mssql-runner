package mssqlRunner

import (
	. "runner"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/denisenkom/go-mssqldb"
)

type MssqlRunner struct{}

func (m MssqlRunner) Metadata() SpaceAgentRunnerMetadata {
	return SpaceAgentRunnerMetadata{
		Name:          "MSSQL Test",
		Description:   "MSSQL Database test. Enter credentials for your existing MSSQL database and see if an app deployed to your ORG and SPACE is able to access that database.",
		Endpoint:      "/mssql_test",
		Params:        parameters(),
		UpsCompatible: true,
	}
}
func parameters() []SpaceAgentRunnerParameter {
	var parameters []SpaceAgentRunnerParameter

	parameters = append(parameters, SpaceAgentRunnerParameter{Name: "host", Secret: false})
	parameters = append(parameters, SpaceAgentRunnerParameter{Name: "port", Secret: false})
	parameters = append(parameters, SpaceAgentRunnerParameter{Name: "username", Secret: false})
	parameters = append(parameters, SpaceAgentRunnerParameter{Name: "password", Secret: true})
	parameters = append(parameters, SpaceAgentRunnerParameter{Name: "databaseName", Secret: false})

	return parameters
}

func (m MssqlRunner) Handler(c *gin.Context) {
	host := c.PostForm("host")
	port := c.PostForm("port")
	username := c.PostForm("username")
	password := c.PostForm("password")
	databaseName := c.PostForm("databaseName")

	result := mssqlAudit(host, port, username, password, databaseName, "5s")
	c.JSON(200, gin.H{
		"canConnect":      result == "",
		"connectionError": result,
	})
}

func mssqlAudit(host, port, username, password, databaseName, timeout string) string {
	//"sqlserver://demouser1:root@10.30.233.100:1433?database=master&connection+timeout=30"
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?readTimeout=%s", username, password, host, port, databaseName, timeout)
	fmt.Println(connectionString)
	db, openError := sql.Open("mssql", connectionString)

	if openError != nil {
		return openError.Error()
	}

	pingError := db.Ping()
	db.Close()
	if pingError != nil {
		return pingError.Error()
	}

	return ""
}
