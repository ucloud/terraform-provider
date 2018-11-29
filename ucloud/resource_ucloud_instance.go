package ucloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudInstanceCreate,
		Read:   resourceUCloudInstanceRead,
		Update: resourceUCloudInstanceUpdate,
		Delete: resourceUCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"root_password": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateInstancePassword,
			},

			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstanceType,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Instance",
				ValidateFunc: validateInstanceName,
			},

			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Month",
				ValidateFunc: validateStringInChoices([]string{"Year", "Month", "Dynamic"}),
			},

			"instance_duration": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"boot_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateDataDiskSize(20, 100),
			},

			"boot_disk_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "LOCAL_NORMAL",
				ValidateFunc: validateStringInChoices([]string{"LOCAL_NORMAL", "LOCAL_SSD", "CLOUD_NORMAL", "CLOUD_SSD"}),
			},

			"data_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateDataDiskSize(0, 2000),
			},

			"data_disk_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "LOCAL_NORMAL",
				ValidateFunc: validateStringInChoices([]string{"LOCAL_NORMAL", "LOCAL_SSD"}),
			},

			"remark": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
				Computed:     true,
			},

			"security_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},

						"disk_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"is_boot": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ip_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"expire_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"auto_renew": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uhostconn

	req := conn.NewCreateUHostInstanceRequest()
	req.LoginMode = ucloud.String("Password")
	req.Zone = ucloud.String(d.Get("availability_zone").(string))
	req.ImageId = ucloud.String(d.Get("image_id").(string))
	req.Password = ucloud.String(d.Get("root_password").(string))
	req.ChargeType = ucloud.String(d.Get("instance_charge_type").(string))
	req.Quantity = ucloud.Int(d.Get("instance_duration").(int))
	req.Name = ucloud.String(d.Get("name").(string))

	// skip error because it has been validated by schema
	t, _ := parseInstanceType(d.Get("instance_type").(string))
	req.CPU = ucloud.Int(t.CPU)
	req.Memory = ucloud.Int(t.Memory)

	imageResp, err := client.DescribeImageById(d.Get("image_id").(string))
	if err != nil {
		return fmt.Errorf("do %s failed in create instance, %s", "DescribeImage", err)
	}

	bootDisk := uhost.UHostDisk{}
	bootDisk.IsBoot = ucloud.String("True")
	bootDisk.Size = ucloud.Int(imageResp.ImageSize)
	bootDisk.Type = ucloud.String(d.Get("boot_disk_type").(string))

	req.Disks = append(req.Disks, bootDisk)

	if val, ok := d.GetOk("data_disk_size"); ok {
		dataDisk := uhost.UHostDisk{}
		dataDisk.IsBoot = ucloud.String("False")
		dataDisk.Type = ucloud.String(d.Get("data_disk_type").(string))
		dataDisk.Size = ucloud.Int(val.(int))

		req.Disks = append(req.Disks, dataDisk)
	}

	if val, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("vpc_id"); ok {
		req.VPCId = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("subnet_id"); ok {
		req.SubnetId = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("security_group"); ok {
		resp, err := client.describeFirewallById(val.(string))
		if err != nil {
			return fmt.Errorf("do %s failed in create instance, %s", "DescribeFirewall", err)
		}

		req.SecurityGroupId = ucloud.String(resp.GroupId)
	}

	resp, err := conn.CreateUHostInstance(req)
	if err != nil {
		return fmt.Errorf("error in create instance, %s", err)
	}

	if len(resp.UHostIds) != 1 {
		return fmt.Errorf("error in create instance, expect exactly one instance, got %v", len(resp.UHostIds))
	}

	d.SetId(resp.UHostIds[0])

	// after instance created, we need to wait it started
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"running"},
		Refresh:    instanceStateRefreshFunc(client, d.Id(), "running"),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("wait for instance start failed in create instance %s, %s", d.Id(), err)
	}

	return resourceUCloudInstanceUpdate(d, meta)
}

func resourceUCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uhostconn
	d.Partial(true)

	if d.HasChange("security_group") && !d.IsNewResource() {
		conn := client.unetconn
		d.SetPartial("security_group")
		req := conn.NewGrantFirewallRequest()
		req.FWId = ucloud.String(d.Get("security_group").(string))
		req.ResourceType = ucloud.String("UHost")
		req.ResourceId = ucloud.String(d.Id())

		_, err := conn.GrantFirewall(req)

		if err != nil {
			return fmt.Errorf("do %s failed in update instance %s, %s", "GrantFirewall", d.Id(), err)
		}
	}

	if d.HasChange("remark") {
		d.SetPartial("remark")
		req := conn.NewModifyUHostInstanceRemarkRequest()
		req.UHostId = ucloud.String(d.Id())
		req.Remark = ucloud.String(d.Get("remark").(string))

		_, err := conn.ModifyUHostInstanceRemark(req)

		if err != nil {
			return fmt.Errorf("error in set remark, %s", err)
		}
	}

	if d.HasChange("tag") && !d.IsNewResource() {
		d.SetPartial("tag")
		req := conn.NewModifyUHostInstanceTagRequest()
		req.UHostId = ucloud.String(d.Id())
		req.Tag = ucloud.String(d.Get("tag").(string))

		_, err := conn.ModifyUHostInstanceTag(req)

		if err != nil {
			return fmt.Errorf("do %s failed in update instance %s, %s", "ModifyUHostInstanceTag", d.Id(), err)
		}
	}

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		req := conn.NewModifyUHostInstanceNameRequest()
		req.UHostId = ucloud.String(d.Id())
		req.Name = ucloud.String(d.Get("name").(string))

		_, err := conn.ModifyUHostInstanceName(req)

		if err != nil {
			return fmt.Errorf("do %s failed in update instance %s, %s", "ModifyUHostInstanceName", d.Id(), err)
		}
	}

	resizeNeedUpdate := false
	resizeReq := conn.NewResizeUHostInstanceRequest()
	resizeReq.UHostId = ucloud.String(d.Id())
	if d.HasChange("instance_type") && !d.IsNewResource() {
		d.SetPartial("instance_type")
		oldType, newType := d.GetChange("instance_type")

		oldInstanceType, err := parseInstanceType(oldType.(string))

		if err != nil {
			return err
		}

		newInstanceType, err := parseInstanceType(newType.(string))

		if err != nil {
			return err
		}

		if oldInstanceType.CPU != newInstanceType.CPU {
			resizeReq.CPU = ucloud.Int(newInstanceType.CPU)
		}

		if oldInstanceType.Memory != newInstanceType.Memory {
			resizeReq.Memory = ucloud.Int(newInstanceType.Memory)
		}

		resizeNeedUpdate = true
	}

	if d.HasChange("data_disk_size") && !d.IsNewResource() {
		d.SetPartial("data_disk_size")
		oldSize, newSize := d.GetChange("data_disk_size")
		if oldSize.(int) > newSize.(int) {
			return fmt.Errorf("reduce data disk size is not supported, new value %d should be larger than the old value %d", newSize.(int), oldSize.(int))
		}
		resizeReq.DiskSpace = ucloud.Int(newSize.(int))
		resizeNeedUpdate = true
	}

	if d.HasChange("boot_disk_size") && !d.IsNewResource() {
		d.SetPartial("boot_disk_size")
		oldSize, newSize := d.GetChange("boot_disk_size")
		if oldSize.(int) > newSize.(int) {
			return fmt.Errorf("reduce boot disk size is not supported, new value %d by user set should be larger than the old value %d allocated by the system", newSize.(int), oldSize.(int))
		}
		resizeReq.BootDiskSpace = ucloud.Int(newSize.(int))
		resizeNeedUpdate = true
	}

	passwordNeedUpdate := false
	if d.HasChange("root_password") && !d.IsNewResource() {
		instance, err := client.describeInstanceById(d.Id())

		if err != nil {
			if isNotFoundError(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("do %s failed in update instance %s, %s", "DescribeUHostInstance", d.Id(), err)
		}

		if instance.BootDiskState == "Normal" {
			d.SetPartial("root_password")
			passwordNeedUpdate = true
		} else {
			return fmt.Errorf("update password must wait 20 minutes after the host starts up and then try again")
		}
	}

	if passwordNeedUpdate || resizeNeedUpdate {
		// instance update these attributes need to wait it stopped
		stopReq := conn.NewStopUHostInstanceRequest()
		stopReq.UHostId = ucloud.String(d.Id())

		instance, err := client.describeInstanceById(d.Id())

		if err != nil {
			if isNotFoundError(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("do %s failed in update instance %s, %s", "DescribeUHostInstance", d.Id(), err)
		}

		if instance.State != "Stopped" {
			_, err := conn.StopUHostInstance(stopReq)

			if err != nil {
				return fmt.Errorf("do %s failed in update instance %s, %s", "StopUHostInstance", d.Id(), err)
			}

			// after stop instance, we need to wait it stopped
			stateConf := &resource.StateChangeConf{
				Pending:    []string{"pending"},
				Target:     []string{"stopped"},
				Refresh:    instanceStateRefreshFunc(client, d.Id(), "stopped"),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      5 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			if _, err = stateConf.WaitForState(); err != nil {
				return fmt.Errorf("wait for instance stop failed in update instance %s, %s", d.Id(), err)
			}
		}

		if passwordNeedUpdate {
			reqPassword := conn.NewResetUHostInstancePasswordRequest()
			reqPassword.UHostId = ucloud.String(d.Id())
			reqPassword.Password = ucloud.String(d.Get("root_password").(string))

			_, err := conn.ResetUHostInstancePassword(reqPassword)

			if err != nil {
				return fmt.Errorf("do %s failed in update instance %s, %s", "ResetUHostInstancePassword", d.Id(), err)
			}
		}

		if resizeNeedUpdate {
			_, err := conn.ResizeUHostInstance(resizeReq)

			if err != nil {
				return fmt.Errorf("do %s failed in update instance %s, %s", "ResizeUHostInstance", d.Id(), err)
			}
		}

		// instance stopped means instance update complete
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending"},
			Target:     []string{"stopped"},
			Refresh:    instanceStateRefreshFunc(client, d.Id(), "stopped"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for instance update failed in update instance %s, %s", d.Id(), err)
		}

		// after instance update, we need to wait it started
		startReq := conn.NewStartUHostInstanceRequest()
		startReq.UHostId = ucloud.String(d.Id())

		if _, err := conn.StartUHostInstance(startReq); err != nil {
			return fmt.Errorf("do %s failed in update instance %s, %s", "StartUHostInstance", d.Id(), err)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"pending"},
			Target:     []string{"running"},
			Refresh:    instanceStateRefreshFunc(client, d.Id(), "running"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for instance start failed in update instance %s, %s", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceUCloudInstanceRead(d, meta)
}

func resourceUCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	instance, err := client.describeInstanceById(d.Id())

	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read instance %s, %s", "DescribeUHostInstance", d.Id(), err)
	}

	// TODO: [API-ERROR] image id is different between twice request
	d.Set("image_id", d.Get("image_id").(string))

	d.Set("name", instance.Name)
	d.Set("instance_charge_type", instance.ChargeType)
	d.Set("availability_zone", instance.Zone)
	d.Set("instance_type", d.Get("instance_type").(string))
	d.Set("root_password", d.Get("root_password").(string))
	d.Set("security_group", d.Get("security_group").(string))
	d.Set("tag", instance.Tag)
	d.Set("cpu", instance.CPU)
	d.Set("memory", instance.Memory)
	d.Set("state", instance.State)
	d.Set("create_time", timestampToString(instance.CreateTime))
	d.Set("expire_time", timestampToString(instance.ExpireTime))
	d.Set("auto_renew", instance.AutoRenew)
	d.Set("remark", instance.Remark)

	ipSet := []map[string]interface{}{}
	for _, item := range instance.IPSet {
		ipSet = append(ipSet, map[string]interface{}{
			"ip":   item.IP,
			"type": item.Type,
		})
	}
	d.Set("ip_set", ipSet)

	diskSet := []map[string]interface{}{}
	for _, item := range instance.DiskSet {
		diskSet = append(diskSet, map[string]interface{}{
			"disk_type": item.DiskType,
			"size":      item.Size,
			"disk_id":   item.DiskId,
			"is_boot":   item.IsBoot,
		})

		if item.IsBoot == "True" {
			d.Set("boot_disk_size", item.Size)
		}

		if item.IsBoot == "False" && checkStringIn(item.DiskType, []string{"LOCAL_NORMAL", "LOCAL_SSD"}) == nil {
			d.Set("data_disk_size", item.Size)
		}
	}

	d.Set("disk_set", diskSet)

	return nil
}

func resourceUCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.uhostconn

	stopReq := conn.NewStopUHostInstanceRequest()
	stopReq.UHostId = ucloud.String(d.Id())

	deleReq := conn.NewTerminateUHostInstanceRequest()
	deleReq.UHostId = ucloud.String(d.Id())

	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		instance, err := client.describeInstanceById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if instance.State != "Stopped" {
			if _, err := conn.StopUHostInstance(stopReq); err != nil {
				return resource.RetryableError(fmt.Errorf("do %s failed in delete instance %s, %s", "StopUHostInstance", d.Id(), err))
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"pending"},
				Target:     []string{"stopped"},
				Refresh:    instanceStateRefreshFunc(client, d.Id(), "stopped"),
				Timeout:    d.Timeout(schema.TimeoutDelete),
				Delay:      5 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			if _, err = stateConf.WaitForState(); err != nil {
				return resource.RetryableError(fmt.Errorf("wait for instance stop faild in delete instance %s, %s", d.Id(), err))
			}
		}

		if _, err := conn.TerminateUHostInstance(deleReq); err != nil {
			return resource.RetryableError(fmt.Errorf("error in delete instance %s, %s", d.Id(), err))
		}

		if _, err := client.describeInstanceById(d.Id()); err != nil {
			if isNotFoundError(err) {
				return nil
			}

			return resource.NonRetryableError(fmt.Errorf("do %s failed in delete instance %s, %s", "DescribeUHostInstance", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("delete instance but it still exists"))
	})
}

func instanceStateRefreshFunc(client *UCloudClient, instanceId, target string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := client.describeInstanceById(instanceId)
		if err != nil {
			if isNotFoundError(err) {
				return nil, "pending", nil
			}
			return nil, "", err
		}

		state := strings.ToLower(instance.State)
		if state != target {
			state = "pending"
		}

		return instance, state, nil
	}
}
