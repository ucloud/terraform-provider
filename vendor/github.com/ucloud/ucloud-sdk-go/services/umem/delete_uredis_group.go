//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UMem DeleteURedisGroup

package umem

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DeleteURedisGroupRequest is request schema for DeleteURedisGroup action
type DeleteURedisGroupRequest struct {
	request.CommonBase

	// 组ID
	GroupId *string `required:"true"`
}

// DeleteURedisGroupResponse is response schema for DeleteURedisGroup action
type DeleteURedisGroupResponse struct {
	response.CommonBase
}

// NewDeleteURedisGroupRequest will create request of DeleteURedisGroup action.
func (c *UMemClient) NewDeleteURedisGroupRequest() *DeleteURedisGroupRequest {
	req := &DeleteURedisGroupRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DeleteURedisGroup - 删除主备redis
func (c *UMemClient) DeleteURedisGroup(req *DeleteURedisGroupRequest) (*DeleteURedisGroupResponse, error) {
	var err error
	var res DeleteURedisGroupResponse

	err = c.client.InvokeAction("DeleteURedisGroup", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
