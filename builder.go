package erd

// NewDiagram creates a new ERD diagram.
func NewDiagram(title string) *Diagram {
	return &Diagram{
		Title:         title,
		Entities:      make(map[string]*Entity),
		Relationships: []*Relationship{},
	}
}

// WithDescription sets the description for the diagram.
func (d *Diagram) WithDescription(description string) *Diagram {
	d.Description = &description
	return d
}

// AddEntity adds an entity to the diagram.
func (d *Diagram) AddEntity(entity *Entity) *Diagram {
	d.Entities[entity.Name] = entity
	return d
}

// AddRelationship adds a relationship to the diagram.
func (d *Diagram) AddRelationship(rel *Relationship) *Diagram {
	d.Relationships = append(d.Relationships, rel)
	return d
}

// NewEntity creates a new entity.
func NewEntity(name string) *Entity {
	return &Entity{
		Name:       name,
		Attributes: []*Attribute{},
	}
}

// WithPackage sets the package name for the entity.
func (e *Entity) WithPackage(pkg string) *Entity {
	e.Package = &pkg
	return e
}

// WithNote adds a note to the entity.
func (e *Entity) WithNote(note string) *Entity {
	e.Note = &note
	return e
}

// AddAttribute adds an attribute to the entity.
func (e *Entity) AddAttribute(attr *Attribute) *Entity {
	e.Attributes = append(e.Attributes, attr)
	return e
}

// NewAttribute creates a new attribute.
func NewAttribute(name, attrType string) *Attribute {
	return &Attribute{
		Name:     name,
		Type:     attrType,
		Nullable: false,
	}
}

// WithKey marks the attribute with a key type.
func (a *Attribute) WithKey(keyType KeyType) *Attribute {
	a.Key = &keyType
	return a
}

// WithPrimaryKey marks the attribute as a primary key.
func (a *Attribute) WithPrimaryKey() *Attribute {
	key := PrimaryKey
	a.Key = &key
	return a
}

// WithForeignKey marks the attribute as a foreign key.
func (a *Attribute) WithForeignKey() *Attribute {
	key := ForeignKey
	a.Key = &key
	return a
}

// WithUnique marks the attribute as unique.
func (a *Attribute) WithUnique() *Attribute {
	key := UniqueKey
	a.Key = &key
	return a
}

// WithNullable marks the attribute as nullable.
func (a *Attribute) WithNullable() *Attribute {
	a.Nullable = true
	return a
}

// WithNote adds a note to the attribute.
func (a *Attribute) WithNote(note string) *Attribute {
	a.Note = &note
	return a
}

// NewRelationship creates a new relationship.
func NewRelationship(from, to, field string, cardinality Cardinality) *Relationship {
	return &Relationship{
		From:        from,
		To:          to,
		Field:       field,
		Cardinality: cardinality,
	}
}

// WithLabel sets a label for the relationship.
func (r *Relationship) WithLabel(label string) *Relationship {
	r.Label = &label
	return r
}

// WithNote adds a note to the relationship.
func (r *Relationship) WithNote(note string) *Relationship {
	r.Note = &note
	return r
}
