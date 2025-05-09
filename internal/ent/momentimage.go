// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"api.us4ever/internal/ent/image"
	"api.us4ever/internal/ent/moment"
	"api.us4ever/internal/ent/momentimage"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// MomentImage is the model entity for the MomentImage schema.
type MomentImage struct {
	config `json:"-"`
	// ID of the ent.
	ID uint `json:"id,omitempty"`
	// ImageId holds the value of the "imageId" field.
	ImageId string `json:"imageId,omitempty"`
	// MomentId holds the value of the "momentId" field.
	MomentId string `json:"momentId,omitempty"`
	// Sort holds the value of the "sort" field.
	Sort int32 `json:"sort,omitempty"`
	// CreatedAt holds the value of the "createdAt" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updatedAt" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MomentImageQuery when eager-loading is set.
	Edges        MomentImageEdges `json:"edges"`
	selectValues sql.SelectValues
}

// MomentImageEdges holds the relations/edges for other nodes in the graph.
type MomentImageEdges struct {
	// Image holds the value of the image edge.
	Image *Image `json:"image,omitempty"`
	// Moment holds the value of the moment edge.
	Moment *Moment `json:"moment,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ImageOrErr returns the Image value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MomentImageEdges) ImageOrErr() (*Image, error) {
	if e.Image != nil {
		return e.Image, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: image.Label}
	}
	return nil, &NotLoadedError{edge: "image"}
}

// MomentOrErr returns the Moment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MomentImageEdges) MomentOrErr() (*Moment, error) {
	if e.Moment != nil {
		return e.Moment, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: moment.Label}
	}
	return nil, &NotLoadedError{edge: "moment"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*MomentImage) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case momentimage.FieldID, momentimage.FieldSort:
			values[i] = new(sql.NullInt64)
		case momentimage.FieldImageId, momentimage.FieldMomentId:
			values[i] = new(sql.NullString)
		case momentimage.FieldCreatedAt, momentimage.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the MomentImage fields.
func (mi *MomentImage) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case momentimage.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			mi.ID = uint(value.Int64)
		case momentimage.FieldImageId:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field imageId", values[i])
			} else if value.Valid {
				mi.ImageId = value.String
			}
		case momentimage.FieldMomentId:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field momentId", values[i])
			} else if value.Valid {
				mi.MomentId = value.String
			}
		case momentimage.FieldSort:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field sort", values[i])
			} else if value.Valid {
				mi.Sort = int32(value.Int64)
			}
		case momentimage.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createdAt", values[i])
			} else if value.Valid {
				mi.CreatedAt = value.Time
			}
		case momentimage.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updatedAt", values[i])
			} else if value.Valid {
				mi.UpdatedAt = value.Time
			}
		default:
			mi.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the MomentImage.
// This includes values selected through modifiers, order, etc.
func (mi *MomentImage) Value(name string) (ent.Value, error) {
	return mi.selectValues.Get(name)
}

// QueryImage queries the "image" edge of the MomentImage entity.
func (mi *MomentImage) QueryImage() *ImageQuery {
	return NewMomentImageClient(mi.config).QueryImage(mi)
}

// QueryMoment queries the "moment" edge of the MomentImage entity.
func (mi *MomentImage) QueryMoment() *MomentQuery {
	return NewMomentImageClient(mi.config).QueryMoment(mi)
}

// Update returns a builder for updating this MomentImage.
// Note that you need to call MomentImage.Unwrap() before calling this method if this MomentImage
// was returned from a transaction, and the transaction was committed or rolled back.
func (mi *MomentImage) Update() *MomentImageUpdateOne {
	return NewMomentImageClient(mi.config).UpdateOne(mi)
}

// Unwrap unwraps the MomentImage entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (mi *MomentImage) Unwrap() *MomentImage {
	_tx, ok := mi.config.driver.(*txDriver)
	if !ok {
		panic("ent: MomentImage is not a transactional entity")
	}
	mi.config.driver = _tx.drv
	return mi
}

// String implements the fmt.Stringer.
func (mi *MomentImage) String() string {
	var builder strings.Builder
	builder.WriteString("MomentImage(")
	builder.WriteString(fmt.Sprintf("id=%v, ", mi.ID))
	builder.WriteString("imageId=")
	builder.WriteString(mi.ImageId)
	builder.WriteString(", ")
	builder.WriteString("momentId=")
	builder.WriteString(mi.MomentId)
	builder.WriteString(", ")
	builder.WriteString("sort=")
	builder.WriteString(fmt.Sprintf("%v", mi.Sort))
	builder.WriteString(", ")
	builder.WriteString("createdAt=")
	builder.WriteString(mi.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updatedAt=")
	builder.WriteString(mi.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// MomentImages is a parsable slice of MomentImage.
type MomentImages []*MomentImage
