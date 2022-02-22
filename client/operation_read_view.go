package client

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"

	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"github.com/aws/aws-sdk-go/service/s3"
)

func (csClient *CSClient) ReadView(ctx context.Context, req *ReadViewRequest) (*ReadViewResponse, error) {
	var resp ReadViewResponse

	//if err := client.readViewAttributesFromBucketTagging(ctx, req, &resp); err != nil {
	//	return nil, err
	//}

	if err := csClient.readViewAttributesFromDatasetEndpoint(ctx, req, &resp); err != nil {
		return nil, err
	}

	log.Printf("ReadViewResponse: %+v", resp)

	return &resp, nil
}

func (csClient *CSClient) readViewAttributesFromDatasetEndpoint(ctx context.Context, req *ReadViewRequest, resp *ReadViewResponse) error {
	method := "GET"
	url := fmt.Sprintf("%s/Bucket/dataset/name/%s", csClient.config.URL, req.ID)
	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	authToken := req.AuthToken
	httpResp, err := csClient.signV2AndDo(authToken, httpReq, nil)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to Close response body  %s", err)
		}
	}(httpResp.Body)

	var read ReadViewResponse

	//var m ReadObjectGroupResponse
	if err := csClient.unmarshalJSONBody(httpResp.Body, &read); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response body: %s", err)
	}

	resp.FilterPredicate = read.FilterPredicate
	resp.Type = read.Type
	resp.MetaData = read.MetaData
	resp.RegionAvailability = read.RegionAvailability
	resp.ID = read.ID
	resp.Bucket = read.Bucket
	resp.Pattern = read.Pattern
	resp.Transforms = read.Transforms
	resp.TimeFieldName = read.TimeFieldName
	resp.Sources = read.Sources
	resp.Cacheable = read.Cacheable
	resp.CaseInsensitive = read.CaseInsensitive
	resp.IndexPattern = read.IndexPattern
	return nil
}

func (csClient *CSClient) readViewAttributesFromBucketTagging(ctx context.Context, req *ReadViewRequest, resp *ReadViewResponse) error {
	log.Printf("readViewAttributesFromBucketTagging")
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
	if err := mapViewBucketTaggingToResponse(tagging, resp); err != nil {
		return fmt.Errorf("failed to unmarshal XML response body: %s", err)
	}

	return nil
}

func mapViewBucketTaggingToResponse(tagging *s3.GetBucketTaggingOutput, v *ReadViewResponse) error {
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
	var Overall int
	if err := readJSONTagValue(tagging, "cs3.index-retention", &Overall); err != nil {
		return err
	}
	v.IndexRetention = Overall
	return nil
}
