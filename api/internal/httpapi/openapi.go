package httpapi

import "fmt"

func buildOpenAPISpecJSON(mode string, defaultTTLSeconds int) []byte {
	if defaultTTLSeconds <= 0 {
		defaultTTLSeconds = 900
	}
	spec := fmt.Sprintf(`{
  "openapi": "3.1.0",
  "info": {
    "title": "Splatmaker API",
    "version": "0.1.0",
    "description": "Core-first API module. Works in standalone mode and is reusable by AWS adapters."
  },
  "servers": [{"url": "/"}],
  "security": [{"bearerAuth": []}],
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "JobSummary": {
        "type": "object",
        "properties": {
          "job_id": {"type": "string"},
          "status": {"type": "string", "enum": ["new", "in_progress", "done", "failed"]},
          "created_at": {"type": "string", "format": "date-time"},
          "updated_at": {"type": "string", "format": "date-time"}
        },
        "required": ["job_id", "status", "created_at", "updated_at"]
      },
      "ResultFileURL": {
        "type": "object",
        "properties": {
          "key": {"type": "string"},
          "file_name": {"type": "string"},
          "url": {"type": "string"},
          "expires_at": {"type": "string", "format": "date-time"}
        },
        "required": ["key", "file_name", "url", "expires_at"]
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {"type": "string"}
        },
        "required": ["error"]
      }
    }
  },
  "paths": {
    "/healthz": {
      "get": {
        "summary": "Health check",
        "security": [],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {"type": "string"},
                    "mode": {"type": "string"}
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/jobs": {
      "get": {
        "summary": "List current user jobs",
        "parameters": [
          {"name": "limit", "in": "query", "schema": {"type": "integer", "default": 20, "minimum": 1, "maximum": 200}},
          {"name": "offset", "in": "query", "schema": {"type": "integer", "default": 0, "minimum": 0}}
        ],
        "responses": {
          "200": {
            "description": "Jobs list",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "items": {"type": "array", "items": {"$ref": "#/components/schemas/JobSummary"}}
                  },
                  "required": ["items"]
                }
              }
            }
          },
          "401": {"description": "Unauthorized", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}}}}
        }
      }
    },
    "/v1/jobs/{job_id}/result-urls": {
      "get": {
        "summary": "Resolve result file URLs for a job",
        "parameters": [
          {"name": "job_id", "in": "path", "required": true, "schema": {"type": "string"}},
          {"name": "ttl_seconds", "in": "query", "schema": {"type": "integer", "default": %d, "minimum": 1}}
        ],
        "responses": {
          "200": {
            "description": "Result URLs",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "items": {"type": "array", "items": {"$ref": "#/components/schemas/ResultFileURL"}}
                  },
                  "required": ["items"]
                }
              }
            }
          },
          "401": {"description": "Unauthorized", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}}}},
          "404": {"description": "Not found", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}}}}
        }
      }
    }
  },
  "x-splatmaker-runtime": {
    "mode": %q
  }
}
`, defaultTTLSeconds, mode)
	return []byte(spec)
}
