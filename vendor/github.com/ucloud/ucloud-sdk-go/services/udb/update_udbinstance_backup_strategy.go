//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDB UpdateUDBInstanceBackupStrategy

package udb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// UpdateUDBInstanceBackupStrategyRequest is request schema for UpdateUDBInstanceBackupStrategy action
type UpdateUDBInstanceBackupStrategyRequest struct {
	request.CommonBase

	// 可用区。参见 [可用区列表](../summary/regionlist.html)
	Zone *string `required:"false"`

	// 主节点的Id
	DBId *string `required:"true"`

	// 备份的整点时间，范围[0,23]
	BackupTime *int `required:"false"`

	// 备份时期标记位。共7位，每一位为一周中一天的备份情况，0表示关闭当天备份，1表示打开当天备份。最右边的一位为星期天的备份开关，其余从右到左依次为星期一到星期六的备份配置开关，每周必须至少设置两天备份。例如：1100000表示打开星期六和星期五的备份功能
	BackupDate *string `required:"false"`

	// 当导出某些数据遇到问题后，是否强制导出其他剩余数据默认是false
	ForceDump *bool `required:"false"`
}

// UpdateUDBInstanceBackupStrategyResponse is response schema for UpdateUDBInstanceBackupStrategy action
type UpdateUDBInstanceBackupStrategyResponse struct {
	response.CommonBase
}

// NewUpdateUDBInstanceBackupStrategyRequest will create request of UpdateUDBInstanceBackupStrategy action.
func (c *UDBClient) NewUpdateUDBInstanceBackupStrategyRequest() *UpdateUDBInstanceBackupStrategyRequest {
	req := &UpdateUDBInstanceBackupStrategyRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// UpdateUDBInstanceBackupStrategy - 修改UDB自动备份策略
func (c *UDBClient) UpdateUDBInstanceBackupStrategy(req *UpdateUDBInstanceBackupStrategyRequest) (*UpdateUDBInstanceBackupStrategyResponse, error) {
	var err error
	var res UpdateUDBInstanceBackupStrategyResponse

	err = c.client.InvokeAction("UpdateUDBInstanceBackupStrategy", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}