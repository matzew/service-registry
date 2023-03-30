package endpoint

type EndpointConfig struct {
	Endpoints []string          `json:"endpoints"`
	Protocol  string            `json:"protocol"`
	Strict    bool              `json:"strict"`
	Options   map[string]string `json:"options"`
}

type Endpoint struct {
	ID              string         `json:"id"`
	Usage           string         `json:"usage"`
	Config          EndpointConfig `json:"config"`
	Channel         string         `json:"channel"`
	DefinitionGroup []string       `json:"definitionGroups"`
	Description     string         `json:"description"`
	Name            string         `json:"name"`
	Version         int64          `json:"version"`
	Format          string         `json:"format"`
	AuthScope       string         `json:"authscope"`
}
