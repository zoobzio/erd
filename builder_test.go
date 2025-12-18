package erd

import "testing"

const (
	testNote    = "test note"
	customLabel = "custom label"
)

func TestBuilderMethods(t *testing.T) {
	t.Run("Entity.WithNote", func(t *testing.T) {
		note := testNote
		entity := NewEntity("User").WithNote(note)
		if entity.Note == nil || *entity.Note != note {
			t.Errorf("WithNote() = %v, want %v", entity.Note, note)
		}
	})

	t.Run("Attribute.WithKey", func(t *testing.T) {
		attr := NewAttribute("ID", "string").WithKey(PrimaryKey)
		if attr.Key == nil || *attr.Key != PrimaryKey {
			t.Errorf("WithKey() = %v, want %v", attr.Key, PrimaryKey)
		}
	})

	t.Run("Attribute.WithNote", func(t *testing.T) {
		note := "attribute note"
		attr := NewAttribute("ID", "string").WithNote(note)
		if attr.Note == nil || *attr.Note != note {
			t.Errorf("WithNote() = %v, want %v", attr.Note, note)
		}
	})

	t.Run("Relationship.WithLabel", func(t *testing.T) {
		label := customLabel
		rel := NewRelationship("User", "Post", "Posts", OneToMany).WithLabel(label)
		if rel.Label == nil || *rel.Label != label {
			t.Errorf("WithLabel() = %v, want %v", rel.Label, label)
		}
	})

	t.Run("Relationship.WithNote", func(t *testing.T) {
		note := "relationship note"
		rel := NewRelationship("User", "Post", "Posts", OneToMany).WithNote(note)
		if rel.Note == nil || *rel.Note != note {
			t.Errorf("WithNote() = %v, want %v", rel.Note, note)
		}
	})
}
