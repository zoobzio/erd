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
