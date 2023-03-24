package registry

import (
	"encoding/json"
	"github.com/perimeterx/marshmallow"
	"io"
	"net/http"
	"strings"
)

var groups = []Group{}

func GroupsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		renderGroups(w)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func renderGroups(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func GroupByIDHandler(w http.ResponseWriter, r *http.Request) {
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

func ParseTargetEventTypes(urls []string) {
	for _, url := range urls {
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

	var g Group
	_, err = marshmallow.Unmarshal(body, &g)
	if err != nil {
		panic(err)
	}

	groups = append(groups, g)
}
