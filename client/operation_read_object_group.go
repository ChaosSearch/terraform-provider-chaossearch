package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type appLogger struct{}

func (l appLogger) Log(args ...interface{}) {
	log.Printf("AWS: %+v", args...)
}

func (csClient *CSClient) ReadObjectGroup(ctx context.Context, req *ReadObjectGroupRequest) (*ReadObjectGroupResponse, error) {
	var resp ReadObjectGroupResponse

	//if err := client.readAttributesFromBucketTagging(ctx, req, &resp); err != nil {
	//	return nil, err
	//}

	if err := csClient.readAttributesFromDatasetEndpoint(ctx, req, &resp); err != nil {
		return nil, err
	}

	log.Printf("ReadObjectGroupResponse: %+v", resp)

	return &resp, nil
}

func (csClient *CSClient) readAttributesFromDatasetEndpoint(ctx context.Context, req *ReadObjectGroupRequest, resp *ReadObjectGroupResponse) error {
	method := "GET"
	url := fmt.Sprintf("%s/Bucket/dataset/name/%s", csClient.config.URL, req.ID)

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	var sessionToken = req.AuthToken
	httpResp, err := csClient.signV2AndDo(sessionToken, httpReq, nil)
	//httpResp, err := client.signV4AndDo(httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var ReadObjectGroup ReadObjectGroupResponse
	if err := csClient.unmarshalJSONBody(httpResp.Body, &ReadObjectGroup); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response body: %s", err)
	}
	resp.Format = ReadObjectGroup.Format
	resp.Filter = ReadObjectGroup.Filter
	resp.Interval = ReadObjectGroup.Interval
	resp.Metadata = ReadObjectGroup.Metadata
	resp.Options = ReadObjectGroup.Options
	resp.RegionAvailability = ReadObjectGroup.RegionAvailability
	resp.Public = ReadObjectGroup.Public
	resp.Realtime = ReadObjectGroup.Realtime
	resp.Type = ReadObjectGroup.Type
	resp.Bucket = ReadObjectGroup.Bucket
	resp.ContentType = ReadObjectGroup.ContentType
	resp.ID = ReadObjectGroup.ID
	resp.Source = ReadObjectGroup.Source
	//resp.IndexRetention = ReadObjectGroup.IndexRetention
	resp.Compression = ReadObjectGroup.Compression
	resp.PartitionBy = ReadObjectGroup.PartitionBy
	resp.Pattern = ReadObjectGroup.Pattern
	resp.SourceBucket = ReadObjectGroup.SourceBucket
	resp.ColumnSelection = ReadObjectGroup.ColumnSelection
	return nil
}

func (csClient *CSClient) readAttributesFromBucketTagging(ctx context.Context, req *ReadObjectGroupRequest, resp *ReadObjectGroupResponse) error {
	session_, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(csClient.config.AccessKeyID, csClient.config.SecretAccessKey, ""),
		Endpoint:         aws.String(fmt.Sprintf("%s/V1", csClient.config.URL)),
		Region:           aws.String(csClient.config.Region),
		S3ForcePathStyle: aws.Bool(true),
		LogLevel:         aws.LogLevel(aws.LogOff),
		Logger:           appLogger{},
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %s", err)
	}

	svc := s3.New(session_)
	input := &s3.GetBucketTaggingInput{
		Bucket: aws.String(req.ID),
	}

	tagging, err := svc.GetBucketTaggingWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to read bucket tagging: %s", err)
	}

	if err := mapBucketTaggingToResponse(tagging, resp); err != nil {
		return fmt.Errorf("failed to unmarshal XML response body: %s", err)
	}

	return nil
}

func mapBucketTaggingToResponse(tagging *s3.GetBucketTaggingOutput, v *ReadObjectGroupResponse) error {
	if err := readStringTagValue(tagging, "cs3.parent", &v.SourceBucket); err != nil {
		return err
	}

	if err := readStringTagValue(tagging, "cs3.compression", &v.Compression); err != nil {
		return err
	}

	if err := readStringTagValue(tagging, "cs3.live-sqs-arn", &v.LiveEventsSqsArn); err != nil {
		return err
	}

	var filterObject struct {
		Type              string `json:"_type"`
		Pattern           string `json:"pattern"`
		ArrayFlattenDepth *int   `json:"arrayFlattenDepth"`
		KeepOriginal      bool   `json:"keepOriginal"`
	}
	if err := readJSONTagValue(tagging, "cs3.dataset-format", &filterObject); err != nil {
		return err
	}
	//v.Format = filterObject.Type
	v.Pattern = filterObject.Pattern
	v.ArrayFlattenDepth = filterObject.ArrayFlattenDepth
	v.KeepOriginal = filterObject.KeepOriginal

	if err := readStringTagValue(tagging, "cs3.predicate", &v.FilterJSON); err != nil {
		return err
	}
	var retentionObject struct {
		Overall int `json:"overall"`
	}
	if err := readJSONTagValue(tagging, "cs3.index-retention", &retentionObject); err != nil {
		return err
	}
	v.IndexRetention = retentionObject.Overall
	return nil
}

func readStringTagValue(tagging *s3.GetBucketTaggingOutput, key string, v *string) error {
	stringValue, err := findTagValue(tagging, key)
	if err != nil {
		return nil
	}

	*v = stringValue
	return nil
}

func readJSONTagValue(tagging *s3.GetBucketTaggingOutput, key string, v interface{}) error {
	valueAsBytes, err := findTagValue(tagging, key)
	if err != nil {
		return nil
	}

	return json.Unmarshal([]byte(valueAsBytes), v)
}

func findTagValue(tagging *s3.GetBucketTaggingOutput, key string) (string, error) {
	for _, tag := range tagging.TagSet {
		if *tag.Key == key {
			return *tag.Value, nil
		}
	}

	return "", fmt.Errorf("no tag found with key: %s", key)
}
