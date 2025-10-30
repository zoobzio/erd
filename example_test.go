package erd_test

import (
	"fmt"

	"github.com/zoobzio/erd"
)

func ExampleDiagram_ToMermaid() {
	// Create a simple domain model diagram
	diagram := erd.NewDiagram("User Management System").
		WithDescription("Core entities for user and order management")

	// Define User entity
	user := erd.NewEntity("User").
		WithPackage("domain").
		AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
		AddAttribute(erd.NewAttribute("Email", "string").WithUnique()).
		AddAttribute(erd.NewAttribute("Name", "string")).
		AddAttribute(erd.NewAttribute("ProfileID", "string").WithForeignKey().WithNullable())

	// Define Profile entity
	profile := erd.NewEntity("Profile").
		WithPackage("domain").
		AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
		AddAttribute(erd.NewAttribute("Bio", "string").WithNullable()).
		AddAttribute(erd.NewAttribute("Avatar", "string").WithNullable())

	// Define Order entity
	order := erd.NewEntity("Order").
		WithPackage("domain").
		AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
		AddAttribute(erd.NewAttribute("UserID", "string").WithForeignKey()).
		AddAttribute(erd.NewAttribute("Total", "float64")).
		AddAttribute(erd.NewAttribute("Status", "string"))

	// Add entities to diagram
	diagram.
		AddEntity(user).
		AddEntity(profile).
		AddEntity(order)

	// Define relationships
	diagram.
		AddRelationship(erd.NewRelationship("User", "Profile", "Profile", erd.OneToOne)).
		AddRelationship(erd.NewRelationship("User", "Order", "Orders", erd.OneToMany))

	// Generate Mermaid diagram
	fmt.Println(diagram.ToMermaid())

	// Output:
	// erDiagram
	//     Order {
	//         string ID PK
	//         string UserID FK
	//         float64 Total
	//         string Status
	//     }
	//     Profile {
	//         string ID PK
	//         string Bio "nullable"
	//         string Avatar "nullable"
	//     }
	//     User {
	//         string ID PK
	//         string Email UK
	//         string Name
	//         string ProfileID FK "nullable"
	//     }
	//     User ||--|| Profile : Profile
	//     User ||--o{ Order : Orders
}

func ExampleDiagram_ToDOT() {
	// Create a simple e-commerce domain model
	diagram := erd.NewDiagram("E-commerce Domain")

	// Define entities
	product := erd.NewEntity("Product").
		AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
		AddAttribute(erd.NewAttribute("Name", "string")).
		AddAttribute(erd.NewAttribute("Price", "float64"))

	cart := erd.NewEntity("Cart").
		AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
		AddAttribute(erd.NewAttribute("UserID", "string").WithForeignKey())

	cartItem := erd.NewEntity("CartItem").
		AddAttribute(erd.NewAttribute("ID", "string").WithPrimaryKey()).
		AddAttribute(erd.NewAttribute("CartID", "string").WithForeignKey()).
		AddAttribute(erd.NewAttribute("ProductID", "string").WithForeignKey()).
		AddAttribute(erd.NewAttribute("Quantity", "int"))

	diagram.
		AddEntity(product).
		AddEntity(cart).
		AddEntity(cartItem)

	// Define relationships
	diagram.
		AddRelationship(erd.NewRelationship("Cart", "CartItem", "Items", erd.OneToMany)).
		AddRelationship(erd.NewRelationship("CartItem", "Product", "Product", erd.ManyToOne))

	// Generate DOT diagram
	fmt.Println(diagram.ToDOT())

	// Output:
	// digraph ERD {
	//     rankdir=LR;
	//     node [shape=record];
	//     labelloc="t";
	//     label="E-commerce Domain";
	//
	//     Cart [label="{Cart|PK ID: string\lFK UserID: string\l}"];
	//     CartItem [label="{CartItem|PK ID: string\lFK CartID: string\lFK ProductID: string\lQuantity: int\l}"];
	//     Product [label="{Product|PK ID: string\lName: string\lPrice: float64\l}"];
	//
	//     Cart -> CartItem [arrowhead=crow, arrowtail=normal, dir=both label="Items"];
	//     CartItem -> Product [arrowhead=normal, arrowtail=crow, dir=both label="Product"];
	// }
}

func ExampleDiagram_Validate() {
	// Create a diagram with validation errors
	diagram := erd.NewDiagram("Invalid Diagram")

	// Entity with no attributes (invalid)
	empty := erd.NewEntity("Empty")
	diagram.AddEntity(empty)

	// Relationship referencing non-existent entity (invalid)
	diagram.AddRelationship(
		erd.NewRelationship("Empty", "NonExistent", "Field", erd.OneToOne),
	)

	// Validate and print errors
	errors := diagram.Validate()
	for _, err := range errors {
		fmt.Println(err.Error())
	}

	// Output:
	// Entity[Empty].Attributes: entity must have at least one attribute
	// Relationship[0].To: entity 'NonExistent' does not exist
}
