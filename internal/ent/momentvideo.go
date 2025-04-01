// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"api.us4ever/internal/ent/moment"
	"api.us4ever/internal/ent/momentvideo"
	"api.us4ever/internal/ent/video"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// MomentVideo is the model entity for the MomentVideo schema.
type MomentVideo struct {
	config `json:"-"`
	// ID of the ent.
	ID uint `json:"id,omitempty"`
	// VideoId holds the value of the "videoId" field.
	VideoId string `json:"videoId,omitempty"`
	// MomentId holds the value of the "momentId" field.
	MomentId string `json:"momentId,omitempty"`
	// Sort holds the value of the "sort" field.
	Sort int32 `json:"sort,omitempty"`
	// CreatedAt holds the value of the "createdAt" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updatedAt" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MomentVideoQuery when eager-loading is set.
	Edges        MomentVideoEdges `json:"edges"`
	selectValues sql.SelectValues
}

// MomentVideoEdges holds the relations/edges for other nodes in the graph.
type MomentVideoEdges struct {
	// Moment holds the value of the moment edge.
	Moment *Moment `json:"moment,omitempty"`
	// Video holds the value of the video edge.
	Video *Video `json:"video,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// MomentOrErr returns the Moment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MomentVideoEdges) MomentOrErr() (*Moment, error) {
	if e.Moment != nil {
		return e.Moment, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: moment.Label}
	}
	return nil, &NotLoadedError{edge: "moment"}
}

// VideoOrErr returns the Video value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MomentVideoEdges) VideoOrErr() (*Video, error) {
	if e.Video != nil {
		return e.Video, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: video.Label}
	}
	return nil, &NotLoadedError{edge: "video"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*MomentVideo) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case momentvideo.FieldID, momentvideo.FieldSort:
			values[i] = new(sql.NullInt64)
		case momentvideo.FieldVideoId, momentvideo.FieldMomentId:
			values[i] = new(sql.NullString)
		case momentvideo.FieldCreatedAt, momentvideo.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the MomentVideo fields.
func (mv *MomentVideo) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case momentvideo.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			mv.ID = uint(value.Int64)
		case momentvideo.FieldVideoId:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field videoId", values[i])
			} else if value.Valid {
				mv.VideoId = value.String
			}
		case momentvideo.FieldMomentId:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field momentId", values[i])
			} else if value.Valid {
				mv.MomentId = value.String
			}
		case momentvideo.FieldSort:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field sort", values[i])
			} else if value.Valid {
				mv.Sort = int32(value.Int64)
			}
		case momentvideo.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createdAt", values[i])
			} else if value.Valid {
				mv.CreatedAt = value.Time
			}
		case momentvideo.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updatedAt", values[i])
			} else if value.Valid {
				mv.UpdatedAt = value.Time
			}
		default:
			mv.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the MomentVideo.
// This includes values selected through modifiers, order, etc.
func (mv *MomentVideo) Value(name string) (ent.Value, error) {
	return mv.selectValues.Get(name)
}

// QueryMoment queries the "moment" edge of the MomentVideo entity.
func (mv *MomentVideo) QueryMoment() *MomentQuery {
	return NewMomentVideoClient(mv.config).QueryMoment(mv)
}

// QueryVideo queries the "video" edge of the MomentVideo entity.
func (mv *MomentVideo) QueryVideo() *VideoQuery {
	return NewMomentVideoClient(mv.config).QueryVideo(mv)
}

// Update returns a builder for updating this MomentVideo.
// Note that you need to call MomentVideo.Unwrap() before calling this method if this MomentVideo
// was returned from a transaction, and the transaction was committed or rolled back.
func (mv *MomentVideo) Update() *MomentVideoUpdateOne {
	return NewMomentVideoClient(mv.config).UpdateOne(mv)
}

// Unwrap unwraps the MomentVideo entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (mv *MomentVideo) Unwrap() *MomentVideo {
	_tx, ok := mv.config.driver.(*txDriver)
	if !ok {
		panic("ent: MomentVideo is not a transactional entity")
	}
	mv.config.driver = _tx.drv
	return mv
}

// String implements the fmt.Stringer.
func (mv *MomentVideo) String() string {
	var builder strings.Builder
	builder.WriteString("MomentVideo(")
	builder.WriteString(fmt.Sprintf("id=%v, ", mv.ID))
	builder.WriteString("videoId=")
	builder.WriteString(mv.VideoId)
	builder.WriteString(", ")
	builder.WriteString("momentId=")
	builder.WriteString(mv.MomentId)
	builder.WriteString(", ")
	builder.WriteString("sort=")
	builder.WriteString(fmt.Sprintf("%v", mv.Sort))
	builder.WriteString(", ")
	builder.WriteString("createdAt=")
	builder.WriteString(mv.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updatedAt=")
	builder.WriteString(mv.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// MomentVideos is a parsable slice of MomentVideo.
type MomentVideos []*MomentVideo
