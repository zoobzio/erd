package erd

import (
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "TestField",
		Message: "test message",
	}
	expected := "TestField: test message"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestDiagram_Validate(t *testing.T) {
	tests := []struct {
		name    string
		diagram *Diagram
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid diagram",
			diagram: NewDiagram("Test").
				AddEntity(NewEntity("User").
					AddAttribute(NewAttribute("ID", "string"))),
			wantErr: false,
		},
		{
			name:    "empty title",
			diagram: NewDiagram(""),
			wantErr: true,
			errMsg:  "Title: diagram title is required",
		},
		{
			name:    "no entities",
			diagram: NewDiagram("Test"),
			wantErr: true,
			errMsg:  "Entities: diagram must have at least one entity",
		},
		{
			name: "entity validation error",
			diagram: NewDiagram("Test").
				AddEntity(NewEntity("")),
			wantErr: true,
			errMsg:  "Entity[].Name: entity name is required",
		},
		{
			name: "relationship validation error",
			diagram: NewDiagram("Test").
				AddEntity(NewEntity("User").
					AddAttribute(NewAttribute("ID", "string"))).
				AddRelationship(NewRelationship("User", "NonExistent", "Ref", OneToOne)),
			wantErr: true,
			errMsg:  "Relationship[0].To: entity 'NonExistent' does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.diagram.Validate()
			if tt.wantErr && len(errs) == 0 {
				t.Error("Validate() expected errors, got none")
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Errorf("Validate() unexpected errors: %v", errs)
			}
			if tt.wantErr && len(errs) > 0 {
				found := false
				for _, err := range errs {
					if err.Error() == tt.errMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Validate() expected error %q, got %v", tt.errMsg, errs)
				}
			}
		})
	}
}

func TestEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entity  *Entity
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid entity",
			entity: NewEntity("User").
				AddAttribute(NewAttribute("ID", "string")),
			wantErr: false,
		},
		{
			name:    "empty name",
			entity:  NewEntity(""),
			wantErr: true,
			errMsg:  "Name: entity name is required",
		},
		{
			name:    "no attributes",
			entity:  NewEntity("User"),
			wantErr: true,
			errMsg:  "Attributes: entity must have at least one attribute",
		},
		{
			name: "attribute validation error",
			entity: NewEntity("User").
				AddAttribute(NewAttribute("", "string")),
			wantErr: true,
			errMsg:  "Attribute[0].Name: attribute name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.entity.Validate()
			if tt.wantErr && len(errs) == 0 {
				t.Error("Validate() expected errors, got none")
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Errorf("Validate() unexpected errors: %v", errs)
			}
			if tt.wantErr && len(errs) > 0 {
				found := false
				for _, err := range errs {
					if err.Error() == tt.errMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Validate() expected error %q, got %v", tt.errMsg, errs)
				}
			}
		})
	}
}

func TestAttribute_Validate(t *testing.T) {
	invalidKey := KeyType("INVALID")
	tests := []struct {
		name      string
		attribute *Attribute
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid attribute",
			attribute: NewAttribute("ID", "string"),
			wantErr:   false,
		},
		{
			name:      "empty name",
			attribute: NewAttribute("", "string"),
			wantErr:   true,
			errMsg:    "Name: attribute name is required",
		},
		{
			name:      "empty type",
			attribute: NewAttribute("ID", ""),
			wantErr:   true,
			errMsg:    "Type: attribute type is required",
		},
		{
			name:      "invalid key type",
			attribute: &Attribute{Name: "ID", Type: "string", Key: &invalidKey},
			wantErr:   true,
			errMsg:    "Key: invalid key type: INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.attribute.Validate()
			if tt.wantErr && len(errs) == 0 {
				t.Error("Validate() expected errors, got none")
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Errorf("Validate() unexpected errors: %v", errs)
			}
			if tt.wantErr && len(errs) > 0 {
				found := false
				for _, err := range errs {
					if err.Error() == tt.errMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Validate() expected error %q, got %v", tt.errMsg, errs)
				}
			}
		})
	}
}

func TestRelationship_ValidateAgainst(t *testing.T) {
	entities := map[string]*Entity{
		"User": NewEntity("User").AddAttribute(NewAttribute("ID", "string")),
		"Post": NewEntity("Post").AddAttribute(NewAttribute("ID", "string")),
	}

	invalidCardinality := Cardinality("invalid")

	tests := []struct {
		name         string
		relationship *Relationship
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid relationship",
			relationship: NewRelationship("User", "Post", "Posts", OneToMany),
			wantErr:      false,
		},
		{
			name:         "empty from",
			relationship: NewRelationship("", "Post", "Posts", OneToMany),
			wantErr:      true,
			errMsg:       "From: from entity is required",
		},
		{
			name:         "empty to",
			relationship: NewRelationship("User", "", "Posts", OneToMany),
			wantErr:      true,
			errMsg:       "To: to entity is required",
		},
		{
			name:         "from entity does not exist",
			relationship: NewRelationship("NonExistent", "Post", "Posts", OneToMany),
			wantErr:      true,
			errMsg:       "From: entity 'NonExistent' does not exist",
		},
		{
			name:         "to entity does not exist",
			relationship: NewRelationship("User", "NonExistent", "Posts", OneToMany),
			wantErr:      true,
			errMsg:       "To: entity 'NonExistent' does not exist",
		},
		{
			name:         "empty field",
			relationship: NewRelationship("User", "Post", "", OneToMany),
			wantErr:      true,
			errMsg:       "Field: field name is required",
		},
		{
			name:         "invalid cardinality",
			relationship: &Relationship{From: "User", To: "Post", Field: "Posts", Cardinality: invalidCardinality},
			wantErr:      true,
			errMsg:       "Cardinality: invalid cardinality: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.relationship.ValidateAgainst(entities)
			if tt.wantErr && len(errs) == 0 {
				t.Error("ValidateAgainst() expected errors, got none")
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Errorf("ValidateAgainst() unexpected errors: %v", errs)
			}
			if tt.wantErr && len(errs) > 0 {
				found := false
				for _, err := range errs {
					if err.Error() == tt.errMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ValidateAgainst() expected error %q, got %v", tt.errMsg, errs)
				}
			}
		})
	}
}

func TestIsValidKeyType(t *testing.T) {
	tests := []struct {
		name string
		kt   KeyType
		want bool
	}{
		{"primary key", PrimaryKey, true},
		{"foreign key", ForeignKey, true},
		{"unique key", UniqueKey, true},
		{"invalid", KeyType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidKeyType(tt.kt); got != tt.want {
				t.Errorf("isValidKeyType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidCardinality(t *testing.T) {
	tests := []struct {
		name string
		c    Cardinality
		want bool
	}{
		{"one to one", OneToOne, true},
		{"one to many", OneToMany, true},
		{"many to one", ManyToOne, true},
		{"many to many", ManyToMany, true},
		{"invalid", Cardinality("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidCardinality(tt.c); got != tt.want {
				t.Errorf("isValidCardinality() = %v, want %v", got, tt.want)
			}
		})
	}
}
