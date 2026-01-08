// Package erd generates Entity Relationship Diagrams from Go domain models.
//
// The package provides two approaches for generating ERD diagrams:
//
// # Automatic Generation (via Sentinel)
//
// Use [FromSchema] to automatically generate diagrams from Go types scanned
// by the sentinel package:
//
//	sentinel.Scan[User]()
//	diagram := erd.FromSchema("My Domain", sentinel.Schema())
//	fmt.Println(diagram.ToMermaid())
//
// # Manual Construction (via Builder)
//
// Use [NewDiagram], [NewEntity], [NewAttribute], and [NewRelationship] to
// construct diagrams programmatically:
//
//	diagram := erd.NewDiagram("E-commerce").
//	    AddEntity(erd.NewEntity("Product").
//	        AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey())).
//	    AddRelationship(erd.NewRelationship("Cart", "Product", "Items", erd.ManyToMany))
//
// # Output Formats
//
// Diagrams can be rendered to multiple formats:
//   - [Diagram.ToMermaid] - Mermaid syntax for web rendering
//   - [Diagram.ToDOT] - GraphViz DOT format for high-quality output
//
// # Struct Tags
//
// When using automatic generation, the erd struct tag controls attribute metadata:
//   - pk: Primary key
//   - fk: Foreign key
//   - uk: Unique key
//   - note:...: Attribute annotation
//
// Tags can be combined: `erd:"pk,note:Auto-generated UUID"`
//
// # Validation
//
// Use [Diagram.Validate] to check diagram structural validity before rendering.
package erd

// Diagram represents an Entity Relationship Diagram.
type Diagram struct {
	Title         string
	Description   *string
	Entities      map[string]*Entity
	Relationships []*Relationship
}

// Entity represents a domain model entity (typically a Go struct).
type Entity struct {
	Package    *string
	Note       *string
	Name       string
	Attributes []*Attribute
}

// Attribute represents a field/property of an entity.
type Attribute struct {
	Key      *KeyType
	Note     *string
	Name     string
	Type     string
	Nullable bool
}

// KeyType represents the key constraint on an attribute.
type KeyType string

// Key type constants.
const (
	PrimaryKey KeyType = "PK"
	ForeignKey KeyType = "FK"
	UniqueKey  KeyType = "UK"
)

// Relationship represents a relationship between entities.
type Relationship struct {
	Label       *string
	Note        *string
	From        string
	To          string
	Field       string
	Cardinality Cardinality
}

// Cardinality represents the type of relationship between entities.
type Cardinality string

// Cardinality constants.
const (
	OneToOne   Cardinality = "one-to-one"
	OneToMany  Cardinality = "one-to-many"
	ManyToOne  Cardinality = "many-to-one"
	ManyToMany Cardinality = "many-to-many"
)
