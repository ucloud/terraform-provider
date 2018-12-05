package ucloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/services/uaccount"
)

func dataSourceUCloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUCloudZonesRead,
		Schema: map[string]*schema.Schema{
			"output_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"zones": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUCloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*UCloudClient).uaccountconn

	req := conn.NewGetRegionRequest()

	resp, err := conn.GetRegion(req)
	if err != nil {
		return fmt.Errorf("error in read region list, %s", err)
	}

	var ids []string
	if v, ok := d.GetOk("ids"); ok {
		ids = ifaceToStringSlice(v)
	}

	var zones []uaccount.RegionInfo
	for _, item := range resp.Regions {
		if len(ids) == 0 || checkStringIn(item.Zone, ids) == nil {
			zones = append(zones, item)
		}
	}

	err = dataSourceUCloudZonesSave(d, zones, meta)
	if err != nil {
		return fmt.Errorf("error in read region list, %s", err)
	}

	return nil
}

func dataSourceUCloudZonesSave(d *schema.ResourceData, zones []uaccount.RegionInfo, meta interface{}) error {
	ids := []string{}
	data := []map[string]interface{}{}
	client := meta.(*UCloudClient)
	for _, item := range zones {
		if item.Region == client.region {
			ids = append(ids, item.Zone)
			data = append(data, map[string]interface{}{
				"id": item.Zone,
			})
		}
	}

	d.SetId(hashStringArray(ids))
	if err := d.Set("zones", data); err != nil {
		return err
	}
	d.Set("ids", ids)

	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), data)
	}

	return nil
}
