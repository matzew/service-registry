package endpoint

import (
	"encoding/json"
	"net/http"
)

// // Define a map to store the endpoints
var endpoints = []Endpoint{
	{
		ID:    "dev.knative.kafkasource.test-namespace.source1",
		Usage: "consumer",
		Config: EndpointConfig{
			Protocol: "Kafka/3.4",
			Endpoints: []string{
				"my-cluster-kafka-bootstrap.kafka:9092",
			},
			Options: map[string]string{
				"topic":          "mytopic",
				"consumer-group": "source1-cg",
			},
		},
		Channel:         "amq-streams.something.mytopic-consumer",
		Format:          "CloudEvents/1.0",
		DefinitionGroup: []string{"myserver/definitionGroups/my.custom.eventdef.grp"},
	},
	{
		ID:    "kafka.mytopic-consumer",
		Usage: "consumer",
		Config: EndpointConfig{
			Protocol: "Kafka/3.4",
			Endpoints: []string{
				"my-cluster-kafka-bootstrap.kafka:9092",
			},
			Options: map[string]string{
				"topic":          "mytopic",
				"consumer-group": "source1-cg",
			},
		},
	},
}

func EndpointsHandler(w http.ResponseWriter, r *http.Request) {

	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the 'endpoints' map as JSON and write it to the response
	json.NewEncoder(w).Encode(endpoints)
}
