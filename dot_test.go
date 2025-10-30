package erd

import (
	"strings"
	"testing"
)

func TestGetDOTCardinality(t *testing.T) {
	tests := []struct {
		name        string
		cardinality Cardinality
		want        string
	}{
		{"one to one", OneToOne, "arrowhead=normal, arrowtail=normal, dir=both"},
		{"one to many", OneToMany, "arrowhead=crow, arrowtail=normal, dir=both"},
		{"many to one", ManyToOne, "arrowhead=normal, arrowtail=crow, dir=both"},
		{"many to many", ManyToMany, "arrowhead=crow, arrowtail=crow, dir=both"},
		{"default/unknown", Cardinality("unknown"), "arrowhead=normal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDOTCardinality(tt.cardinality); got != tt.want {
				t.Errorf("getDOTCardinality() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDOTRelationship_WithLabel(t *testing.T) {
	label := "custom label"
	rel := NewRelationship("User", "Post", "Posts", OneToMany).WithLabel(label)
	output := formatDOTRelationship(rel)

	if !strings.Contains(output, label) {
		t.Errorf("formatDOTRelationship() should contain label %q, got %q", label, output)
	}
}

func TestFormatDOTAttribute_WithNullable(t *testing.T) {
	attr := NewAttribute("Name", "string").WithNullable()
	output := formatDOTAttribute(attr)

	if !strings.Contains(output, "?") {
		t.Errorf("formatDOTAttribute() should contain '?' for nullable, got %q", output)
	}
}

func TestFormatDOTAttribute_WithoutKey(t *testing.T) {
	attr := NewAttribute("Name", "string")
	output := formatDOTAttribute(attr)

	// Should not start with PK/FK/UK
	if strings.HasPrefix(output, "PK") || strings.HasPrefix(output, "FK") || strings.HasPrefix(output, "UK") {
		t.Errorf("formatDOTAttribute() should not have key prefix for non-key field, got %q", output)
	}
	if !strings.Contains(output, "Name") {
		t.Errorf("formatDOTAttribute() should contain field name, got %q", output)
	}
}

func TestEscapeDOT(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"backslash", "test\\value", "test\\\\value"},
		{"quote", `test"value`, `test\"value`},
		{"newline", "test\nvalue", "test\\nvalue"},
		{"multiple", "test\"with\\new\nline", "test\\\"with\\\\new\\nline"},
		{"no special chars", "test value", "test value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeDOT(tt.input); got != tt.want {
				t.Errorf("escapeDOT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDOT_WithEmptyTitle(t *testing.T) {
	diagram := &Diagram{
		Title:    "",
		Entities: make(map[string]*Entity),
	}
	diagram.AddEntity(NewEntity("User").AddAttribute(NewAttribute("ID", "string")))

	output := diagram.ToDOT()

	if strings.Contains(output, "labelloc") {
		t.Error("ToDOT() should not include labelloc for empty title")
	}
	// Note: "label=" can appear in entity labels, so we check for title label specifically
	if strings.Contains(output, "label=\"\"") || (strings.Contains(output, "    label=") && !strings.Contains(output, "[label=")) {
		t.Error("ToDOT() should not include title label for empty title")
	}
}
