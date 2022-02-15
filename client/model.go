package client

//import "container/list"

type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

type BucketCollection struct {
	Buckets []Bucket `xml:"Bucket"`
}

type ListBucketsResponse struct {
	BucketsCollection BucketCollection `xml:"Buckets"`
}

type ReadObjectGroupRequest struct {
	ID string
}

type ReadObjectGroupResponse struct {
	Compression       string
	FilterJSON        string
	Format            string
	Pattern           string
	LiveEventsSqsArn  string
	PartitionBy       string
	SourceBucket      string
	IndexRetention    int
	KeepOriginal      bool
	ArrayFlattenDepth *int
	ColumnRenames     map[string]string
	ColumnSelection   []map[string]interface{}
}

type CreateObjectGroupRequest struct {
	Bucket         string
	Source         string
	Format         *Format
	Interval       *Interval
	IndexRetention *IndexRetention
	Filter         *[]Filter
	Options        *Options
	Realtime       bool
}

type Format struct {
	Type            string
	ColumnDelimiter string
	RowDelimiter    string
	HeaderRow       bool
}

type Interval struct {
	Mode   int
	Column int
}

type IndexRetention struct {
	ForPartition []interface {}
	Overall       int
}

type Filter struct {
	Obj1 *ClassOne
	Obj2 *ClassTwo
}

type ClassOne struct {
	Field  string
	Prefix string
}

type ClassTwo struct {
	Field  string
	Regex string
}

type Options struct {
	IgnoreIrregular bool
}
type UpdateIndexingStateRequest struct {
	ObjectGroupName string
	Active          bool
}

type DeleteObjectGroupRequest struct {
	Name string
}

type UpdateObjectGroupRequest struct {
	Name           string
	IndexRetention int
}

type ReadIndexingStateRequest struct {
	ObjectGroupName string
}

type readBucketMetadataRequest struct {
	BucketName string `json:"BucketName"`
	Stats      bool   `json:"Stats"`
}

type IndexingState struct {
	ObjectGroupName string
	Active          bool
}

type CreateViewRequest struct {
	AuthToken string

	Bucket            string
	FilterJSON        string
	TimeFieldName     string
	Pattern           string
	CaseInsensitive   bool
	ArrayFlattenDepth *int
	IndexRetention    int
	// IndexRetention    map[string]interface{}
	Cacheable bool
	Overwrite bool
	// Sources           map[string]string
	Sources    []interface{}
	Transforms []interface{}
}

type RequestHeaders struct {
	Headers map[string]interface{}
}
