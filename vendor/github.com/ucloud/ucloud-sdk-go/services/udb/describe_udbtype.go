//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDB DescribeUDBType

package udb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DescribeUDBTypeRequest is request schema for DescribeUDBType action
type DescribeUDBTypeRequest struct {
	request.CommonBase

	// 可用区。参见 [可用区列表](../summary/regionlist.html)
	Zone *string `required:"true"`
}

// DescribeUDBTypeResponse is response schema for DescribeUDBType action
type DescribeUDBTypeResponse struct {
	response.CommonBase

	// DB类型列表 参数见 UDBTypeSet
	DataSet []UDBTypeSet
}

// NewDescribeUDBTypeRequest will create request of DescribeUDBType action.
func (c *UDBClient) NewDescribeUDBTypeRequest() *DescribeUDBTypeRequest {
	req := &DescribeUDBTypeRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DescribeUDBType - 获取UDB支持的类型信息
func (c *UDBClient) DescribeUDBType(req *DescribeUDBTypeRequest) (*DescribeUDBTypeResponse, error) {
	var err error
	var res DescribeUDBTypeResponse

	err = c.client.InvokeAction("DescribeUDBType", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
