package ucloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudLBListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudLBListenerCreate,
		Update: resourceUCloudLBListenerUpdate,
		Read:   resourceUCloudLBListenerRead,
		Delete: resourceUCloudLBListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP",
					"HTTPS",
					"TCP",
					"UDP",
				}, false),
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      resource.PrefixedUniqueId("tf-listener-"),
				ValidateFunc: validateName,
			},

			"listen_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RequestProxy",
				ValidateFunc: validation.StringInSlice([]string{
					"RequestProxy",
					"PacketsTransmit",
				}, false),
			},

			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validation.IntBetween(1, 65535),
			},

			"idle_timeout": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},

			"method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Roundrobin",
				ValidateFunc: validation.StringInSlice([]string{
					"Roundrobin",
					"Path",
					"SourcePort",
					"ConsistentHash",
					"ConsistentHashPort",
				}, false),
			},

			"persistence_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "None",
				ValidateFunc: validation.StringInSlice([]string{
					"ServerInsert",
					"UserDefined",
					"None",
				}, false),
			},

			"persistence": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"health_check_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Port",
					"Path",
				}, false),
			},

			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudLBListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.ulbconn

	lbId := d.Get("load_balancer_id").(string)

	req := conn.NewCreateVServerRequest()
	req.ULBId = ucloud.String(lbId)
	req.Protocol = ucloud.String(d.Get("protocol").(string))
	req.ListenType = ucloud.String(d.Get("listen_type").(string))
	req.FrontendPort = ucloud.Int(d.Get("port").(int))
	req.Method = ucloud.String(d.Get("method").(string))
	req.VServerName = ucloud.String(d.Get("name").(string))

	if val, ok := d.GetOk("idle_timeout"); ok {
		req.ClientTimeout = ucloud.Int(val.(int))
	}

	if val, ok := d.GetOk("persistence_type"); ok {
		req.PersistenceType = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("persistence"); ok {
		req.PersistenceInfo = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("health_check_type"); ok {
		req.MonitorType = ucloud.String(val.(string))
		if val == "Path" {

			if val, ok := d.GetOk("domain"); ok {
				req.Domain = ucloud.String(val.(string))
			}

			if val, ok := d.GetOk("path"); ok {
				req.Path = ucloud.String(val.(string))
			}

		}
	}

	resp, err := conn.CreateVServer(req)
	if err != nil {
		return fmt.Errorf("error in create lb listener, %s", err)
	}

	d.SetId(resp.VServerId)

	// after create lb listener, we need to wait it initialized
	stateConf := lbListenerWaitForState(client, lbId, d.Id())

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("wait for lb listener initialize failed in create lb listener %s, %s", d.Id(), err)
	}

	return resourceUCloudLBListenerUpdate(d, meta)
}

func resourceUCloudLBListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*UCloudClient).ulbconn

	d.Partial(true)

	isChanged := false
	req := conn.NewUpdateVServerAttributeRequest()
	req.ULBId = ucloud.String(d.Get("load_balancer_id").(string))
	req.VServerId = ucloud.String(d.Id())

	if d.HasChange("name") && !d.IsNewResource() {
		isChanged = true
		req.VServerName = ucloud.String(d.Get("name").(string))
	}

	if d.HasChange("protocol") && !d.IsNewResource() {
		isChanged = true
		req.Protocol = ucloud.String(d.Get("protocol").(string))
	}

	if d.HasChange("method") && !d.IsNewResource() {
		isChanged = true
		req.Method = ucloud.String(d.Get("method").(string))
	}

	if d.HasChange("persistence_type") && !d.IsNewResource() {
		isChanged = true
		req.PersistenceType = ucloud.String(d.Get("persistence_type").(string))
	}

	if d.HasChange("persistence") && !d.IsNewResource() {
		isChanged = true
		req.PersistenceInfo = ucloud.String(d.Get("persistence").(string))
	}

	if d.HasChange("idle_timeout") && !d.IsNewResource() {
		isChanged = true
		req.ClientTimeout = ucloud.Int(d.Get("idle_timeout").(int))
	}

	if d.HasChange("health_check_type") && !d.IsNewResource() {
		isChanged = true
		req.MonitorType = ucloud.String(d.Get("health_check_type").(string))
	}

	if d.HasChange("domain") && !d.IsNewResource() {
		isChanged = true
		req.Domain = ucloud.String(d.Get("domain").(string))
	}

	if d.HasChange("path") && !d.IsNewResource() {
		isChanged = true
		req.Path = ucloud.String(d.Get("path").(string))
	}

	if isChanged {
		_, err := conn.UpdateVServerAttribute(req)
		if err != nil {
			return fmt.Errorf("do %s failed in update lb listener %s, %s", "UpdateVServerAttribute", d.Id(), err)
		}

		d.SetPartial("name")
		d.SetPartial("protocol")
		d.SetPartial("method")
		d.SetPartial("persistence_type")
		d.SetPartial("persistence")
		d.SetPartial("idle_timeout")
		d.SetPartial("health_check_type")
		d.SetPartial("domain")
		d.SetPartial("path")
	}

	d.Partial(false)

	return resourceUCloudLBListenerRead(d, meta)
}

func resourceUCloudLBListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	lbId := d.Get("load_balancer_id").(string)
	vserverSet, err := client.describeVServerById(lbId, d.Id())

	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read lb listener %s, %s", "DescribeVServer", d.Id(), err)
	}

	d.Set("name", vserverSet.VServerName)
	d.Set("protocol", vserverSet.Protocol)
	d.Set("listen_type", vserverSet.ListenType)
	d.Set("port", vserverSet.FrontendPort)
	d.Set("idle_timeout", vserverSet.ClientTimeout)
	d.Set("method", vserverSet.Method)
	d.Set("persistence_type", vserverSet.PersistenceType)
	d.Set("persistence", vserverSet.PersistenceInfo)
	d.Set("health_check_type", vserverSet.MonitorType)
	d.Set("domain", vserverSet.Domain)
	d.Set("path", vserverSet.Path)
	d.Set("status", listenerStatusCvt.mustConvert(vserverSet.Status))

	return nil
}

func resourceUCloudLBListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.ulbconn
	lbId := d.Get("load_balancer_id").(string)

	req := conn.NewDeleteVServerRequest()
	req.ULBId = ucloud.String(lbId)
	req.VServerId = ucloud.String(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeleteVServer(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error in delete lb listener %s, %s", d.Id(), err))
		}

		_, err := client.describeVServerById(lbId, d.Id())

		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("do %s failed in delete lb listener %s, %s", "DescribeVServer", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("delete lb listener but it still exists"))
	})
}

func lbListenerWaitForState(client *UCloudClient, lbId, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{statusInitialized},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
		Refresh: func() (interface{}, string, error) {
			vserverSet, err := client.describeVServerById(lbId, id)
			if err != nil {
				if isNotFoundError(err) {
					return nil, statusPending, nil
				}
				return nil, "", err
			}

			state := listenerStatusCvt.mustConvert(vserverSet.Status)
			if state != "allNormal" {
				state = statusPending
			} else {
				state = statusInitialized
			}

			return vserverSet, state, nil
		},
	}
}
