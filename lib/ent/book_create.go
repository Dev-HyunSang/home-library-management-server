// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/dev-hyunsang/home-library/lib/ent/book"
	"github.com/dev-hyunsang/home-library/lib/ent/user"
	"github.com/google/uuid"
)

// BookCreate is the builder for creating a Book entity.
type BookCreate struct {
	config
	mutation *BookMutation
	hooks    []Hook
}

// SetBookTitle sets the "book_title" field.
func (bc *BookCreate) SetBookTitle(s string) *BookCreate {
	bc.mutation.SetBookTitle(s)
	return bc
}

// SetAuthor sets the "author" field.
func (bc *BookCreate) SetAuthor(s string) *BookCreate {
	bc.mutation.SetAuthor(s)
	return bc
}

// SetBookIsbn sets the "book_isbn" field.
func (bc *BookCreate) SetBookIsbn(i int) *BookCreate {
	bc.mutation.SetBookIsbn(i)
	return bc
}

// SetNillableBookIsbn sets the "book_isbn" field if the given value is not nil.
func (bc *BookCreate) SetNillableBookIsbn(i *int) *BookCreate {
	if i != nil {
		bc.SetBookIsbn(*i)
	}
	return bc
}

// SetRegisteredAt sets the "registered_at" field.
func (bc *BookCreate) SetRegisteredAt(t time.Time) *BookCreate {
	bc.mutation.SetRegisteredAt(t)
	return bc
}

// SetNillableRegisteredAt sets the "registered_at" field if the given value is not nil.
func (bc *BookCreate) SetNillableRegisteredAt(t *time.Time) *BookCreate {
	if t != nil {
		bc.SetRegisteredAt(*t)
	}
	return bc
}

// SetComplatedAt sets the "complated_at" field.
func (bc *BookCreate) SetComplatedAt(t time.Time) *BookCreate {
	bc.mutation.SetComplatedAt(t)
	return bc
}

// SetNillableComplatedAt sets the "complated_at" field if the given value is not nil.
func (bc *BookCreate) SetNillableComplatedAt(t *time.Time) *BookCreate {
	if t != nil {
		bc.SetComplatedAt(*t)
	}
	return bc
}

// SetID sets the "id" field.
func (bc *BookCreate) SetID(u uuid.UUID) *BookCreate {
	bc.mutation.SetID(u)
	return bc
}

// SetOwnerID sets the "owner" edge to the User entity by ID.
func (bc *BookCreate) SetOwnerID(id uuid.UUID) *BookCreate {
	bc.mutation.SetOwnerID(id)
	return bc
}

// SetOwner sets the "owner" edge to the User entity.
func (bc *BookCreate) SetOwner(u *User) *BookCreate {
	return bc.SetOwnerID(u.ID)
}

// Mutation returns the BookMutation object of the builder.
func (bc *BookCreate) Mutation() *BookMutation {
	return bc.mutation
}

// Save creates the Book in the database.
func (bc *BookCreate) Save(ctx context.Context) (*Book, error) {
	bc.defaults()
	return withHooks(ctx, bc.sqlSave, bc.mutation, bc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (bc *BookCreate) SaveX(ctx context.Context) *Book {
	v, err := bc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bc *BookCreate) Exec(ctx context.Context) error {
	_, err := bc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bc *BookCreate) ExecX(ctx context.Context) {
	if err := bc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bc *BookCreate) defaults() {
	if _, ok := bc.mutation.RegisteredAt(); !ok {
		v := book.DefaultRegisteredAt
		bc.mutation.SetRegisteredAt(v)
	}
	if _, ok := bc.mutation.ComplatedAt(); !ok {
		v := book.DefaultComplatedAt
		bc.mutation.SetComplatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bc *BookCreate) check() error {
	if _, ok := bc.mutation.BookTitle(); !ok {
		return &ValidationError{Name: "book_title", err: errors.New(`ent: missing required field "Book.book_title"`)}
	}
	if v, ok := bc.mutation.BookTitle(); ok {
		if err := book.BookTitleValidator(v); err != nil {
			return &ValidationError{Name: "book_title", err: fmt.Errorf(`ent: validator failed for field "Book.book_title": %w`, err)}
		}
	}
	if _, ok := bc.mutation.Author(); !ok {
		return &ValidationError{Name: "author", err: errors.New(`ent: missing required field "Book.author"`)}
	}
	if v, ok := bc.mutation.Author(); ok {
		if err := book.AuthorValidator(v); err != nil {
			return &ValidationError{Name: "author", err: fmt.Errorf(`ent: validator failed for field "Book.author": %w`, err)}
		}
	}
	if _, ok := bc.mutation.RegisteredAt(); !ok {
		return &ValidationError{Name: "registered_at", err: errors.New(`ent: missing required field "Book.registered_at"`)}
	}
	if _, ok := bc.mutation.ComplatedAt(); !ok {
		return &ValidationError{Name: "complated_at", err: errors.New(`ent: missing required field "Book.complated_at"`)}
	}
	if len(bc.mutation.OwnerIDs()) == 0 {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Book.owner"`)}
	}
	return nil
}

func (bc *BookCreate) sqlSave(ctx context.Context) (*Book, error) {
	if err := bc.check(); err != nil {
		return nil, err
	}
	_node, _spec := bc.createSpec()
	if err := sqlgraph.CreateNode(ctx, bc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	bc.mutation.id = &_node.ID
	bc.mutation.done = true
	return _node, nil
}

func (bc *BookCreate) createSpec() (*Book, *sqlgraph.CreateSpec) {
	var (
		_node = &Book{config: bc.config}
		_spec = sqlgraph.NewCreateSpec(book.Table, sqlgraph.NewFieldSpec(book.FieldID, field.TypeUUID))
	)
	if id, ok := bc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := bc.mutation.BookTitle(); ok {
		_spec.SetField(book.FieldBookTitle, field.TypeString, value)
		_node.BookTitle = value
	}
	if value, ok := bc.mutation.Author(); ok {
		_spec.SetField(book.FieldAuthor, field.TypeString, value)
		_node.Author = value
	}
	if value, ok := bc.mutation.BookIsbn(); ok {
		_spec.SetField(book.FieldBookIsbn, field.TypeInt, value)
		_node.BookIsbn = value
	}
	if value, ok := bc.mutation.RegisteredAt(); ok {
		_spec.SetField(book.FieldRegisteredAt, field.TypeTime, value)
		_node.RegisteredAt = value
	}
	if value, ok := bc.mutation.ComplatedAt(); ok {
		_spec.SetField(book.FieldComplatedAt, field.TypeTime, value)
		_node.ComplatedAt = value
	}
	if nodes := bc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   book.OwnerTable,
			Columns: []string{book.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_books = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// BookCreateBulk is the builder for creating many Book entities in bulk.
type BookCreateBulk struct {
	config
	err      error
	builders []*BookCreate
}

// Save creates the Book entities in the database.
func (bcb *BookCreateBulk) Save(ctx context.Context) ([]*Book, error) {
	if bcb.err != nil {
		return nil, bcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(bcb.builders))
	nodes := make([]*Book, len(bcb.builders))
	mutators := make([]Mutator, len(bcb.builders))
	for i := range bcb.builders {
		func(i int, root context.Context) {
			builder := bcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*BookMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, bcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, bcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, bcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (bcb *BookCreateBulk) SaveX(ctx context.Context) []*Book {
	v, err := bcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bcb *BookCreateBulk) Exec(ctx context.Context) error {
	_, err := bcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bcb *BookCreateBulk) ExecX(ctx context.Context) {
	if err := bcb.Exec(ctx); err != nil {
		panic(err)
	}
}
