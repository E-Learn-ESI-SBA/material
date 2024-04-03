// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Seif Hanachi",
            "url": "http://www.swagger.io/support",
            "email": "s.hannachi@esi-sba.dz"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/courses/admin": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Protected Route used to get the courses (chapters) by admin id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Courses"
                ],
                "summary": "Getting Course By Admin",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Course"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/courses/create": {
            "post": {
                "description": "Protected Route used to create a course (chapter)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Courses"
                ],
                "summary": "Create Course",
                "parameters": [
                    {
                        "description": "Course Object",
                        "name": "course",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/courses/delete/{id}": {
            "delete": {
                "description": "Protected Route used to delete a course (chapter)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Courses"
                ],
                "summary": "Delete Course",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Course ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/courses/teacher": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Protected Route used to get the courses (chapters) by teacher id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Courses"
                ],
                "summary": "Getting Course By teacher",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Course"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/courses/update": {
            "put": {
                "description": "Protected Route used to update a course (chapter)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Courses"
                ],
                "summary": "Update Course",
                "parameters": [
                    {
                        "description": "Course Object",
                        "name": "course",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Course"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/modules/create": {
            "post": {
                "description": "Protected Route used to create a module",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Modules"
                ],
                "summary": "Create Module",
                "parameters": [
                    {
                        "description": "Module Object",
                        "name": "module",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/modules/delete/{id}": {
            "delete": {
                "description": "Protected Route used to delete a module",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Modules"
                ],
                "summary": "Delete Module",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Module ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/modules/public": {
            "post": {
                "description": "Protected Route used to get public modules",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Modules"
                ],
                "summary": "Get Public Modules",
                "parameters": [
                    {
                        "description": "Module Filter",
                        "name": "filter",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/interfaces.ModuleFilter"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/modules/teacher": {
            "get": {
                "description": "Protected Route used to get teacher modules",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Modules"
                ],
                "summary": "Get Teacher Modules",
                "parameters": [
                    {
                        "description": "Module Filter",
                        "name": "filter",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/interfaces.ModuleFilter"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/modules/update": {
            "put": {
                "description": "Protected Route used to update a module",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Modules"
                ],
                "summary": "Update Module",
                "parameters": [
                    {
                        "description": "Module Object",
                        "name": "module",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/modules/visibility/{id}": {
            "put": {
                "description": "Protected Route used to edit module visibility",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Modules"
                ],
                "summary": "Edit Module Visibility",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Module ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Module Visibility",
                        "name": "visibility",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Module"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        },
        "/modules/{id}": {
            "get": {
                "description": "Get Module By ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Modules"
                ],
                "summary": "Get Module By ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Module ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ExtendedModule"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/interfaces.APiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "interfaces.APiError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "interfaces.ModuleFilter": {
            "type": "object",
            "properties": {
                "semester": {
                    "type": "integer",
                    "maximum": 2,
                    "minimum": 1
                },
                "speciality": {
                    "type": "string"
                },
                "year": {
                    "type": "integer",
                    "maximum": 5,
                    "minimum": 1
                }
            }
        },
        "models.Course": {
            "type": "object",
            "required": [
                "plan"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "module_id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "plan": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.ExtendCourse": {
            "type": "object",
            "required": [
                "plan"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "module_id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "plan": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                },
                "sections": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Section"
                    }
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.ExtendedModule": {
            "type": "object",
            "required": [
                "teacher_id"
            ],
            "properties": {
                "coefficient": {
                    "type": "integer"
                },
                "courses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ExtendCourse"
                    }
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "instructors": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "isPublic": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "semester": {
                    "type": "integer"
                },
                "speciality": {
                    "type": "string"
                },
                "teacher_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "models.Files": {
            "type": "object",
            "required": [
                "group",
                "section_id",
                "teacher_id",
                "url"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "group": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "section_id": {
                    "type": "string"
                },
                "teacher_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "models.Lecture": {
            "type": "object",
            "required": [
                "content",
                "section_id",
                "teacher_id"
            ],
            "properties": {
                "content": {
                    "type": "string",
                    "minLength": 250
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_public": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "section_id": {
                    "type": "string"
                },
                "teacher_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.Module": {
            "type": "object",
            "required": [
                "teacher_id"
            ],
            "properties": {
                "coefficient": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "instructors": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "isPublic": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "semester": {
                    "type": "integer"
                },
                "speciality": {
                    "type": "string"
                },
                "teacher_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "models.Section": {
            "type": "object",
            "required": [
                "course_id",
                "name",
                "teacher_id"
            ],
            "properties": {
                "course_id": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Files"
                    }
                },
                "id": {
                    "type": "string"
                },
                "lectures": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Lecture"
                    }
                },
                "name": {
                    "type": "string"
                },
                "order": {
                    "type": "integer",
                    "minimum": 1
                },
                "teacher_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "videos": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Video"
                    }
                }
            }
        },
        "models.Video": {
            "type": "object",
            "required": [
                "section_id",
                "teacher_id",
                "url"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "section_id": {
                    "type": "string"
                },
                "teacher_id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Madaurus Material service",
	Description:      "This Service is for managing the material of the Madaurus Platform",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
