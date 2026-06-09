// Package docs contains Swagger metadata for the subscription tracker API.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "schemes": {{ marshal .Schemes }},
  "swagger": "2.0",
  "info": {
    "description": "{{escape .Description}}",
    "title": "{{.Title}}",
    "contact": {},
    "version": "{{.Version}}"
  },
  "host": "{{.Host}}",
  "basePath": "{{.BasePath}}",
  "paths": {
    "/subscriptions": {
      "get": {
        "produces": ["application/json"],
        "tags": ["subscriptions"],
        "summary": "List subscriptions",
        "parameters": [
          {"type": "string", "description": "Filter by user ID", "name": "user_id", "in": "query"},
          {"type": "integer", "default": 1, "description": "Page number", "name": "page", "in": "query"},
          {"type": "integer", "default": 20, "description": "Page size", "name": "page_size", "in": "query"}
        ],
        "responses": {"200": {"description": "OK", "schema": {"type": "object"}}}
      },
      "post": {
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "tags": ["subscriptions"],
        "summary": "Create a new subscription",
        "parameters": [
          {"description": "Subscription data", "name": "subscription", "in": "body", "required": true, "schema": {"$ref": "#/definitions/models.CreateSubscriptionRequest"}}
        ],
        "responses": {
          "201": {"description": "Created", "schema": {"$ref": "#/definitions/models.SubscriptionResponse"}},
          "400": {"description": "Bad Request", "schema": {"type": "object", "additionalProperties": {"type": "string"}}},
          "500": {"description": "Internal Server Error", "schema": {"type": "object", "additionalProperties": {"type": "string"}}}
        }
      }
    },
    "/subscriptions/{id}": {
      "get": {
        "produces": ["application/json"],
        "tags": ["subscriptions"],
        "summary": "Get subscription by ID",
        "parameters": [{"type": "string", "description": "Subscription ID", "name": "id", "in": "path", "required": true}],
        "responses": {
          "200": {"description": "OK", "schema": {"$ref": "#/definitions/models.SubscriptionResponse"}},
          "404": {"description": "Not Found", "schema": {"type": "object", "additionalProperties": {"type": "string"}}}
        }
      },
      "put": {
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "tags": ["subscriptions"],
        "summary": "Update subscription",
        "parameters": [
          {"type": "string", "description": "Subscription ID", "name": "id", "in": "path", "required": true},
          {"description": "Update data", "name": "subscription", "in": "body", "required": true, "schema": {"$ref": "#/definitions/models.UpdateSubscriptionRequest"}}
        ],
        "responses": {
          "200": {"description": "OK", "schema": {"type": "object", "additionalProperties": {"type": "string"}}},
          "400": {"description": "Bad Request", "schema": {"type": "object", "additionalProperties": {"type": "string"}}},
          "404": {"description": "Not Found", "schema": {"type": "object", "additionalProperties": {"type": "string"}}}
        }
      },
      "delete": {
        "produces": ["application/json"],
        "tags": ["subscriptions"],
        "summary": "Delete subscription",
        "parameters": [{"type": "string", "description": "Subscription ID", "name": "id", "in": "path", "required": true}],
        "responses": {
          "204": {"description": "No Content"},
          "404": {"description": "Not Found", "schema": {"type": "object", "additionalProperties": {"type": "string"}}}
        }
      }
    },
    "/calculate-cost": {
      "post": {
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "tags": ["calculations"],
        "summary": "Calculate total cost for period",
        "parameters": [
          {"description": "Period and filters", "name": "request", "in": "body", "required": true, "schema": {"$ref": "#/definitions/models.CalculateCostRequest"}}
        ],
        "responses": {
          "200": {"description": "OK", "schema": {"$ref": "#/definitions/models.CalculateCostResponse"}},
          "400": {"description": "Bad Request", "schema": {"type": "object", "additionalProperties": {"type": "string"}}}
        }
      }
    }
  },
  "definitions": {
    "models.CalculateCostRequest": {
      "type": "object",
      "required": ["start_date", "end_date"],
      "properties": {
        "end_date": {"type": "string", "example": "12-2025"},
        "service_name": {"type": "string", "example": "Yandex Plus"},
        "start_date": {"type": "string", "example": "07-2025"},
        "user_id": {"type": "string", "example": "60601fee-2bf1-4721-ae6f-7636e79a0cba"}
      }
    },
    "models.CalculateCostResponse": {
      "type": "object",
      "properties": {"total_cost": {"type": "integer", "example": 2400}}
    },
    "models.CreateSubscriptionRequest": {
      "type": "object",
      "required": ["service_name", "price", "user_id", "start_date"],
      "properties": {
        "end_date": {"type": "string", "example": "12-2025"},
        "price": {"type": "integer", "example": 400},
        "service_name": {"type": "string", "example": "Yandex Plus"},
        "start_date": {"type": "string", "example": "07-2025"},
        "user_id": {"type": "string", "example": "60601fee-2bf1-4721-ae6f-7636e79a0cba"}
      }
    },
    "models.SubscriptionResponse": {
      "type": "object",
      "properties": {
        "created_at": {"type": "string"},
        "end_date": {"type": "string"},
        "id": {"type": "string"},
        "price": {"type": "integer"},
        "service_name": {"type": "string"},
        "start_date": {"type": "string"},
        "updated_at": {"type": "string"},
        "user_id": {"type": "string"}
      }
    },
    "models.UpdateSubscriptionRequest": {
      "type": "object",
      "properties": {
        "end_date": {"type": "string", "example": "12-2025"},
        "price": {"type": "integer", "example": 500},
        "service_name": {"type": "string", "example": "Yandex Plus"},
        "start_date": {"type": "string", "example": "07-2025"}
      }
    }
  }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "Subscription Service API",
	Description:      "REST service for aggregating user online subscriptions.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
