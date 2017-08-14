package runner

import "github.com/gin-gonic/gin"

type SpaceAgentRunnerParameter struct {
	Name          string   `json:"name"`
	Secret        bool     `json:"secret"`
}

type SpaceAgentRunnerMetadata struct {
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Endpoint      string                      `json:"endpoint"`
	UpsCompatible bool                        `json:"ups_compatible"`
	Params        []SpaceAgentRunnerParameter `json:"params"`
}

type SpaceAgentRunner interface {
	Metadata() SpaceAgentRunnerMetadata
	Handler(c *gin.Context)
}
