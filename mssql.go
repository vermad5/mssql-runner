package mssqlRunner

import (
	. "runner"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
	"strconv"
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
	port,err := strconv.Atoi(c.PostForm("port"))
	username := c.PostForm("username")
	password := c.PostForm("password")
	databaseName := c.PostForm("databaseName")

    if err == nil {
		result := mssqlAudit(host, username, password, databaseName, "5s", port)
		c.JSON(200, gin.H{
			"canConnect":      result == "",
			"connectionError": result,
		})
	}

	if err != nil{
		fmt.Println("Error caught")
	}
}

func mssqlAudit(host, username, password, databaseName, timeout string, port int) string {
	//"sqlserver://demouser1:root@10.30.233.100:1433?database=master&connection+timeout=30"
	query := url.Values{}
	query.Add("database", fmt.Sprintf("%s", databaseName))
	query.Add("connection timeout", fmt.Sprintf("%d", 15))

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%d", host, port),
		RawQuery: query.Encode(),
	}

	connectionString := u.String()
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
