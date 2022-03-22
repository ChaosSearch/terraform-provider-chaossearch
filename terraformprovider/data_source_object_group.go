package cs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceObjectGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceObjectGroupRead,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"_public": {
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
				Optional: true,
			},
			"_realtime": {
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
				Optional: true,
			},
			"_type": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
			},
			"format": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"column_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"header_row": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"row_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"regex_filter": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"regex": {
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
			"interval": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"mode": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_irregular": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
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
			"region_availability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}
