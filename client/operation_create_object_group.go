package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (client *Client) CreateObjectGroup(ctx context.Context, req *CreateObjectGroupRequest) error {
	method := "POST"
	url := fmt.Sprintf("%s/Bucket/createObjectGroup", client.config.URL)

	bodyAsBytes, err := marshalCreateObjectGroupRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	httpResp, err := client.signAndDo(httpReq, bodyAsBytes)
	if err != nil {
		return fmt.Errorf("failed to %s to %s: %s", method, url, err)
	}
	defer httpResp.Body.Close()
	return nil
}

func marshalCreateObjectGroupRequest(req *CreateObjectGroupRequest) ([]byte, error) {
	log.Info("marshal class1 = ", req.Filter.ClassOne)

	log.Info("marshal class1 field = ", req.Filter.ClassOne.Field)
	log.Info("marshal class1 prefix = ", req.Filter.ClassOne.Prefix)

	log.Info("marshal class22 = ", req.Filter.ClassTwo)
	log.Info("marshal class22 field = ", req.Filter.ClassTwo.Field)
	log.Info("marshal class22 regx = ", req.Filter.ClassTwo.Regex)
// req.Filter.ClassOne,req.Filter.ClassTwo,
	body := map[string]interface{}{
		"bucket": req.Bucket,
		"source": req.Source,
		"format": map[string]interface{}{
			"_type":           req.Format.Type,
			"columnDelimiter": req.Format.ColumnDelimiter,
			"rowDelimiter":    req.Format.RowDelimiter,
			"headerRow":  req.Format.HeaderRow,
		},
		"filter" : []interface{}{
			req.Filter.ClassOne,req.Filter.ClassTwo,
			// {
			// 	"field":req.Filter.ClassOne.Field,
			// 	"prefix":req.Filter.ClassOne.Prefix,
			// },{
			// 	"field":req.Filter.ClassTwo.Field,
			// 	"regex":req.Filter.ClassTwo.Regex,
			// },
		},
		// "filter": []interface{}{
		// 	req.Filter.Obj1, req.Filter.Obj2,
		// },
		"indexRetention": map[string]interface{}{
			"forPartition": req.IndexRetention.ForPartition,
			"overall": req.IndexRetention.Overall,
		},
		"options": map[string]interface{}{
			"ignoreIrregular": req.Options.IgnoreIrregular,
		},
		"interval": map[string]interface{}{
			"mode":   req.Interval.Mode,
			"column": req.Interval.Column,
		},
		"realtime": req.Realtime,
	}

	log.Info("req body =========---> ", body)

	bodyAsBytes, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	return bodyAsBytes, nil
}
