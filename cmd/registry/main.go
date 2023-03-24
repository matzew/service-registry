package main

import (
	"encoding/json"
	"io"
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

var federatedUrls = []string{
	"https://raw.githubusercontent.com/matzew/service-registry/main/wa_sample_group.json",
	"https://raw.githubusercontent.com/matzew/service-registry/main/telegram_sample_group.json",
}

var groups = []Group{}

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

func parseTargetEventTypes() {
	for _, url := range federatedUrls {
		registerGroup(url)
	}
}

func registerGroup(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var group Group
	err = json.Unmarshal(body, &group)
	if err != nil {
		panic(err)
	}

	groups = append(groups, group)
}

func main() {

	parseTargetEventTypes()

	m := http.NewServeMux()

	m.HandleFunc("/groups", groupsHandler)
	m.HandleFunc("/groups/", groupByIDHandler)
	m.HandleFunc("/", addGroupCloudEventHandler)

	http.ListenAndServe(":8080", m)
}
