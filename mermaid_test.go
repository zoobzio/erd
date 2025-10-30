package erd

import (
	"strings"
	"testing"
)

func TestGetMermaidCardinality(t *testing.T) {
	tests := []struct {
		name        string
		cardinality Cardinality
		want        string
	}{
		{"one to one", OneToOne, "||--||"},
		{"one to many", OneToMany, "||--o{"},
		{"many to one", ManyToOne, "}o--||"},
		{"many to many", ManyToMany, "}o--o{"},
		{"default/unknown", Cardinality("unknown"), "||--||"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMermaidCardinality(tt.cardinality); got != tt.want {
				t.Errorf("getMermaidCardinality() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatMermaidRelationship_WithLabel(t *testing.T) {
	label := "custom label"
	rel := NewRelationship("User", "Post", "Posts", OneToMany).WithLabel(label)
	output := formatMermaidRelationship(rel)

	if !strings.Contains(output, label) {
		t.Errorf("formatMermaidRelationship() should contain label %q, got %q", label, output)
	}
}

func TestSanitizeType(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     string
	}{
		{"simple type", "string", "string"},
		{"with package", "github.com/user/pkg.Type", "Type"},
		{"with slice", "[]string", "Array_string"},
		{"with pointer", "*User", "User"},
		{"with space", "some type", "some_type"},
		{"complex", "[]*github.com/user/pkg.Type", "Type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeType(tt.typeName); got != tt.want {
				t.Errorf("sanitizeType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatMermaidAttribute_WithNote(t *testing.T) {
	note := "test note"
	attr := NewAttribute("ID", "string").WithNote(note)
	output := formatMermaidAttribute(attr)

	if !strings.Contains(output, note) {
		t.Errorf("formatMermaidAttribute() should contain note %q, got %q", note, output)
	}
}

func TestFormatMermaidAttribute_WithNullableAndNote(t *testing.T) {
	note := "test note"
	attr := NewAttribute("Name", "string").WithNullable().WithNote(note)
	output := formatMermaidAttribute(attr)

	if !strings.Contains(output, "nullable") {
		t.Errorf("formatMermaidAttribute() should contain 'nullable', got %q", output)
	}
	if !strings.Contains(output, note) {
		t.Errorf("formatMermaidAttribute() should contain note %q, got %q", note, output)
	}
}
