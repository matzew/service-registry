package registry

type Group struct {
	ID          string                `json:"id"`
	Version     int                   `json:"version"`
	Format      string                `json:"format"`
	Definitions map[string]Definition `json:"definitions"`
}

type Definition struct {
	ID          string   `json:"id"`
	Version     int      `json:"version"`
	Description string   `json:"description"`
	SchemaURL   string   `json:"schemaUrl"`
	Format      string   `json:"format"`
	Metadata    Metadata `json:"metadata"`
}

type Metadata struct {
	Attributes map[string]Attribute `json:"attributes"`
}

type Attribute struct {
	Required bool        `json:"required"`
	Value    interface{} `json:"value,omitempty"`
}
