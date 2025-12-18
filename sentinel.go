package erd

import (
	"strings"

	"github.com/zoobzio/sentinel"
)

func init() {
	sentinel.Tag("erd")
}

// FromSchema converts a sentinel schema to an ERD diagram.
// The schema is typically obtained via sentinel.Schema() after scanning types.
func FromSchema(title string, schema map[string]sentinel.ModelMetadata) *Diagram {
	diagram := NewDiagram(title)

	// Add entities, filtering out relationship fields
	for _, meta := range schema {
		entity := fromMetadataFiltered(meta)
		diagram.AddEntity(entity)
	}

	// Add relationships
	for _, meta := range schema {
		for _, rel := range meta.Relationships {
			diagram.AddRelationship(relationshipFromSentinel(rel))
		}
	}

	return diagram
}

// fromMetadataFiltered converts sentinel ModelMetadata to an ERD Entity,
// filtering out fields that are represented as relationships.
func fromMetadataFiltered(meta sentinel.ModelMetadata) *Entity {
	// Build set of relationship field names
	relFields := make(map[string]bool)
	for _, rel := range meta.Relationships {
		relFields[rel.Field] = true
	}

	entity := NewEntity(meta.TypeName)

	if meta.PackageName != "" {
		entity.WithPackage(meta.PackageName)
	}

	for _, field := range meta.Fields {
		// Skip fields that are relationships (shown as lines, not attributes)
		if relFields[field.Name] {
			continue
		}
		attr := attributeFromField(field)
		entity.AddAttribute(attr)
	}

	return entity
}

// FromMetadata converts a single sentinel ModelMetadata to an ERD Entity.
func FromMetadata(meta sentinel.ModelMetadata) *Entity {
	entity := NewEntity(meta.TypeName)

	if meta.PackageName != "" {
		entity.WithPackage(meta.PackageName)
	}

	for _, field := range meta.Fields {
		attr := attributeFromField(field)
		entity.AddAttribute(attr)
	}

	return entity
}

// attributeFromField converts a sentinel FieldMetadata to an ERD Attribute.
func attributeFromField(field sentinel.FieldMetadata) *Attribute {
	attrType := field.Type

	// Detect nullability from pointer types
	nullable := strings.HasPrefix(field.Type, "*")
	if nullable {
		attrType = strings.TrimPrefix(field.Type, "*")
	}

	attr := NewAttribute(field.Name, attrType)

	if nullable {
		attr.WithNullable()
	}

	// Parse erd tag for key types and notes
	if erdTag, ok := field.Tags["erd"]; ok {
		parseErdTag(attr, erdTag)
	}

	return attr
}

// parseErdTag parses the erd struct tag and applies settings to the attribute.
// Supported values: pk, fk, uk, note:...
func parseErdTag(attr *Attribute, tag string) {
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case part == "pk":
			attr.WithPrimaryKey()
		case part == "fk":
			attr.WithForeignKey()
		case part == "uk":
			attr.WithUnique()
		case strings.HasPrefix(part, "note:"):
			note := strings.TrimPrefix(part, "note:")
			attr.WithNote(note)
		}
	}
}

// relationshipFromSentinel converts a sentinel TypeRelationship to an ERD Relationship.
func relationshipFromSentinel(rel sentinel.TypeRelationship) *Relationship {
	cardinality := cardinalityFromKind(rel.Kind)
	return NewRelationship(rel.From, rel.To, rel.Field, cardinality)
}

// cardinalityFromKind maps sentinel relationship kinds to ERD cardinalities.
func cardinalityFromKind(kind string) Cardinality {
	switch kind {
	case sentinel.RelationshipReference:
		return OneToOne
	case sentinel.RelationshipCollection:
		return OneToMany
	case sentinel.RelationshipEmbedding:
		return OneToOne
	case sentinel.RelationshipMap:
		return ManyToMany
	default:
		return OneToOne
	}
}
