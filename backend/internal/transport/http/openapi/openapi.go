// Package openapi serves the StreamPulse OpenAPI 3.0 specification.
//
// The spec is hand-maintained as YAML alongside the handlers. Using a static
// spec keeps the binary lean and avoids requiring `swag` at build time. For a
// fully generated approach, run:
//
//	go install github.com/swaggo/swag/cmd/swag@latest
//	swag init -g cmd/api/main.go -o internal/transport/http/openapi/generated
//
// then point Handler at `generated/swagger.json`.
package openapi

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed spec.yaml
var spec []byte

// SpecHandler returns the OpenAPI document.
func SpecHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", spec)
	}
}

// UIHandler returns a Swagger UI page that loads /openapi.yaml.
func UIHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUI))
	}
}

const swaggerUI = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>StreamPulse API · Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
  <style>html,body{margin:0;background:#fafafa}</style>
</head>
<body>
  <div id="ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.onload = () => {
      SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: "#ui",
        deepLinking: true,
        layout: "BaseLayout",
        presets: [SwaggerUIBundle.presets.apis]
      });
    };
  </script>
</body>
</html>`
