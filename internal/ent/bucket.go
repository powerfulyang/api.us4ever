// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"api.us4ever/internal/ent/bucket"
	"api.us4ever/internal/ent/user"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Bucket is the model entity for the Bucket schema.
type Bucket struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// BucketName holds the value of the "bucketName" field.
	BucketName string `json:"bucketName,omitempty"`
	// Provider holds the value of the "provider" field.
	Provider bucket.Provider `json:"provider,omitempty"`
	// Region holds the value of the "region" field.
	Region string `json:"region,omitempty"`
	// Endpoint holds the value of the "endpoint" field.
	Endpoint string `json:"endpoint,omitempty"`
	// PublicUrl holds the value of the "publicUrl" field.
	PublicUrl string `json:"publicUrl,omitempty"`
	// AccessKey holds the value of the "accessKey" field.
	AccessKey string `json:"accessKey,omitempty"`
	// SecretKey holds the value of the "secretKey" field.
	SecretKey string `json:"secretKey,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// CreatedAt holds the value of the "createdAt" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updatedAt" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// OwnerId holds the value of the "ownerId" field.
	OwnerId string `json:"ownerId,omitempty"`
	// ExtraData holds the value of the "extraData" field.
	ExtraData json.RawMessage `json:"extraData,omitempty"`
	// Category holds the value of the "category" field.
	Category string `json:"category,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BucketQuery when eager-loading is set.
	Edges        BucketEdges `json:"edges"`
	selectValues sql.SelectValues
}

// BucketEdges holds the relations/edges for other nodes in the graph.
type BucketEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// Files holds the value of the files edge.
	Files []*File `json:"files,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BucketEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// FilesOrErr returns the Files value or an error if the edge
// was not loaded in eager-loading.
func (e BucketEdges) FilesOrErr() ([]*File, error) {
	if e.loadedTypes[1] {
		return e.Files, nil
	}
	return nil, &NotLoadedError{edge: "files"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Bucket) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case bucket.FieldExtraData:
			values[i] = new([]byte)
		case bucket.FieldID, bucket.FieldName, bucket.FieldBucketName, bucket.FieldProvider, bucket.FieldRegion, bucket.FieldEndpoint, bucket.FieldPublicUrl, bucket.FieldAccessKey, bucket.FieldSecretKey, bucket.FieldDescription, bucket.FieldOwnerId, bucket.FieldCategory:
			values[i] = new(sql.NullString)
		case bucket.FieldCreatedAt, bucket.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Bucket fields.
func (b *Bucket) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case bucket.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				b.ID = value.String
			}
		case bucket.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				b.Name = value.String
			}
		case bucket.FieldBucketName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field bucketName", values[i])
			} else if value.Valid {
				b.BucketName = value.String
			}
		case bucket.FieldProvider:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field provider", values[i])
			} else if value.Valid {
				b.Provider = bucket.Provider(value.String)
			}
		case bucket.FieldRegion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field region", values[i])
			} else if value.Valid {
				b.Region = value.String
			}
		case bucket.FieldEndpoint:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field endpoint", values[i])
			} else if value.Valid {
				b.Endpoint = value.String
			}
		case bucket.FieldPublicUrl:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field publicUrl", values[i])
			} else if value.Valid {
				b.PublicUrl = value.String
			}
		case bucket.FieldAccessKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field accessKey", values[i])
			} else if value.Valid {
				b.AccessKey = value.String
			}
		case bucket.FieldSecretKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field secretKey", values[i])
			} else if value.Valid {
				b.SecretKey = value.String
			}
		case bucket.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				b.Description = value.String
			}
		case bucket.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createdAt", values[i])
			} else if value.Valid {
				b.CreatedAt = value.Time
			}
		case bucket.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updatedAt", values[i])
			} else if value.Valid {
				b.UpdatedAt = value.Time
			}
		case bucket.FieldOwnerId:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ownerId", values[i])
			} else if value.Valid {
				b.OwnerId = value.String
			}
		case bucket.FieldExtraData:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field extraData", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &b.ExtraData); err != nil {
					return fmt.Errorf("unmarshal field extraData: %w", err)
				}
			}
		case bucket.FieldCategory:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field category", values[i])
			} else if value.Valid {
				b.Category = value.String
			}
		default:
			b.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Bucket.
// This includes values selected through modifiers, order, etc.
func (b *Bucket) Value(name string) (ent.Value, error) {
	return b.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the Bucket entity.
func (b *Bucket) QueryUser() *UserQuery {
	return NewBucketClient(b.config).QueryUser(b)
}

// QueryFiles queries the "files" edge of the Bucket entity.
func (b *Bucket) QueryFiles() *FileQuery {
	return NewBucketClient(b.config).QueryFiles(b)
}

// Update returns a builder for updating this Bucket.
// Note that you need to call Bucket.Unwrap() before calling this method if this Bucket
// was returned from a transaction, and the transaction was committed or rolled back.
func (b *Bucket) Update() *BucketUpdateOne {
	return NewBucketClient(b.config).UpdateOne(b)
}

// Unwrap unwraps the Bucket entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (b *Bucket) Unwrap() *Bucket {
	_tx, ok := b.config.driver.(*txDriver)
	if !ok {
		panic("ent: Bucket is not a transactional entity")
	}
	b.config.driver = _tx.drv
	return b
}

// String implements the fmt.Stringer.
func (b *Bucket) String() string {
	var builder strings.Builder
	builder.WriteString("Bucket(")
	builder.WriteString(fmt.Sprintf("id=%v, ", b.ID))
	builder.WriteString("name=")
	builder.WriteString(b.Name)
	builder.WriteString(", ")
	builder.WriteString("bucketName=")
	builder.WriteString(b.BucketName)
	builder.WriteString(", ")
	builder.WriteString("provider=")
	builder.WriteString(fmt.Sprintf("%v", b.Provider))
	builder.WriteString(", ")
	builder.WriteString("region=")
	builder.WriteString(b.Region)
	builder.WriteString(", ")
	builder.WriteString("endpoint=")
	builder.WriteString(b.Endpoint)
	builder.WriteString(", ")
	builder.WriteString("publicUrl=")
	builder.WriteString(b.PublicUrl)
	builder.WriteString(", ")
	builder.WriteString("accessKey=")
	builder.WriteString(b.AccessKey)
	builder.WriteString(", ")
	builder.WriteString("secretKey=")
	builder.WriteString(b.SecretKey)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(b.Description)
	builder.WriteString(", ")
	builder.WriteString("createdAt=")
	builder.WriteString(b.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updatedAt=")
	builder.WriteString(b.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("ownerId=")
	builder.WriteString(b.OwnerId)
	builder.WriteString(", ")
	builder.WriteString("extraData=")
	builder.WriteString(fmt.Sprintf("%v", b.ExtraData))
	builder.WriteString(", ")
	builder.WriteString("category=")
	builder.WriteString(b.Category)
	builder.WriteByte(')')
	return builder.String()
}

// Buckets is a parsable slice of Bucket.
type Buckets []*Bucket
