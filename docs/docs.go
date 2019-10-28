// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-10-28 09:46:19.293384 +0800 CST m=+0.056374399

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/customers": {
            "get": {
                "description": "get customer list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "客户"
                ],
                "summary": "获取客户列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "联系电话",
                        "name": "mobile",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page_no",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "页数",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/customers/detail": {
            "get": {
                "description": "get customer detail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "客户"
                ],
                "summary": "获取客户详情",
                "parameters": [
                    {
                        "type": "string",
                        "description": "客户ID",
                        "name": "customer_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/files/download": {
            "get": {
                "description": "download file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文件"
                ],
                "summary": "下载文件",
                "parameters": [
                    {
                        "type": "string",
                        "description": "filename",
                        "name": "filename",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "\"avatar\", \"poster\"",
                        "description": "file type",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/files/upload": {
            "post": {
                "description": "upload file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文件"
                ],
                "summary": "上传文件",
                "parameters": [
                    {
                        "type": "file",
                        "description": "upload file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "\"avatar\", \"poster\"",
                        "description": "file type",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/merchants": {
            "get": {
                "description": "get merchant list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商户"
                ],
                "summary": "获取商户列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "商户名",
                        "name": "store_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "联系人",
                        "name": "contact_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "联系电话",
                        "name": "contact_phone",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page_no",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "页数",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "create new merchant",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商户"
                ],
                "summary": "新增商户",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.MerchantVO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/merchants/detail": {
            "get": {
                "description": "get merchant detail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商户"
                ],
                "summary": "获取商户详情信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "商户token",
                        "name": "access_token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/merchants/login": {
            "post": {
                "description": "merchant login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商户"
                ],
                "summary": "商户登录",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.MerchantLoginVO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/merchants/writeoff": {
            "get": {
                "description": "get merchant write off",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商户"
                ],
                "summary": "获取商户核销信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "商户token",
                        "name": "access_token",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "客户ID",
                        "name": "customer_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "merchant exec write off",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "商户"
                ],
                "summary": "商户核销",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.MerchantExecWriteOffVO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "user login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "后台用户"
                ],
                "summary": "后台用户登录",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "args",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.UserVO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        },
        "/verify_code": {
            "get": {
                "description": "send verify code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "验证码"
                ],
                "summary": "发送验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "手机号",
                        "name": "mobile",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":true}",
                        "schema": {
                            "$ref": "#/definitions/http.baseResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.baseResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {
                    "type": "object"
                },
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "model.MerchantExecWriteOffVO": {
            "type": "object",
            "required": [
                "access_token"
            ],
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "商户token"
                }
            }
        },
        "model.MerchantLoginVO": {
            "type": "object",
            "required": [
                "code",
                "contact_phone"
            ],
            "properties": {
                "code": {
                    "type": "string",
                    "example": "验证码"
                },
                "contact_phone": {
                    "type": "string",
                    "example": "手机号"
                }
            }
        },
        "model.MerchantVO": {
            "type": "object",
            "required": [
                "address",
                "contact_phone",
                "lat",
                "lon",
                "store_avatar",
                "store_name",
                "total_receive"
            ],
            "properties": {
                "address": {
                    "description": "地址",
                    "type": "string"
                },
                "catering_type": {
                    "description": "餐饮类型",
                    "type": "string"
                },
                "checkin_days": {
                    "description": "签到天数多少天可领取礼品",
                    "type": "integer"
                },
                "checkin_num": {
                    "description": "达到指定签到天数后，可领取的礼品数量",
                    "type": "integer"
                },
                "contact_name": {
                    "description": "联系人",
                    "type": "string"
                },
                "contact_phone": {
                    "description": "联系人电话",
                    "type": "string"
                },
                "lat": {
                    "description": "纬度",
                    "type": "number"
                },
                "lon": {
                    "description": "经度",
                    "type": "number"
                },
                "poster": {
                    "description": "商户海报",
                    "type": "string"
                },
                "received": {
                    "description": "已领礼品数量",
                    "type": "integer"
                },
                "store_avatar": {
                    "description": "店铺头像",
                    "type": "string"
                },
                "store_name": {
                    "description": "店名",
                    "type": "string"
                },
                "total_receive": {
                    "description": "该店礼品一共可领取总数",
                    "type": "integer"
                }
            }
        },
        "model.UserVO": {
            "type": "object",
            "required": [
                "name",
                "password"
            ],
            "properties": {
                "name": {
                    "description": "用户名",
                    "type": "string"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "福利签API文档",
	Description: "福利签API文档",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
