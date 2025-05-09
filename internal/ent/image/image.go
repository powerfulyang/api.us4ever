// Code generated by ent, DO NOT EDIT.

package image

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the image type in the database.
	Label = "image"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldSize holds the string denoting the size field in the database.
	FieldSize = "size"
	// FieldWidth holds the string denoting the width field in the database.
	FieldWidth = "width"
	// FieldHeight holds the string denoting the height field in the database.
	FieldHeight = "height"
	// FieldExif holds the string denoting the exif field in the database.
	FieldExif = "exif"
	// FieldHash holds the string denoting the hash field in the database.
	FieldHash = "hash"
	// FieldAddress holds the string denoting the address field in the database.
	FieldAddress = "address"
	// FieldIsPublic holds the string denoting the ispublic field in the database.
	FieldIsPublic = "isPublic"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldTags holds the string denoting the tags field in the database.
	FieldTags = "tags"
	// FieldThumbnail10x holds the string denoting the thumbnail_10x field in the database.
	FieldThumbnail10x = "thumbnail_10x"
	// FieldThumbnail320xID holds the string denoting the thumbnail_320x_id field in the database.
	FieldThumbnail320xID = "thumbnail_320x_id"
	// FieldThumbnail768xID holds the string denoting the thumbnail_768x_id field in the database.
	FieldThumbnail768xID = "thumbnail_768x_id"
	// FieldCompressedID holds the string denoting the compressed_id field in the database.
	FieldCompressedID = "compressed_id"
	// FieldOriginalID holds the string denoting the original_id field in the database.
	FieldOriginalID = "original_id"
	// FieldCreatedAt holds the string denoting the createdat field in the database.
	FieldCreatedAt = "createdAt"
	// FieldUpdatedAt holds the string denoting the updatedat field in the database.
	FieldUpdatedAt = "updatedAt"
	// FieldUploadedBy holds the string denoting the uploadedby field in the database.
	FieldUploadedBy = "uploadedBy"
	// FieldCategory holds the string denoting the category field in the database.
	FieldCategory = "category"
	// FieldExtraData holds the string denoting the extradata field in the database.
	FieldExtraData = "extraData"
	// FieldDescriptionVector holds the string denoting the description_vector field in the database.
	FieldDescriptionVector = "description_vector"
	// EdgeCompressed holds the string denoting the compressed edge name in mutations.
	EdgeCompressed = "compressed"
	// EdgeOriginal holds the string denoting the original edge name in mutations.
	EdgeOriginal = "original"
	// EdgeThumbnail320x holds the string denoting the thumbnail320x edge name in mutations.
	EdgeThumbnail320x = "thumbnail320x"
	// EdgeThumbnail768x holds the string denoting the thumbnail768x edge name in mutations.
	EdgeThumbnail768x = "thumbnail768x"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeMomentImages holds the string denoting the moment_images edge name in mutations.
	EdgeMomentImages = "moment_images"
	// Table holds the table name of the image in the database.
	Table = "images"
	// CompressedTable is the table that holds the compressed relation/edge.
	CompressedTable = "images"
	// CompressedInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	CompressedInverseTable = "files"
	// CompressedColumn is the table column denoting the compressed relation/edge.
	CompressedColumn = "compressed_id"
	// OriginalTable is the table that holds the original relation/edge.
	OriginalTable = "images"
	// OriginalInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	OriginalInverseTable = "files"
	// OriginalColumn is the table column denoting the original relation/edge.
	OriginalColumn = "original_id"
	// Thumbnail320xTable is the table that holds the thumbnail320x relation/edge.
	Thumbnail320xTable = "images"
	// Thumbnail320xInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	Thumbnail320xInverseTable = "files"
	// Thumbnail320xColumn is the table column denoting the thumbnail320x relation/edge.
	Thumbnail320xColumn = "thumbnail_320x_id"
	// Thumbnail768xTable is the table that holds the thumbnail768x relation/edge.
	Thumbnail768xTable = "images"
	// Thumbnail768xInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	Thumbnail768xInverseTable = "files"
	// Thumbnail768xColumn is the table column denoting the thumbnail768x relation/edge.
	Thumbnail768xColumn = "thumbnail_768x_id"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "images"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "uploadedBy"
	// MomentImagesTable is the table that holds the moment_images relation/edge.
	MomentImagesTable = "moment_images"
	// MomentImagesInverseTable is the table name for the MomentImage entity.
	// It exists in this package in order to avoid circular dependency with the "momentimage" package.
	MomentImagesInverseTable = "moment_images"
	// MomentImagesColumn is the table column denoting the moment_images relation/edge.
	MomentImagesColumn = "imageId"
)

// Columns holds all SQL columns for image fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldType,
	FieldSize,
	FieldWidth,
	FieldHeight,
	FieldExif,
	FieldHash,
	FieldAddress,
	FieldIsPublic,
	FieldDescription,
	FieldTags,
	FieldThumbnail10x,
	FieldThumbnail320xID,
	FieldThumbnail768xID,
	FieldCompressedID,
	FieldOriginalID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldUploadedBy,
	FieldCategory,
	FieldExtraData,
	FieldDescriptionVector,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Image queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// BySize orders the results by the size field.
func BySize(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSize, opts...).ToFunc()
}

// ByWidth orders the results by the width field.
func ByWidth(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWidth, opts...).ToFunc()
}

// ByHeight orders the results by the height field.
func ByHeight(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHeight, opts...).ToFunc()
}

// ByHash orders the results by the hash field.
func ByHash(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHash, opts...).ToFunc()
}

// ByAddress orders the results by the address field.
func ByAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAddress, opts...).ToFunc()
}

// ByIsPublic orders the results by the isPublic field.
func ByIsPublic(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsPublic, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByThumbnail320xID orders the results by the thumbnail_320x_id field.
func ByThumbnail320xID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldThumbnail320xID, opts...).ToFunc()
}

// ByThumbnail768xID orders the results by the thumbnail_768x_id field.
func ByThumbnail768xID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldThumbnail768xID, opts...).ToFunc()
}

// ByCompressedID orders the results by the compressed_id field.
func ByCompressedID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCompressedID, opts...).ToFunc()
}

// ByOriginalID orders the results by the original_id field.
func ByOriginalID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOriginalID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the createdAt field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updatedAt field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByUploadedBy orders the results by the uploadedBy field.
func ByUploadedBy(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUploadedBy, opts...).ToFunc()
}

// ByCategory orders the results by the category field.
func ByCategory(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCategory, opts...).ToFunc()
}

// ByCompressedField orders the results by compressed field.
func ByCompressedField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCompressedStep(), sql.OrderByField(field, opts...))
	}
}

// ByOriginalField orders the results by original field.
func ByOriginalField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOriginalStep(), sql.OrderByField(field, opts...))
	}
}

// ByThumbnail320xField orders the results by thumbnail320x field.
func ByThumbnail320xField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newThumbnail320xStep(), sql.OrderByField(field, opts...))
	}
}

// ByThumbnail768xField orders the results by thumbnail768x field.
func ByThumbnail768xField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newThumbnail768xStep(), sql.OrderByField(field, opts...))
	}
}

// ByUserField orders the results by user field.
func ByUserField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserStep(), sql.OrderByField(field, opts...))
	}
}

// ByMomentImagesCount orders the results by moment_images count.
func ByMomentImagesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMomentImagesStep(), opts...)
	}
}

// ByMomentImages orders the results by moment_images terms.
func ByMomentImages(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMomentImagesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newCompressedStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CompressedInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, CompressedTable, CompressedColumn),
	)
}
func newOriginalStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OriginalInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OriginalTable, OriginalColumn),
	)
}
func newThumbnail320xStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Thumbnail320xInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, Thumbnail320xTable, Thumbnail320xColumn),
	)
}
func newThumbnail768xStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Thumbnail768xInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, Thumbnail768xTable, Thumbnail768xColumn),
	)
}
func newUserStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
	)
}
func newMomentImagesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MomentImagesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MomentImagesTable, MomentImagesColumn),
	)
}
