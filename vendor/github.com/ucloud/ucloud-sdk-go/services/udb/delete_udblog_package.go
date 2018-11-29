//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDB DeleteUDBLogPackage

package udb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DeleteUDBLogPackageRequest is request schema for DeleteUDBLogPackage action
type DeleteUDBLogPackageRequest struct {
	request.CommonBase

	// 可用区。参见 [可用区列表](../summary/regionlist.html)
	Zone *string `required:"true"`

	// 日志包id，可通过DescribeUDBLogPackage获得
	BackupId *int `required:"true"`

	// 跨可用区高可用备库所在可用区
	BackupZone *string `required:"false"`
}

// DeleteUDBLogPackageResponse is response schema for DeleteUDBLogPackage action
type DeleteUDBLogPackageResponse struct {
	response.CommonBase
}

// NewDeleteUDBLogPackageRequest will create request of DeleteUDBLogPackage action.
func (c *UDBClient) NewDeleteUDBLogPackageRequest() *DeleteUDBLogPackageRequest {
	req := &DeleteUDBLogPackageRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DeleteUDBLogPackage - 删除UDB日志包
func (c *UDBClient) DeleteUDBLogPackage(req *DeleteUDBLogPackageRequest) (*DeleteUDBLogPackageResponse, error) {
	var err error
	var res DeleteUDBLogPackageResponse

	err = c.client.InvokeAction("DeleteUDBLogPackage", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
