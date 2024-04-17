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
	Name        string     `xml:"Name"`
	KeyCount    int        `xml:"KeyCount"`
	MaxKeys     int        `xml:"MaxKeys"`
	Delimiter   string     `xml:"Delimiter"`
	IsTruncated bool       `xml:"IsTruncated"`
	Contents    []Contents `xml:"Contents"`
}

type Contents struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int    `xml:"Size"`
	StorageClass string `xml:"StorageClass"`
}

type BasicRequest struct {
	AuthToken string
	Id        string
}

type Metadata struct {
	CreationDate int64 `json:"creationDate"`
}

type ReadBucketDatasetResp struct {
	Bucket             string       `json:"bucket"`
	ContentType        string       `json:"contentType"`
	ObjectFilter       ObjectFilter `json:"filter"`
	Format             Format       `json:"format"`
	Id                 string       `json:"id"`
	Interval           Interval     `json:"interval"`
	Options            Options      `json:"options"`
	RegionAvailability []string     `json:"regionAvailability"`
	Source             string       `json:"source"`
	State              State        `json:"state"`
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
	Pattern            string
	PartitionBy        interface{} `json:"partitionBy"`
	SourceBucket       string
	IndexRetention     int
	KeepOriginal       bool
	ArrayFlattenDepth  *int
	ColumnRenames      map[string]string
	ColumnSelection    []map[string]interface{}
}

type CreateObjectGroupRequest struct {
	AuthToken         string
	Bucket            string
	Source            string
	Format            *Format
	Interval          *Interval
	IndexRetention    *IndexRetention
	Filter            []Filter
	Options           *Options
	Realtime          bool
	LiveEventsAws     string
	LiveEventsGcp     *LiveEventsGcp
	PartitionBy       string
	TargetActiveIndex *int
}

type LiveEventsGcp struct {
	ProjectId      string
	SubscriptionId string
}

type Format struct {
	Type              string                   `json:"_type"`
	ColumnDelimiter   string                   `json:"columnDelimiter"`
	RowDelimiter      string                   `json:"rowDelimiter"`
	HeaderRow         bool                     `json:"headerRow"`
	ArrayFlattenDepth int                      `json:"arrayFlattenDepth"`
	StripPrefix       bool                     `json:"stripPrefix"`
	Horizontal        bool                     `json:"horizontal"`
	ArraySelection    []map[string]interface{} `json:"arraySelection"`
	FieldSelection    []map[string]interface{} `json:"fieldSelection"`
	VerticalSelection []map[string]interface{} `json:"verticalSelection"`
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
	IgnoreIrregular bool                     `json:"ignoreIrregular"`
	Compression     string                   `json:"compression"`
	ColTypes        map[string]string        `json:"colTypes,omitempty"`
	ColRenames      map[string]string        `json:"colRenames,omitempty"`
	ColSelection    []map[string]interface{} `json:"colSelection,omitempty"`
}

type UpdateObjectGroupRequest struct {
	AuthToken             string
	Bucket                string `json:"bucket"`
	IndexParallelism      *int   `json:"indexParallelism,omitempty"`
	IndexRetention        *int   `json:"indexRetention,omitempty"`
	TargetActiveIndex     *int   `json:"targetActiveIndex,omitempty"`
	LiveEventsParallelism *int   `json:"liveEventsParallelism,omitempty"`
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
	Transforms        []Transform
}

type Transform struct {
	Type         string          `json:"_type"`
	InputField   string          `json:"inputField"`
	OutputFields []ViewFieldSpec `json:"outputFields"`
	KeyPart      int             `json:"keyPart"`
	Pattern      string          `json:"pattern,omitempty"`
	Paths        []string        `json:"paths,omitempty"`
	Queries      []string        `json:"queries,omitempty"`
	Vertical     []string        `json:"vertical,omitempty"`
	Format       float32         `json:"format,omitempty"`
}

type ViewFieldSpec struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ReadViewResponse struct {
	Type               string           `json:"_type"`
	Bucket             string           `json:"bucket"`
	FilterPredicate    *FilterPredicate `json:"filter"`
	TimeFieldName      string           `json:"timeFieldName"`
	IndexPattern       string           `json:"indexPattern"`
	CaseInsensitive    bool             `json:"_caseInsensitive"`
	ArrayFlattenDepth  *int
	IndexRetention     int  `json:"overall"`
	Cacheable          bool `json:"_cacheable"`
	Overwrite          bool
	Sources            []string    `json:"sources"`
	Transforms         []Transform `json:"transforms"`
	ID                 string      `json:"id"`
	MetaData           *Metadata   `json:"metadata"`
	RegionAvailability []string    `json:"regionAvailability"`
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
	Type              string  `json:"_type,omitempty"`
	Status            *string `json:"status,omitempty"`
	TargetActiveIndex *int    `json:"targetActiveIndex,omitempty"`
}

// user group create related models
type Condition struct {
	StartsWith interface{} `json:"StartsWith,omitempty"`
	Equals     interface{} `json:"Equals,omitempty"`
	NotEquals  interface{} `json:"NotEquals,omitempty"`
	Like       interface{} `json:"Like,omitempty"`
}

type ConditionGroup struct {
	Conditions []*Condition `json:"Conditions,omitempty"`
}

type Permission struct {
	Effect         string
	Version        string
	Actions        []interface{}
	Resources      []interface{}
	ConditionGroup *ConditionGroup `json:"Condition,omitempty"`
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

type CreateMonitorRequest struct {
	Id         string      `json:"-"`
	AuthToken  string      `json:"-"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Enabled    bool        `json:"enabled"`
	Schedule   Schedule    `json:"schedule"`
	Inputs     []Input     `json:"inputs"`
	Triggers   []Trigger   `json:"triggers"`
	UIMetadata interface{} `json:"ui_metadata"`
}

type Schedule struct {
	Period Period `json:"period"`
}

type Period struct {
	Interval int    `json:"interval"`
	Unit     string `json:"unit"`
}

type Input struct {
	Search Search `json:"search"`
}

type Search struct {
	Indices []string    `json:"indices"`
	Query   interface{} `json:"query"`
}

type Trigger struct {
	Name      string           `json:"name"`
	Severity  string           `json:"severity"`
	MinTime   string           `json:"min_time_between_executions"`
	Condition MonitorCondition `json:"condition"`
	Actions   []Action         `json:"actions"`
}

type MonitorCondition struct {
	Script Script `json:"script"`
}

type Script struct {
	Lang   string `json:"lang"`
	Source string `json:"source"`
}

type Action struct {
	Name            string   `json:"name"`
	DestinationId   string   `json:"destination_id"`
	SubjectTemplate Script   `json:"subject_template"`
	MessageTemplate Script   `json:"message_template"`
	ThrottleEnabled bool     `json:"throttle_enabled,omitempty"`
	Throttle        Throttle `json:"throttle,omitempty"`
}

type Throttle struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type CreateMonitorResponse struct {
	Ok   bool        `json:"ok"`
	Resp MonitorResp `json:"resp"`
}

type MonitorResp struct {
	Id string `json:"_id"`
}

type CreateDestinationRequest struct {
	Id            string         `json:"-"`
	AuthToken     string         `json:"-"`
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	Slack         *Slack         `json:"slack,omitempty"`
	CustomWebhook *CustomWebhook `json:"custom_webhook,omitempty"`
}

type Slack struct {
	Url string `json:"url"`
}

type CustomWebhook struct {
	Scheme       string            `json:"scheme"`
	Method       string            `json:"method"`
	Url          string            `json:"url"`
	Host         string            `json:"host"`
	Port         int               `json:"port"`
	Path         string            `json:"path"`
	HeaderParams map[string]string `json:"header_params,omitempty"`
	QueryParams  map[string]string `json:"query_params,omitempty"`
}

type CreateDestinationResponse struct {
	Id      string `json:"id"`
	Ok      bool   `json:"ok"`
	Version int    `json:"version"`
}

type ReadDestinationResponse struct {
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	Id            string         `json:"id"`
	Type          string         `json:"type"`
	Name          string         `json:"name"`
	Slack         *Slack         `json:"slack,omitempty"`
	CustomWebhook *CustomWebhook `json:"custom_webhook,omitempty"`
}
