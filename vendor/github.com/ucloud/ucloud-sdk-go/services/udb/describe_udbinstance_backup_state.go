//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDB DescribeUDBInstanceBackupState

package udb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DescribeUDBInstanceBackupStateRequest is request schema for DescribeUDBInstanceBackupState action
type DescribeUDBInstanceBackupStateRequest struct {
	request.CommonBase

	// 可用区。参见 [可用区列表](../summary/regionlist.html)
	Zone *string `required:"true"`

	// 备份记录ID
	BackupId *int `required:"true"`

	// 跨可用区高可用备库所在可用区，参见［可用区列表］
	BackupZone *string `required:"false"`
}

// DescribeUDBInstanceBackupStateResponse is response schema for DescribeUDBInstanceBackupState action
type DescribeUDBInstanceBackupStateResponse struct {
	response.CommonBase

	// 备份状态 0 Backuping // 备份中 1 Success // 备份成功 2 Failed // 备份失败 3 Expired // 备份过期
	State string

	// 备份所占空间大小
	BackupSize int

	// 备份截止时间
	BackupEndTime int
}

// NewDescribeUDBInstanceBackupStateRequest will create request of DescribeUDBInstanceBackupState action.
func (c *UDBClient) NewDescribeUDBInstanceBackupStateRequest() *DescribeUDBInstanceBackupStateRequest {
	req := &DescribeUDBInstanceBackupStateRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DescribeUDBInstanceBackupState - 获取UDB实例备份状态
func (c *UDBClient) DescribeUDBInstanceBackupState(req *DescribeUDBInstanceBackupStateRequest) (*DescribeUDBInstanceBackupStateResponse, error) {
	var err error
	var res DescribeUDBInstanceBackupStateResponse

	err = c.client.InvokeAction("DescribeUDBInstanceBackupState", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
