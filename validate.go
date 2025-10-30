package erd

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate checks the diagram for structural validity.
func (d *Diagram) Validate() []ValidationError {
	var errors []ValidationError

	// Check title is present
	if strings.TrimSpace(d.Title) == "" {
		errors = append(errors, ValidationError{
			Field:   "Title",
			Message: "diagram title is required",
		})
	}

	// Check entities exist
	if len(d.Entities) == 0 {
		errors = append(errors, ValidationError{
			Field:   "Entities",
			Message: "diagram must have at least one entity",
		})
	}

	// Validate each entity
	for name, entity := range d.Entities {
		if entityErrors := entity.Validate(); len(entityErrors) > 0 {
			for _, err := range entityErrors {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("Entity[%s].%s", name, err.Field),
					Message: err.Message,
				})
			}
		}
	}

	// Validate relationships reference valid entities
	for i, rel := range d.Relationships {
		if relErrors := rel.ValidateAgainst(d.Entities); len(relErrors) > 0 {
			for _, err := range relErrors {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("Relationship[%d].%s", i, err.Field),
					Message: err.Message,
				})
			}
		}
	}

	return errors
}

// Validate checks the entity for structural validity.
func (e *Entity) Validate() []ValidationError {
	var errors []ValidationError

	// Check name is present
	if strings.TrimSpace(e.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "entity name is required",
		})
	}

	// Check attributes exist
	if len(e.Attributes) == 0 {
		errors = append(errors, ValidationError{
			Field:   "Attributes",
			Message: "entity must have at least one attribute",
		})
	}

	// Validate each attribute
	for i, attr := range e.Attributes {
		if attrErrors := attr.Validate(); len(attrErrors) > 0 {
			for _, err := range attrErrors {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("Attribute[%d].%s", i, err.Field),
					Message: err.Message,
				})
			}
		}
	}

	return errors
}

// Validate checks the attribute for structural validity.
func (a *Attribute) Validate() []ValidationError {
	var errors []ValidationError

	// Check name is present
	if strings.TrimSpace(a.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "attribute name is required",
		})
	}

	// Check type is present
	if strings.TrimSpace(a.Type) == "" {
		errors = append(errors, ValidationError{
			Field:   "Type",
			Message: "attribute type is required",
		})
	}

	// Validate key type if present
	if a.Key != nil {
		if !isValidKeyType(*a.Key) {
			errors = append(errors, ValidationError{
				Field:   "Key",
				Message: fmt.Sprintf("invalid key type: %s", *a.Key),
			})
		}
	}

	return errors
}

// ValidateAgainst checks the relationship references valid entities.
func (r *Relationship) ValidateAgainst(entities map[string]*Entity) []ValidationError {
	var errors []ValidationError

	// Check from entity exists
	if strings.TrimSpace(r.From) == "" {
		errors = append(errors, ValidationError{
			Field:   "From",
			Message: "from entity is required",
		})
	} else if _, exists := entities[r.From]; !exists {
		errors = append(errors, ValidationError{
			Field:   "From",
			Message: fmt.Sprintf("entity '%s' does not exist", r.From),
		})
	}

	// Check to entity exists
	if strings.TrimSpace(r.To) == "" {
		errors = append(errors, ValidationError{
			Field:   "To",
			Message: "to entity is required",
		})
	} else if _, exists := entities[r.To]; !exists {
		errors = append(errors, ValidationError{
			Field:   "To",
			Message: fmt.Sprintf("entity '%s' does not exist", r.To),
		})
	}

	// Check field is present
	if strings.TrimSpace(r.Field) == "" {
		errors = append(errors, ValidationError{
			Field:   "Field",
			Message: "field name is required",
		})
	}

	// Validate cardinality
	if !isValidCardinality(r.Cardinality) {
		errors = append(errors, ValidationError{
			Field:   "Cardinality",
			Message: fmt.Sprintf("invalid cardinality: %s", r.Cardinality),
		})
	}

	return errors
}

// isValidKeyType checks if a key type is valid.
func isValidKeyType(kt KeyType) bool {
	switch kt {
	case PrimaryKey, ForeignKey, UniqueKey:
		return true
	default:
		return false
	}
}

// isValidCardinality checks if a cardinality is valid.
func isValidCardinality(c Cardinality) bool {
	switch c {
	case OneToOne, OneToMany, ManyToOne, ManyToMany:
		return true
	default:
		return false
	}
}
