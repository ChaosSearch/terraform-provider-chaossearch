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
	AuthToken string
	ID        string
}

type ReadViewRequest struct {
	AuthToken string
	ID        string
}
type Metadata struct {
	CreationDate int64 `json:"creationDate"`
}

type ObjectFilter struct {
	And []interface{} `json:"AND"`
}
type ReadObjectGroupResponse struct {
	Public      bool   `json:"_public"`
	Realtime    bool   `json:"_realtime"`
	Type        string `json:"_type"`
	Bucket      string `json:"bucket"`
	ContentType string `json:"contentType"`
	//Filter      struct {
	//	And []struct {
	//		Field  string `json:"field"`
	//		Prefix string `json:"prefix,omitempty"`
	//		Regex  string `json:"regex,omitempty"`
	//	} `json:"AND"`
	//} `json:"filter"`

	ObjectFilter ObjectFilter `json:"filter"`
	//Filter             *Filter
	Format             *Format   `json:"format"`
	ID                 string    `json:"id"`
	Interval           *Interval `json:"interval"`
	Metadata           *Metadata `json:"metadata"`
	Options            *Options  `json:"options"`
	RegionAvailability []string  `json:"regionAvailability"`
	Source             string    `json:"source"`

	Compression string
	FilterJSON  string
	//Format            string
	Pattern           string
	LiveEventsSqsArn  string
	PartitionBy       string
	SourceBucket      string
	IndexRetention    int
	KeepOriginal      bool
	ArrayFlattenDepth *int
	ColumnRenames     map[string]string
	ColumnSelection   []map[string]interface {
	}
}

type CreateObjectGroupRequest struct {
	AuthToken      string
	Bucket         string
	Source         string
	Format         *Format
	Interval       *Interval
	IndexRetention *IndexRetention
	Filter         *Filter
	Options        *Options
	Realtime       bool
}

//TODO add json value
type Format struct {
	Type            string `json:"_type"`
	ColumnDelimiter string `json:"columnDelimiter"`
	RowDelimiter    string `json:"rowDelimiter"`
	HeaderRow       bool   `json:"headerRow"`
}

type Interval struct {
	Mode   int
	Column int
}

type IndexRetention struct {
	ForPartition []interface{}
	Overall      int
}

type Filter struct {
	PrefixFilter *PrefixFilter
	RegexFilter  *RegexFilter
}

type PrefixFilter struct {
	Field  string `json:"field"`
	Prefix string `json:"prefix"`
}

type RegexFilter struct {
	Field string `json:"field"`
	Regex string `json:"regex"`
}

type Options struct {
	IgnoreIrregular bool
}
type UpdateIndexingStateRequest struct {
	ObjectGroupName string
	Active          bool
}

type DeleteObjectGroupRequest struct {
	AuthToken string
	Name      string
}

type DeleteViewRequest struct {
	AuthToken string
	Name      string
}

type UpdateObjectGroupRequest struct {
	AuthToken             string
	Bucket                string `json:"bucket"`
	IndexParallelism      int    `json:"indexParallelism"`
	IndexRetention        int    `json:"indexRetention"`
	TargetActiveIndex     int    `json:"targetActiveIndex"`
	LiveEventsParallelism int    `json:"liveEventsParallelism"`
}

type ReadIndexingStateRequest struct {
	ObjectGroupName string
}

type readBucketMetadataRequest struct {
	AuthToken  string
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
	FilterPredicate   *FilterPredicate `json:"filter"`
	TimeFieldName     string
	IndexPattern      string
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

type ReadViewResponse struct {
	Type               string `json:"_type"`
	Bucket             string
	FilterPredicate    *FilterPredicate `json:"filter"`
	TimeFieldName      string
	IndexPattern       string
	CaseInsensitive    bool
	ArrayFlattenDepth  *int
	IndexRetention     int `json:"overall"`
	Cacheable          bool
	Overwrite          bool
	Sources            []interface{}
	Transforms         []interface{}
	ID                 string    `json:"id"`
	MetaData           *Metadata `json:"metadata"`
	RegionAvailability []string  `json:"regionAvailability"`
	Compression        string
	LiveEventsSqsArn   string
	SourceBucket       string
	FilterJSON         string
	Pattern            string
	KeepOriginal       bool
}

type RequestHeaders struct {
	Headers map[string]interface{}
}

type FilterPredicate struct {
	Predicate *Predicate `json:"predicate"`
}

type Pred struct {
	Field string `json:"field"`
	Query string `json:"query"`
	State State  `json:"state"`
	Type_ string `json:"_type"`
}

type Predicate struct {
	Pred  Pred   `json:"pred"`
	Type_ string `json:"_type"`
}

type State struct {
	Type_ string `json:"_type"`
}

//user group create related models

type StartsWith struct {
	ChaosDocumentAttributesTitle string `json:"chaos:document/attributes.title"`
}
type Equals struct {
	ChaosDocumentAttributesTitle string `json:"chaos:document/attributes.title"`
}
type NotEquals struct {
	ChaosDocumentAttributesTitle string `json:"chaos:document/attributes.title"`
}
type Like struct {
	ChaosDocumentAttributesTitle string `json:"chaos:document/attributes.title"`
}

type Condition struct {
	StartsWith StartsWith
	Equals     Equals
	NotEquals  NotEquals
	Like       Like
}
type ConditionGroup struct {
	Condition []Condition `json:"Condition"`
}
type Permission struct {
	Effect         string
	Version        string
	Actions        []string
	Resources      []string
	ConditionGroup ConditionGroup `json:"Condition"`
}

type CreateUserGroupRequest struct {
	AuthToken  string
	Id         string
	Name       string
	Permission []Permission `json:"GroupIds"`
}

type UserInfoBlock struct {
	Username string `json:"Username"`
	FullName string `json:"FullName"`
	Email    string `json:"Email"`
}

type CreateSubAccountRequest struct {
	AuthToken     string
	UserInfoBlock UserInfoBlock `json:"UserInfoBlock"`
	GroupIds      []interface{} `json:"GroupIds"`
	Password      string
	HoCon         []interface{} `json:"HoCon"`
}

type ImportBucketRequest struct {
	AuthToken  string
	Bucket     string `json:"bucket"`
	HideBucket bool   `json:"hideBucket"`
}

type DeleteSubAccountRequest struct {
	AuthToken string
	Username  string
}

type ListUsersResponse struct {
	Users []User
}

type User struct {
	SubAccounts []SubAccount
}

type SubAccount struct {
	FullName  string
	Hocon     string
	Uid       string
	Username  string
	GroupIds  []interface{}
	Activated bool
}
