# API

## GraphQL

### Request

```http request
POST /graphql
Content-Type: application/json
```

[GraphQL schema](../schema.graphql) is a document 

## Upload a File

### Request

```http request
POST /files

file1.png
```

### Response

```
Status: 201 Created
Content-Type: application/json; charset=utf-8

{
    "fileName": "bd693ed1-b2e3-42d8-80d6-a7696847939f.png"
}
```

## Download a File

### Request

```http request
GET /files/bd693ed1-b2e3-42d8-80d6-a7696847939f.png
```
