package erd

import (
	"testing"

	"github.com/zoobzio/sentinel"
)

type User struct {
	ID      string `erd:"pk"`
	Email   string `erd:"uk"`
	Name    string
	Profile *Profile
	Orders  []Order
}

type Profile struct {
	ID     string `erd:"pk"`
	Bio    *string
	Avatar *string `erd:"note:Profile picture URL"`
}

type Order struct {
	ID     string `erd:"pk"`
	UserID string `erd:"fk"`
	Total  float64
	Status string
}

// Additional types to test embedding and map relationships.
type Address struct {
	Street string
	City   string
}

type Company struct {
	ID        string `erd:"pk"`
	Name      string
	Address   // Embedded struct
	Employees map[string]Employee
}

type Employee struct {
	ID   string `erd:"pk"`
	Name string
}

func TestFromSchema(t *testing.T) {
	sentinel.Scan[User]()
	schema := sentinel.Schema()

	diagram := FromSchema("Test Domain", schema)

	if diagram.Title != "Test Domain" {
		t.Errorf("expected title 'Test Domain', got %q", diagram.Title)
	}

	if len(diagram.Entities) != 3 {
		t.Errorf("expected 3 entities, got %d", len(diagram.Entities))
	}

	// Check User entity
	user, ok := diagram.Entities["User"]
	if !ok {
		t.Fatal("User entity not found")
	}

	if len(user.Attributes) != 3 {
		t.Errorf("expected 3 attributes on User (excluding relationships), got %d", len(user.Attributes))
	}

	// Check that ID has primary key
	var idAttr *Attribute
	for _, attr := range user.Attributes {
		if attr.Name == "ID" {
			idAttr = attr
			break
		}
	}
	if idAttr == nil {
		t.Fatal("ID attribute not found on User")
	}
	if idAttr.Key == nil || *idAttr.Key != PrimaryKey {
		t.Error("expected ID to be marked as primary key")
	}

	// Check relationships
	if len(diagram.Relationships) < 2 {
		t.Errorf("expected at least 2 relationships, got %d", len(diagram.Relationships))
	}

	// Check Profile relationship is OneToOne
	var profileRel *Relationship
	for _, rel := range diagram.Relationships {
		if rel.Field == "Profile" {
			profileRel = rel
			break
		}
	}
	if profileRel == nil {
		t.Fatal("Profile relationship not found")
	}
	if profileRel.Cardinality != OneToOne {
		t.Errorf("expected Profile relationship to be OneToOne, got %s", profileRel.Cardinality)
	}

	// Check Orders relationship is OneToMany
	var ordersRel *Relationship
	for _, rel := range diagram.Relationships {
		if rel.Field == "Orders" {
			ordersRel = rel
			break
		}
	}
	if ordersRel == nil {
		t.Fatal("Orders relationship not found")
	}
	if ordersRel.Cardinality != OneToMany {
		t.Errorf("expected Orders relationship to be OneToMany, got %s", ordersRel.Cardinality)
	}
}

func TestFromMetadata(t *testing.T) {
	meta := sentinel.Inspect[Profile]()

	entity := FromMetadata(meta)

	if entity.Name != "Profile" {
		t.Errorf("expected entity name 'Profile', got %q", entity.Name)
	}

	// Check Bio is nullable (pointer type)
	var bioAttr *Attribute
	for _, attr := range entity.Attributes {
		if attr.Name == "Bio" {
			bioAttr = attr
			break
		}
	}
	if bioAttr == nil {
		t.Fatal("Bio attribute not found")
	}
	if !bioAttr.Nullable {
		t.Error("expected Bio to be nullable (pointer type)")
	}
	if bioAttr.Type != "string" {
		t.Errorf("expected Bio type to be 'string' (without pointer), got %q", bioAttr.Type)
	}

	// Check Avatar has note
	var avatarAttr *Attribute
	for _, attr := range entity.Attributes {
		if attr.Name == "Avatar" {
			avatarAttr = attr
			break
		}
	}
	if avatarAttr == nil {
		t.Fatal("Avatar attribute not found")
	}
	if avatarAttr.Note == nil || *avatarAttr.Note != "Profile picture URL" {
		t.Error("expected Avatar to have note 'Profile picture URL'")
	}
}

func TestParseErdTag(t *testing.T) {
	tests := []struct {
		tag      string
		wantKey  *KeyType
		wantNote *string
	}{
		{"pk", ptr(PrimaryKey), nil},
		{"fk", ptr(ForeignKey), nil},
		{"uk", ptr(UniqueKey), nil},
		{"pk,note:Primary identifier", ptr(PrimaryKey), ptr("Primary identifier")},
		{"note:Just a note", nil, ptr("Just a note")},
	}

	for _, tt := range tests {
		attr := NewAttribute("test", "string")
		parseErdTag(attr, tt.tag)

		if tt.wantKey == nil && attr.Key != nil {
			t.Errorf("tag %q: expected no key, got %v", tt.tag, *attr.Key)
		} else if tt.wantKey != nil && (attr.Key == nil || *attr.Key != *tt.wantKey) {
			t.Errorf("tag %q: expected key %v, got %v", tt.tag, *tt.wantKey, attr.Key)
		}

		if tt.wantNote == nil && attr.Note != nil {
			t.Errorf("tag %q: expected no note, got %q", tt.tag, *attr.Note)
		} else if tt.wantNote != nil && (attr.Note == nil || *attr.Note != *tt.wantNote) {
			t.Errorf("tag %q: expected note %q, got %v", tt.tag, *tt.wantNote, attr.Note)
		}
	}
}

func TestCardinalityFromKind(t *testing.T) {
	// Scan Company to exercise embedding and map relationships
	sentinel.Scan[Company]()
	schema := sentinel.Schema()

	diagram := FromSchema("Company Domain", schema)

	// Check Company entity exists
	company, ok := diagram.Entities["Company"]
	if !ok {
		t.Fatal("Company entity not found")
	}
	if company == nil {
		t.Fatal("Company entity is nil")
	}

	// Check for embedding relationship (Address)
	var embeddingRel *Relationship
	for _, rel := range diagram.Relationships {
		if rel.From == "Company" && rel.To == "Address" {
			embeddingRel = rel
			break
		}
	}
	if embeddingRel == nil {
		t.Fatal("Address embedding relationship not found")
	}
	if embeddingRel.Cardinality != OneToOne {
		t.Errorf("expected embedding to be OneToOne, got %s", embeddingRel.Cardinality)
	}

	// Check for map relationship (Employees)
	var mapRel *Relationship
	for _, rel := range diagram.Relationships {
		if rel.Field == "Employees" {
			mapRel = rel
			break
		}
	}
	if mapRel == nil {
		t.Fatal("Employees map relationship not found")
	}
	if mapRel.Cardinality != ManyToMany {
		t.Errorf("expected map to be ManyToMany, got %s", mapRel.Cardinality)
	}
}

func TestCardinalityFromKindAllCases(t *testing.T) {
	tests := []struct {
		kind string
		want Cardinality
	}{
		{sentinel.RelationshipReference, OneToOne},
		{sentinel.RelationshipCollection, OneToMany},
		{sentinel.RelationshipEmbedding, OneToOne},
		{sentinel.RelationshipMap, ManyToMany},
		{"unknown", OneToOne}, // default case
	}

	for _, tt := range tests {
		got := cardinalityFromKind(tt.kind)
		if got != tt.want {
			t.Errorf("cardinalityFromKind(%q) = %s, want %s", tt.kind, got, tt.want)
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}
