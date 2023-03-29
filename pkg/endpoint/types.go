package endpoint

type EndpointConfig struct {
	Endpoints []string          `json:"endpoints"`
	Protocol  string            `json:"protocol"`
	Strict    bool              `json:"strict"`
	Options   map[string]string `json:"options"`
}

type Endpoint struct {
	Format          string         `json:"format"`
	AuthScope       string         `json:"authscope"`
	Config          EndpointConfig `json:"config"`
	DefinitionGroup []string       `json:"definitionGroups"`
	Description     string         `json:"description"`
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Usage           string         `json:"usage"`
	Channel         string         `json:"channel"`
	Version         int64          `json:"version"`
}
