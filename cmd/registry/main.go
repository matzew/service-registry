package main

import (
	"github.com/matzew/service-registry/pkg/endpoint"
	"github.com/matzew/service-registry/pkg/registry"
	"net/http"
)

// TODO: read from configmap
var federatedUrls = []string{
	// Clement's sample data
	"https://cediscoveryinterop.azurewebsites.net/registry/definitiongroups/microsoft.eventgrid?code=VyRba-8O1MYZ6EPVoV34u3RSpZXqnlaefnVTsxT7p0BLAzFuCziSzw==",
	"https://cediscoveryinterop.azurewebsites.net/registry/definitiongroups/microsoft.eventhub?code=VyRba-8O1MYZ6EPVoV34u3RSpZXqnlaefnVTsxT7p0BLAzFuCziSzw==",

	// Matthias' sample data
	"https://raw.githubusercontent.com/matzew/service-registry/main/sample-data/knative_k8s_source.json",
	"https://raw.githubusercontent.com/matzew/service-registry/main/sample-data/knative_ping_source.json",
	"https://raw.githubusercontent.com/matzew/service-registry/main/sample-data/telegram_sample_group.json",
	"https://raw.githubusercontent.com/matzew/service-registry/main/sample-data/wa_sample_group.json",
}

func main() {

	registry.ParseTargetEventTypes(federatedUrls)

	m := http.NewServeMux()

	m.HandleFunc("/groups", registry.GroupsHandler)
	m.HandleFunc("/groups/", registry.GroupByIDHandler)

	m.HandleFunc("/endpoints", endpoint.EndpointsHandler)

	http.ListenAndServe(":8080", m)
}
