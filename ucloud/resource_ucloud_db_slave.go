package ucloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudDBSlave() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudDBSlaveCreate,
		Read:   resourceUCloudDBSlaveRead,
		Update: resourceUCloudDBSlaveUpdate,
		Delete: resourceUCloudDBSlaveDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"master_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"is_lock": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validateInstancePassword,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDBInstanceName,
			},

			"instance_storage": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateDataDiskSize(20, 3000),
			},

			"parameter_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDBInstanceType,
			},

			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(3306, 65535),
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"expire_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"modify_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudDBSlaveCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	masterId := d.Get("master_id").(string)
	db, err := client.describeDBInstanceById(masterId)
	if err != nil {
		return fmt.Errorf("do %s failed in create db slave, %s", "DescribeUDBInstance", err)
	}
	arr := strings.Split(db.DBTypeId, "-")
	engine := arr[0]
	// skip error because it has been validated by schema
	dbType, _ := parseDBInstanceType(d.Get("instance_type").(string))
	if dbType.Engine != engine {
		return fmt.Errorf("error in create db slave, engine of instance type %s must be same as engine of master db %s", dbType.Engine, engine)
	}

	if dbType.Type == "ha" {
		return fmt.Errorf("error in create db slave, create high availability db slave is not supported")
	}

	req := conn.NewCreateUDBSlaveRequest()

	req.InstanceMode = ucloud.String(dbMap.convert(dbType.Type))
	req.SrcId = ucloud.String(masterId)
	req.Name = ucloud.String(d.Get("name").(string))
	req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
	req.MemoryLimit = ucloud.Int(dbType.Memory * 1000)
	req.InstanceType = ucloud.String("SATA_SSD")

	if val, ok := d.GetOk("port"); ok {
		req.Port = ucloud.Int(val.(int))
	} else {
		if engine == "mysql" || engine == "percona" {
			req.Port = ucloud.Int(3306)
		}
		if engine == "postgresql" {
			req.Port = ucloud.Int(5432)
		}
	}

	if val, ok := d.GetOk("is_lock"); ok {
		req.IsLock = ucloud.Bool(val.(bool))
	}

	resp, err := conn.CreateUDBSlave(req)

	if err != nil {
		return fmt.Errorf("error in create slave db, %s", err)
	}

	d.SetId(resp.DBId)

	// after create db slave, we need to wait it initialized
	stateConf := client.dbWaitForState(d.Id(), "Running")

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("wait for db slave initialize failed in create db slave %s, %s", d.Id(), err)
	}

	return resourceUCloudDBSlaveUpdate(d, meta)
}

func resourceUCloudDBSlaveUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	d.Partial(true)

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		req := conn.NewModifyUDBInstanceNameRequest()
		req.DBId = ucloud.String(d.Id())
		req.Name = ucloud.String(d.Get("name").(string))

		if _, err := conn.ModifyUDBInstanceName(req); err != nil {
			return fmt.Errorf("do %s failed in update db slave %s, %s", "ModifyUDBInstanceName", d.Id(), err)
		}

		// after update db slave name, we need to wait it completed
		stateConf := client.dbWaitForState(d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for update db name failed in update db slave %s, %s", d.Id(), err)
		}
	}

	if d.HasChange("password") && !d.IsNewResource() {
		d.SetPartial("password")
		req := conn.NewModifyUDBInstancePasswordRequest()
		req.DBId = ucloud.String(d.Id())
		req.Password = ucloud.String(d.Get("password").(string))

		if _, err := conn.ModifyUDBInstancePassword(req); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "ModifyUDBInstancePassword", d.Id(), err)
		}

		// after update db password, we need to wait it completed
		stateConf := client.dbWaitForState(d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for update db password failed in update db slave %s, %s", d.Id(), err)
		}
	}

	isChanged := false
	req := conn.NewResizeUDBInstanceRequest()
	req.DBId = ucloud.String(d.Id())

	if d.HasChange("instance_type") && !d.IsNewResource() {
		d.SetPartial("instance_type")
		old, new := d.GetChange("instance_type")

		oldType, _ := parseDBInstanceType(old.(string))

		newType, _ := parseDBInstanceType(new.(string))

		db, err := client.describeDBInstanceById(d.Get("master_id").(string))
		if err != nil {
			return fmt.Errorf("do %s failed in update db slave, %s", "DescribeUDBInstance", err)
		}
		arr := strings.Split(db.DBTypeId, "-")
		engine := arr[0]

		if newType.Engine != engine {
			return fmt.Errorf("error in update db slave, engine of instance type %s must be same as engine of master db %s", newType.Engine, engine)
		}

		if newType.Type == "ha" {
			return fmt.Errorf("error in update db slave, create high availability db slave is not supported")
		}

		if oldType.Memory != newType.Memory {
			req.MemoryLimit = ucloud.Int(newType.Memory * 1000)
		}

		isChanged = true
	}

	if d.HasChange("instance_storage") && !d.IsNewResource() {
		d.SetPartial("instance_storage")
		req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
		isChanged = true
	}

	if isChanged {
		//update these attributes of db slave, we need to wait it stopped
		stopReq := conn.NewStopUDBInstanceRequest()
		stopReq.DBId = ucloud.String(d.Id())
		stopReq.Zone = ucloud.String(d.Get("availability_zone").(string))

		dbSlave, err := client.describeDBInstanceById(d.Id())

		if err != nil {
			if isNotFoundError(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("do %s failed in update db slave %s, %s", "DescribeUDBInstance", d.Id(), err)
		}

		if dbSlave.State != "Shutoff" {
			_, err := conn.StopUDBInstance(stopReq)

			if err != nil {
				return fmt.Errorf("do %s failed in update db slave %s, %s", "StopUDBInstance", d.Id(), err)
			}

			// after stop db slave, we need to wait it stopped
			stateConf := client.dbWaitForState(d.Id(), "Shutoff")

			if _, err := stateConf.WaitForState(); err != nil {
				return fmt.Errorf("wait for stop db slave failed in update db slave %s, %s", d.Id(), err)
			}
		}

		if _, err := conn.ResizeUDBInstance(req); err != nil {
			return fmt.Errorf("do %s failed in update db slave %s, %s", "ResizeUDBInstance", d.Id(), err)
		}

		// after resize db slave, we need to wait it completed
		stateConf := client.dbWaitForState(d.Id(), "Shutoff")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for resize db slave failed in update db %s, %s", d.Id(), err)
		}

		// after update these attributes of db slave, we need to wait it started
		startReq := conn.NewStartUDBInstanceRequest()
		startReq.DBId = ucloud.String(d.Id())
		startReq.Zone = ucloud.String(d.Get("availability_zone").(string))

		_, err = conn.StartUDBInstance(startReq)

		if err != nil {
			return fmt.Errorf("do %s failed in update db slave %s, %s", "StartUDBInstance", d.Id(), err)
		}

		// after db slave update, we need to wait it started
		stateConf = client.dbWaitForState(d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for start db slave failed in update db slave %s, %s", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceUCloudDBSlaveRead(d, meta)
}

func resourceUCloudDBSlaveRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	db, err := client.describeDBInstanceById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read db slave %s, %s", "DescribeUDBInstance", d.Id(), err)
	}

	arr := strings.Split(db.DBTypeId, "-")
	d.Set("name", db.Name)
	d.Set("parameter_group_id", strconv.Itoa(db.ParamGroupId))
	d.Set("port", db.Port)
	d.Set("status", db.State)
	d.Set("instance_storage", db.DiskSpace)
	d.Set("availability_zone", db.Zone)
	d.Set("create_time", timestampToString(db.CreateTime))
	d.Set("expire_time", timestampToString(db.ExpiredTime))
	d.Set("modify_time", timestampToString(db.ModifyTime))
	var dbType dbInstanceType
	dbType.Memory = db.MemoryLimit / 1000
	dbType.Engine = arr[0]
	dbType.Type = dbMap.unconvert(db.InstanceMode)

	d.Set("instance_type", fmt.Sprintf("%s-%s-%d", dbType.Engine, dbType.Type, dbType.Memory))

	return nil
}

func resourceUCloudDBSlaveDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	req := conn.NewDeleteUDBInstanceRequest()
	req.DBId = ucloud.String(d.Id())
	stopReq := conn.NewStopUDBInstanceRequest()
	stopReq.DBId = ucloud.String(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		db, err := client.describeDBInstanceById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if db.State != "Shutoff" {
			if _, err := conn.StopUDBInstance(stopReq); err != nil {
				return resource.RetryableError(fmt.Errorf("do %s failed in delete db slave %s, %s", "StopUDBInstance", d.Id(), err))
			}

			// after db slave stop, we need to wait it stoped
			stateConf := client.dbWaitForState(d.Id(), "Shutoff")

			if _, err := stateConf.WaitForState(); err != nil {
				return resource.RetryableError(fmt.Errorf("wait for db slave stop failed in delete db %s, %s", d.Id(), err))
			}
		}

		if _, err := conn.DeleteUDBInstance(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error in delete db slave %s, %s", d.Id(), err))
		}

		if _, err := client.describeDBInstanceById(d.Id()); err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("do %s failed in delete db slave %s, %s", "DescribeUDBInstance", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("delete db slave but it still exists"))
	})
}
