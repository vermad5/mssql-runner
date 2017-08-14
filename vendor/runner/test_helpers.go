package runner

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	 . "github.com/onsi/gomega"
)

func FindParameterWithName(params []SpaceAgentRunnerParameter, name string) SpaceAgentRunnerParameter {
	var paramWithName SpaceAgentRunnerParameter
	for _, param := range params {
		if param.Name == name {
			paramWithName = param
		}
	}
	return paramWithName
}

func ParameterNames(parameters []SpaceAgentRunnerParameter) []string {
	var parameterNames []string
	for _, param := range parameters {
		parameterNames = append(parameterNames, param.Name)
	}
	return parameterNames
}

func DoPost(router *gin.Engine, endpoint string, params map[string]string) string {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", endpoint, encodeParams(params))
	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded; param=value",
	)
	router.ServeHTTP(response, request)

	return strings.TrimSpace(response.Body.String())
}

func encodeParams(params map[string]string) *strings.Reader {
	var values = url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	return strings.NewReader(values.Encode())
}

var AssertThat = Ω
