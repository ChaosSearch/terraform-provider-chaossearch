package datasources

import (
	"context"
	"cs-tf-provider/provider/models"
	"cs-tf-provider/provider/resources"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceObjectGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: resources.ResourceObjectGroupRead,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"format": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"column_delimiter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"header_row": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"row_delimiter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"array_flatten_depth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"strip_prefix": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"horizontal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"array_selection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_selection": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"live_events": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"index_retention": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"overall": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"equals": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"regex": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"interval": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ignore_irregular": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"col_types": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region_availability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"index_parallelism": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_active_index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"live_events_parallelism": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DataSourceObjectGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectGroupsRead,
		Schema: map[string]*schema.Schema{
			"object_groups": {
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
						"case_insensitive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_retention": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"visible": {
							Type:     schema.TypeString,
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
						"transform": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceObjectGroupsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*models.ProviderMeta).CSClient
	tokenValue := meta.(*models.ProviderMeta).Token
	clientResponse, err := client.ListBuckets(ctx, tokenValue)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("object_groups", GetBucketData(clientResponse)); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}
