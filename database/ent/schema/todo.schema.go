package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("title").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.Bool("completed").
			Default(false),
		field.Enum("priority").
			Values("low", "medium", "high").
			Optional(),
		field.Time("due_date").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Optional().
			UpdateDefault(time.Now),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return nil
}
