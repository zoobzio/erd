package erd

import (
	"fmt"
	"sort"
	"strings"
)

// ToMermaid generates a Mermaid ERD diagram from the diagram structure.
func (d *Diagram) ToMermaid() string {
	var sb strings.Builder

	sb.WriteString("erDiagram\n")

	// Write entities with their attributes (sorted by name for deterministic output)
	entityNames := make([]string, 0, len(d.Entities))
	for name := range d.Entities {
		entityNames = append(entityNames, name)
	}
	sort.Strings(entityNames)

	for _, name := range entityNames {
		entity := d.Entities[name]
		sb.WriteString(fmt.Sprintf("    %s {\n", sanitizeName(entity.Name)))
		for _, attr := range entity.Attributes {
			sb.WriteString(formatMermaidAttribute(attr))
		}
		sb.WriteString("    }\n")
	}

	// Write relationships
	for _, rel := range d.Relationships {
		sb.WriteString(formatMermaidRelationship(rel))
	}

	return sb.String()
}

// formatMermaidAttribute formats an attribute for Mermaid syntax.
func formatMermaidAttribute(attr *Attribute) string {
	var parts []string

	// Type comes first in Mermaid
	parts = append(parts, sanitizeType(attr.Type))

	// Then the name
	parts = append(parts, attr.Name)

	// Add key constraint if present
	if attr.Key != nil {
		parts = append(parts, string(*attr.Key))
	}

	// Add comment for nullable or notes
	var comments []string
	if attr.Nullable {
		comments = append(comments, "nullable")
	}
	if attr.Note != nil {
		comments = append(comments, *attr.Note)
	}

	line := fmt.Sprintf("        %s", strings.Join(parts, " "))
	if len(comments) > 0 {
		line = fmt.Sprintf("%s %q", line, strings.Join(comments, ", "))
	}

	return line + "\n"
}

// formatMermaidRelationship formats a relationship for Mermaid syntax.
func formatMermaidRelationship(rel *Relationship) string {
	symbol := getMermaidCardinality(rel.Cardinality)
	label := rel.Field
	if rel.Label != nil {
		label = *rel.Label
	}

	return fmt.Sprintf("    %s %s %s : %s\n",
		sanitizeName(rel.From),
		symbol,
		sanitizeName(rel.To),
		label)
}

// getMermaidCardinality converts cardinality to Mermaid relationship syntax.
func getMermaidCardinality(c Cardinality) string {
	switch c {
	case OneToOne:
		return "||--||"
	case OneToMany:
		return "||--o{"
	case ManyToOne:
		return "}o--||"
	case ManyToMany:
		return "}o--o{"
	default:
		return "||--||"
	}
}

// sanitizeName ensures names are valid for Mermaid syntax.
func sanitizeName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	return name
}

// sanitizeType simplifies type names for display.
func sanitizeType(typeName string) string {
	// Remove package prefixes
	parts := strings.Split(typeName, ".")
	if len(parts) > 1 {
		typeName = parts[len(parts)-1]
	}

	// Simplify common patterns
	typeName = strings.ReplaceAll(typeName, "[]", "Array_")
	typeName = strings.ReplaceAll(typeName, "*", "")
	typeName = strings.ReplaceAll(typeName, " ", "_")

	return typeName
}
