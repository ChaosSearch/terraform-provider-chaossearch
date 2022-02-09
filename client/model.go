package client

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
	Name              string
	Compression       string
	FilterJSON        string
	Format            string
	LiveEventsSqsArn  string
	PartitionBy       string
	SourceBucket      string
	Pattern           string
	IndexRetention    int
	KeepOriginal      bool
	ArrayFlattenDepth *int
	ColumnRenames     map[string]interface{}
	ColumnSelection   map[string]interface{}
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
	AuthToken         string
	
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
