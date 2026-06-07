// Package docs holds the OpenAPI specification served by the gateway under
// the /swagger endpoint.
package docs

import _ "embed"

// OpenAPISpec is the raw OpenAPI 3.0 document describing the gateway REST API.
//
//go:embed openapi.json
var OpenAPISpec []byte
