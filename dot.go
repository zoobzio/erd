package erd

import (
	"fmt"
	"sort"
	"strings"
)

// ToDOT generates a GraphViz DOT diagram from the diagram structure.
func (d *Diagram) ToDOT() string {
	var sb strings.Builder

	sb.WriteString("digraph ERD {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape=record];\n")

	// Add title if present
	if d.Title != "" {
		sb.WriteString("    labelloc=\"t\";\n")
		sb.WriteString(fmt.Sprintf("    label=%q;\n", escapeDOT(d.Title)))
	}

	sb.WriteString("\n")

	// Write entities with their attributes (sorted by name for deterministic output)
	entityNames := make([]string, 0, len(d.Entities))
	for name := range d.Entities {
		entityNames = append(entityNames, name)
	}
	sort.Strings(entityNames)

	for _, name := range entityNames {
		entity := d.Entities[name]
		sb.WriteString(formatDOTEntity(entity))
	}

	sb.WriteString("\n")

	// Write relationships
	for _, rel := range d.Relationships {
		sb.WriteString(formatDOTRelationship(rel))
	}

	sb.WriteString("}\n")
	return sb.String()
}

// formatDOTEntity formats an entity as a DOT record node.
func formatDOTEntity(entity *Entity) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("    %s [label=\"{%s|",
		sanitizeName(entity.Name),
		escapeDOT(entity.Name)))

	attrs := make([]string, 0, len(entity.Attributes))
	for _, attr := range entity.Attributes {
		attrs = append(attrs, formatDOTAttribute(attr))
	}

	sb.WriteString(strings.Join(attrs, "\\l"))
	sb.WriteString("\\l}\"];\n")

	return sb.String()
}

// formatDOTAttribute formats an attribute for DOT syntax.
func formatDOTAttribute(attr *Attribute) string {
	var parts []string

	// Add key indicator
	if attr.Key != nil {
		parts = append(parts, string(*attr.Key))
	}

	// Add name and type
	parts = append(parts, fmt.Sprintf("%s: %s",
		attr.Name,
		sanitizeType(attr.Type)))

	// Add nullable indicator
	if attr.Nullable {
		parts = append(parts, "?")
	}

	return escapeDOT(strings.Join(parts, " "))
}

// formatDOTRelationship formats a relationship for DOT syntax.
func formatDOTRelationship(rel *Relationship) string {
	edgeStyle := getDOTCardinality(rel.Cardinality)
	label := rel.Field
	if rel.Label != nil {
		label = *rel.Label
	}

	return fmt.Sprintf("    %s -> %s [%s label=%q];\n",
		sanitizeName(rel.From),
		sanitizeName(rel.To),
		edgeStyle,
		escapeDOT(label))
}

// getDOTCardinality returns DOT edge styling for cardinality.
func getDOTCardinality(c Cardinality) string {
	switch c {
	case OneToOne:
		return "arrowhead=normal, arrowtail=normal, dir=both"
	case OneToMany:
		return "arrowhead=crow, arrowtail=normal, dir=both"
	case ManyToOne:
		return "arrowhead=normal, arrowtail=crow, dir=both"
	case ManyToMany:
		return "arrowhead=crow, arrowtail=crow, dir=both"
	default:
		return "arrowhead=normal"
	}
}

// escapeDOT escapes special characters for DOT syntax.
func escapeDOT(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
