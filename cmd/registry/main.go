package main

import (
	"encoding/json"
	"net/http"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Group struct {
	ID          string       `json:"id"`
	Version     int          `json:"version"`
	Definitions []Definition `json:"definitions"`
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

var groups = []Group{
	{
		ID:      "apiserversources.sources.knative.dev",
		Version: 1,
		Definitions: []Definition{
			{
				ID:          "dev.knative.apiserver.resource.add",
				Version:     133179207676972370,
				Description: "ApiServerSource CloudEvent type for adds",
				Format:      "CloudEvents/1.0",
				Metadata: Metadata{
					Attributes: map[string]Attribute{
						"datacontenttype": {
							Required: true,
							Value:    "application/json",
						},
						"my-add-extension": {
							Required: false,
						},
						"dataschema": {
							Required: true,
							Value:    "https://k8s.io/some/schema.json",
						},
						"id": {
							Required: true,
						},
						"source": {
							Required: false,
							Value:    "https://10.96.0.1:443",
						},
						"time": {
							Required: true,
						},
						"type": {
							Required: true,
							Value:    "dev.knative.apiserver.resource.add",
						},
					},
				},
			},
			{
				ID:          "dev.knative.apiserver.resource.update",
				Version:     133179207676972370,
				Description: "ApiServerSource CloudEvent type for updates",
				Format:      "CloudEvents/1.0",
				Metadata: Metadata{
					Attributes: map[string]Attribute{
						"datacontenttype": {
							Required: true,
							Value:    "application/json",
						},
						"my-update-extension": {
							Required: false,
						},
						"dataschema": {
							Required: true,
							Value:    "https://k8s.io/some/schema.json",
						},
						"id": {
							Required: true,
						},
						"source": {
							Required: false,
							Value:    "https://10.96.0.1:443",
						},
						"time": {
							Required: true,
						},
						"type": {
							Required: true,
							Value:    "dev.knative.apiserver.resource.update",
						},
					},
				},
			},
			{
				ID:          "dev.knative.apiserver.resource.delete",
				Version:     133179207676972370,
				Description: "ApiServerSource CloudEvent type for deletions",
				Format:      "CloudEvents/1.0",
				Metadata: Metadata{
					Attributes: map[string]Attribute{
						"datacontenttype": {
							Required: true,
							Value:    "application/json",
						},
						"my-delete-extension": {
							Required: false,
						},
						"dataschema": {
							Required: true,
							Value:    "https://k8s.io/some/schema.json",
						},
						"id": {
							Required: true,
						},
						"source": {
							Required: false,
							Value:    "https://10.96.0.1:443",
						},
						"time": {
							Required: true,
						},
						"type": {
							Required: true,
							Value:    "dev.knative.apiserver.resource.delete",
						},
					},
				},
			},
		},
	},
	{
		ID:      "pingsources.sources.knative.dev",
		Version: 1,
		Definitions: []Definition{
			{
				ID:          "dev.knative.sources.ping",
				Version:     133179207676972370,
				Description: "The default PingSource CloudEvent type",
				Format:      "CloudEvents/1.0",
				Metadata: Metadata{
					Attributes: map[string]Attribute{
						"datacontenttype": {
							Required: true,
							Value:    "application/json",
						},
						"dataschema": {
							Required: false,
							Value:    "https://k8s.io/some/schema.json",
						},
						"id": {
							Required: true,
						},
						"source": {
							Required: false,
							Value:    "/apis/v1/namespaces/default/pingsources/test-ping-source",
						},
						"time": {
							Required: true,
						},
						"type": {
							Required: true,
							Value:    "dev.knative.sources.ping",
						},
					},
				},
			},
		},
	},
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		renderGroups(w)
	case http.MethodPost:
		addGroupHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func renderGroups(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func groupByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groupID := strings.TrimPrefix(r.URL.Path, "/groups/")
	var foundGroup *Group
	for _, group := range groups {
		if group.ID == groupID {
			foundGroup = &group
			break
		}
	}

	if foundGroup == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundGroup)
}

func groupDefinitionsHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/definitions") {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groupID := strings.TrimPrefix(r.URL.Path, "/groups/")
	groupID = strings.TrimSuffix(groupID, "/definitions")
	var foundGroup *Group
	for _, group := range groups {
		if group.ID == groupID {
			foundGroup = &group
			break
		}
	}

	if foundGroup == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundGroup.Definitions)
}

func groupDefinitionByIDHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/groups/"), "/")
	if len(pathParts) != 3 || pathParts[1] != "definitions" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groupID := pathParts[0]
	defID := pathParts[2]
	var foundDefinition *Definition

	for _, group := range groups {
		if group.ID == groupID {
			for _, definition := range group.Definitions {
				if definition.ID == defID {
					foundDefinition = &definition
					break
				}
			}
		}
	}

	if foundDefinition == nil {
		http.Error(w, "Definition not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundDefinition)
}

func addGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type, expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Read and decode the request body
	var newGroup Group
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newGroup)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check if a group with the same ID already exists
	for _, group := range groups {
		if group.ID == newGroup.ID {
			http.Error(w, "Group with the same ID already exists", http.StatusConflict)
			return
		}
	}

	// Add the new group to the global "groups" array
	groups = append(groups, newGroup)

	// Send a response with the created group
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGroup)
}

func addGroupCloudEventHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a CloudEvent
	if r.Header.Get("Content-Type") == "application/cloudevents+json" {
		// Parse the CloudEvent
		event := &cloudevents.Event{}
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Extract the group and definitions from the CloudEvent data
		var newGroup Group
		if err := json.Unmarshal(event.Data(), &newGroup); err != nil {
			http.Error(w, "Invalid CloudEvent data", http.StatusBadRequest)
			return
		}

		// Check if a group with the same ID already exists
		for _, group := range groups {
			if group.ID == newGroup.ID {
				http.Error(w, "Group with the same ID already exists", http.StatusConflict)
				return
			}
		}

		// Add the new group to the global "groups" array
		groups = append(groups, newGroup)

		// Send a response with the created group
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newGroup)

	} else if r.Header.Get("Content-Type") == "application/json" {
		// The request is JSON, use the existing addGroupHandler function
		addGroupHandler(w, r)

	} else {
		http.Error(w, "Invalid content type, expected application/json or application/cloudevents+json", http.StatusUnsupportedMediaType)
		return
	}
}

func main() {

	m := http.NewServeMux()

	m.HandleFunc("/groups", groupsHandler)
	m.HandleFunc("/groups/", groupByIDHandler)
	m.HandleFunc("/", addGroupCloudEventHandler)

	//http.HandleFunc("/groups/", groupDefinitionsHandler)
	//http.HandleFunc("/groups/", groupDefinitionByIDHandler)
	http.ListenAndServe(":8080", m)
}
