package client

type ListBucketsResponse struct {
	BucketsCollection BucketCollection `xml:"Buckets"`
}

type BucketCollection struct {
	Buckets []Bucket `xml:"Bucket"`
}

type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
	Tags         []Tag  `xml:"Tagging>TagSet>Tag"`
}

type Tag struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

type ListBucketResponse struct {
	Name        string    `xml:"Name"`
	KeyCount    int       `xml:"KeyCount"`
	MaxKeys     int       `xml:"MaxKeys"`
	Delimiter   string    `xml:"Delimiter"`
	IsTruncated bool      `xml:"IsTruncated"`
	Contents    *Contents `xml:"Contents"`
}

type Contents struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int    `xml:"Size"`
	StorageClass string `xml:"StorageClass"`
}

type ReadObjGroupReq struct {
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

type ReadObjGroupResp struct {
	Public             bool         `json:"_public"`
	Realtime           bool         `json:"_realtime"`
	Type               string       `json:"_type"`
	Bucket             string       `json:"bucket"`
	ContentType        string       `json:"contentType"`
	ObjectFilter       ObjectFilter `json:"filter"`
	Format             *Format      `json:"format"`
	ID                 string       `json:"id"`
	Interval           *Interval    `json:"interval"`
	Metadata           *Metadata    `json:"metadata"`
	Options            *Options     `json:"options"`
	RegionAvailability []string     `json:"regionAvailability"`
	Source             string       `json:"source"`
	Compression        string
	Pattern            string
	PartitionBy        string
	SourceBucket       string
	IndexRetention     int
	KeepOriginal       bool
	ArrayFlattenDepth  *int
	ColumnRenames      map[string]string
	ColumnSelection    []map[string]interface{}
}

type CreateObjectGroupRequest struct {
	AuthToken      string
	Bucket         string
	Source         string
	Format         *Format
	Interval       *Interval
	IndexRetention *IndexRetention
	Filter         []Filter
	Options        *Options
	Realtime       bool
	LiveEvents     string
}

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
	Field  string `json:"field"`
	Prefix string `json:"prefix,omitempty"`
	Regex  string `json:"regex,omitempty"`
	Equals string `json:"equals,omitempty"`
	Range  Range  `json:"range,omitempty"`
}

type Range struct {
	Min string `json:"min,omitempty"`
	Max string `json:"max,omitempty"`
}

type Options struct {
	IgnoreIrregular bool
	Compression     string
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

type CreateViewRequest struct {
	AuthToken         string
	Bucket            string
	FilterPredicate   *FilterPredicate `json:"filter"`
	TimeFieldName     string
	IndexPattern      string
	CaseInsensitive   bool
	ArrayFlattenDepth *int
	IndexRetention    int
	Cacheable         bool
	Overwrite         bool
	Sources           []interface{}
	Transforms        []interface{}
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

type FilterPredicate struct {
	Predicate *Predicate `json:"predicate"`
}

type Pred struct {
	Preds []Pred `json:"preds,omitempty"`
	Field string `json:"field,omitempty"`
	Query string `json:"query,omitempty"`
	State *State `json:"state,omitempty"`
	Type  string `json:"_type,omitempty"`
}

type Predicate struct {
	Pred  *Pred  `json:"pred,omitempty"`
	Preds []Pred `json:"preds,omitempty"`
	Type  string `json:"_type,omitempty"`
}

type State struct {
	Type string `json:"_type,omitempty"`
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
	Conditions []Condition `json:"Conditions"`
}

type Permission struct {
	Effect         string
	Version        string
	Actions        []interface{}
	Resources      []interface{}
	ConditionGroup ConditionGroup `json:"Condition"`
}

type CreateUserGroupRequest struct {
	AuthToken   string
	ID          string
	Name        string
	Permissions []Permission
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

type DeleteSubAccountRequest struct {
	AuthToken string
	Username  string
}

type ListUsersResponse struct {
	Users []User `json:"Users"`
}

type User struct {
	Activated   bool         `json:"Activated"`
	Deployed    bool         `json:"Deployed"`
	Email       string       `json:"Email"`
	FullName    string       `json:"FullName"`
	UserGroups  []UserGroup  `json:"Groups"`
	Hocon       string       `json:"Hocon"`
	InternalUid string       `json:"InternalUid"`
	IsHead      bool         `json:"IsHead"`
	Readonly    bool         `json:"Readonly"`
	Regions     []Region     `json:"Regions"`
	ServiceType string       `json:"ServiceType"`
	SubAccounts []SubAccount `json:"SubAccounts"`
	Uid         string       `json:"Uid"`
	Username    string       `json:"Username"`
}

type Region struct {
	Region string `json:"Region"`
	Uid    string `json:"Uid"`
}

type SubAccount struct {
	FullName  string   `json:"FullName"`
	Hocon     string   `json:"Hocon"`
	UID       string   `json:"Uid"`
	Username  string   `json:"Username"`
	GroupIds  []string `json:"GroupIds"`
	Activated bool     `json:"Activated"`
}

type UserGroup struct {
	ID          string       `json:"Id"`
	Name        string       `json:"Name"`
	Permissions []Permission `json:"permissions"`
}

type ReadUserGroupRequest struct {
	AuthToken string
	ID        string
}

type DeleteUserGroupRequest struct {
	AuthToken string
	ID        string
}

type IndexModelRequest struct {
	AuthToken  string
	BucketName string `json:"BucketName"`
	ModelMode  int    `json:"ModelMode"`
}

type IndexModelResponse struct {
	BucketName string `json:"BucketName"`
	Result     bool   `json:"Result"`
}

type IndexStatusResponse struct {
	Indexed bool `json:"indexed"`
}
