//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UMem DeleteUMemSpace

package umem

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DeleteUMemSpaceRequest is request schema for DeleteUMemSpace action
type DeleteUMemSpaceRequest struct {
	request.CommonBase

	// 可用区。参见 [可用区列表](../summary/regionlist.html)
	Zone *string `required:"false"`

	// UMem内存空间ID
	SpaceId *string `required:"true"`
}

// DeleteUMemSpaceResponse is response schema for DeleteUMemSpace action
type DeleteUMemSpaceResponse struct {
	response.CommonBase
}

// NewDeleteUMemSpaceRequest will create request of DeleteUMemSpace action.
func (c *UMemClient) NewDeleteUMemSpaceRequest() *DeleteUMemSpaceRequest {
	req := &DeleteUMemSpaceRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DeleteUMemSpace - 删除UMem内存空间
func (c *UMemClient) DeleteUMemSpace(req *DeleteUMemSpaceRequest) (*DeleteUMemSpaceResponse, error) {
	var err error
	var res DeleteUMemSpaceResponse

	err = c.client.InvokeAction("DeleteUMemSpace", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}