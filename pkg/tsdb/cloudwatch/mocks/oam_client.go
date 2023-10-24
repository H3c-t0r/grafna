package mocks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/oam"
	"github.com/stretchr/testify/mock"
)

type FakeOAMClient struct {
	mock.Mock
}

func (o *FakeOAMClient) ListSinksWithContext(ctx aws.Context, input *oam.ListSinksInput, opts ...request.Option) (*oam.ListSinksOutput, error) {
	args := o.Called(input)
	return args.Get(0).(*oam.ListSinksOutput), args.Error(1)
}

func (o *FakeOAMClient) ListAttachedLinksWithContext(ctx aws.Context, input *oam.ListAttachedLinksInput, opts ...request.Option) (*oam.ListAttachedLinksOutput, error) {
	args := o.Called(input)
	return args.Get(0).(*oam.ListAttachedLinksOutput), args.Error(1)
}
