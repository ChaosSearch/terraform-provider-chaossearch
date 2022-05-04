package datasources

import (
	"context"
	"cs-tf-provider/client"
	"cs-tf-provider/provider/models"
	"cs-tf-provider/provider/resources"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceView() *schema.Resource {
	return &schema.Resource{
		ReadContext: resources.ResourceViewRead,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"cacheable": {
				Type:     schema.TypeBool,
				ForceNew: false,
				Optional: true,
			},
			"case_insensitive": {
				Type:     schema.TypeBool,
				ForceNew: false,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
			},
			"index_pattern": {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
			},
			"time_field_name": {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"predicate": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pred": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"type": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"query": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"state": {
													Type:     schema.TypeSet,
													Required: true,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:     schema.TypeString,
																Optional: true,
																ForceNew: true,
															},
														},
													},
												},
											},
										},
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_date": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"transforms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"region_availability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"array_flatten_depth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DataSourceViews() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceViewsRead,
		Schema: map[string]*schema.Schema{
			"views": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"case_insensitive": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"index_retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bucket_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"visible": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"time_field": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cacheable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_object_groups": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceViewsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	clientResponse, err := client.ListBuckets(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	objectGroups := GetBucketData(clientResponse)
	if err := data.Set("views", objectGroups); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func GetBucketData(clientResponse *client.ListBucketsResponse) []map[string]interface{} {
	result := make([]map[string]interface{}, len(clientResponse.BucketsCollection.Buckets))
	for i, bucket := range clientResponse.BucketsCollection.Buckets {
		tagMap := convertTagSetToMap(bucket.Tagging.TagSet)
		result[i] = map[string]interface{}{
			"name":                 bucket.Name,
			"creation_date":        bucket.CreationDate,
			"filter":               tagMap["cs3.filter"],
			"transform":            tagMap["cs3.transform"],
			"case_insensitive":     tagMap["cs3.case-insensitive"],
			"index_retention":      tagMap["cs3.index-retention"],
			"bucket_type":          tagMap["cs3.bucket-type"],
			"visible":              tagMap["cs3.visible"],
			"time_field":           tagMap["cs3.time-field"],
			"id":                   tagMap["cs3.dataset-id"],
			"index_pattern":        tagMap["cs3.index-pattern"],
			"cacheable":            tagMap["cs3.cacheable"],
			"parent_object_groups": tagMap["cs3.parent"],
		}
	}
	return result
}

func convertTagSetToMap(tags []client.Tag) map[string]interface{} {
	tagMap := make(map[string]interface{})
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}

	return tagMap
}
