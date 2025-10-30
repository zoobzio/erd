[![CI Status](https://github.com/zoobzio/erd/workflows/CI/badge.svg)](https://github.com/zoobzio/erd/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zoobzio/erd/graph/badge.svg?branch=main)](https://codecov.io/gh/zoobzio/erd)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoobzio/erd)](https://goreportcard.com/report/github.com/zoobzio/erd)
[![CodeQL](https://github.com/zoobzio/erd/workflows/CodeQL/badge.svg)](https://github.com/zoobzio/erd/security/code-scanning)
[![Go Reference](https://pkg.go.dev/badge/github.com/zoobzio/erd.svg)](https://pkg.go.dev/github.com/zoobzio/erd)
[![License](https://img.shields.io/github/license/zoobzio/erd)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zoobzio/erd)](go.mod)
[![Release](https://img.shields.io/github/v/release/zoobzio/erd)](https://github.com/zoobzio/erd/releases)

# erd

Go package for defining and generating Entity Relationship Diagrams from domain models.

## Overview

`erd` provides a declarative way to define domain model relationships using Go structs, following the same pattern as the `openapi` and `dbml` packages. Unlike DBML (which models database schemas), ERD models Go type relationships - struct fields, embeddings, and collections.

## Features

- **Declarative API**: Define entities, attributes, and relationships using builder patterns
- **Multiple Output Formats**: Generate Mermaid and GraphViz DOT diagrams
- **Validation**: Built-in validation for diagram structure
- **Type-safe**: Strongly-typed cardinality and key constraints

## Installation

```bash
go get github.com/zoobzio/erd
```

## Usage

### Creating a Diagram

```go
diagram := erd.NewDiagram("User Management System").
    WithDescription("Core entities for user and order management")

// Define entities
user := erd.NewEntity("User").
    AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
    AddAttribute(erd.NewAttribute("Email", "string").WithUnique()).
    AddAttribute(erd.NewAttribute("Name", "string"))

profile := erd.NewEntity("Profile").
    AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
    AddAttribute(erd.NewAttribute("Bio", "string").WithNullable())

diagram.
    AddEntity(user).
    AddEntity(profile).
    AddRelationship(erd.NewRelationship("User", "Profile", "Profile", erd.OneToOne))
```

### Generating Output

```go
// Generate Mermaid diagram
mermaid := diagram.ToMermaid()
fmt.Println(mermaid)

// Generate GraphViz DOT diagram
dot := diagram.ToDOT()
fmt.Println(dot)
```

### Validation

```go
errors := diagram.Validate()
if len(errors) > 0 {
    for _, err := range errors {
        fmt.Println(err)
    }
}
```

## Key Concepts

### Entity
Represents a domain model entity (typically a Go struct):
- Name
- Package (optional)
- Attributes
- Notes

### Attribute
Represents a field/property of an entity:
- Name and Type
- Key constraints (PK, FK, UK)
- Nullable flag
- Notes

### Relationship
Represents relationships between entities:
- From/To entities
- Field name
- Cardinality (one-to-one, one-to-many, many-to-one, many-to-many)
- Optional label and notes

### Cardinality Types

- `OneToOne`: Single reference (`User.Profile`)
- `OneToMany`: Collection (`User.Orders`)
- `ManyToOne`: Reverse of one-to-many
- `ManyToMany`: Bi-directional collection

## Integration with Sentinel

This package is designed to work with [sentinel](https://github.com/zoobzio/sentinel) for automatic ERD generation from Go types:

```go
// sentinel extracts type metadata -> populates erd structs -> generates diagrams
```

## License

MIT
