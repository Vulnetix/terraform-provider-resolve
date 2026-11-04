package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// kindSpec describes one config kind as a Terraform resource. cmd/genresources
// writes one per kind from the config-schemas JSON Schemas (zz_generated_*.go);
// nothing in it is written by hand.
type kindSpec struct {
	// Kind is the document kind, such as QualityGate.
	Kind string
	// TypeName is the resource type without the provider prefix.
	TypeName string
	// Description is the kind schema's description.
	Description string
	// Singleton kinds have exactly one object per tenant, named default.
	Singleton bool
	// NamePattern is the metadata.name pattern for non-singleton kinds.
	NamePattern string
	// DependsOn are the kinds this kind can reference (x-vulnetix.dependsOn).
	DependsOn []string
	// OwnerOnly kinds can only be applied with an owner-role token.
	OwnerOnly bool
	// Attributes returns the spec attributes. They sit at the top level of
	// the resource beside name, labels and annotations.
	Attributes func() map[string]schema.Attribute
	// Fields maps the same attributes onto document spec keys.
	Fields []field
}

// fieldKind is how a value is carried between Terraform and a document.
type fieldKind int

const (
	fieldString fieldKind = iota
	fieldInt
	fieldFloat
	fieldBool
	fieldObject     // SingleNestedAttribute; Children
	fieldList       // ListAttribute of Elem
	fieldMap        // MapAttribute of Elem
	fieldListObject // ListNestedAttribute; Children
	fieldMapObject  // MapNestedAttribute; Children
	fieldSet        // SetAttribute of Elem, for arrays whose order the server ignores
	fieldSetObject  // SetNestedAttribute; Children
	fieldJSON       // a JSON-encoded string for shapes HCL cannot type
)

// field maps one Terraform attribute onto one spec key.
type field struct {
	Attr     string
	JSON     string
	Kind     fieldKind
	Elem     fieldKind
	Children []field
	// WriteOnly marks a secret's <attr>_wo argument. Its value is sent in the
	// document under JSON, never kept in state, and only sent when Version
	// names a sibling whose value changed (or on create).
	WriteOnly bool
	Version   string
	// Local attributes exist only in Terraform, such as <attr>_wo_version.
	Local bool
	// Default is a scalar's schema default, JSON-encoded ("\"full\"", "true",
	// "86400000"), or empty. An import reads a value equal to it as unset.
	Default string
}
