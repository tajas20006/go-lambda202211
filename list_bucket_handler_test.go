package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListObjects(t *testing.T) {
	cases := []struct {
		name      string
		output    s3.ListObjectsV2Output
		err       error
		expectErr bool
		expected  []Object
	}{
		{
			name:      "no output",
			output:    s3.ListObjectsV2Output{},
			err:       nil,
			expectErr: false,
			expected:  []Object{},
		},
		{
			name:      "err",
			output:    s3.ListObjectsV2Output{},
			err:       fmt.Errorf("failed"),
			expectErr: true,
		},
		{
			name: "one output",
			output: s3.ListObjectsV2Output{
				Contents: []types.Object{
					{
						Key:  aws.String("object1"),
						Size: 123,
					},
				},
			},
			err:       nil,
			expectErr: false,
			expected: []Object{
				{
					Key:   "object1",
					Bytes: 123,
				},
			},
		},
		{
			name: "two outputs",
			output: s3.ListObjectsV2Output{
				Contents: []types.Object{
					{
						Key:  aws.String("object1"),
						Size: 123,
					},
					{
						Key:  aws.String("object2"),
						Size: 234,
					},
				},
			},
			err:       nil,
			expectErr: false,
			expected: []Object{
				{
					Key:   "object1",
					Bytes: 123,
				},
				{
					Key:   "object2",
					Bytes: 234,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			ctrl := gomock.NewController(t)
			// `defer ctrl.Finish()` は v1.5.0 以降不要になっている
			assert := assert.New(t)
			ctx := context.Background()

			s3mock := NewMocks3API(ctrl)
			s3mock.EXPECT().ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(S3BUCKET)},
			).Return(&tt.output, tt.err)

			// Execute
			actual, err := listObjects(ctx, s3mock, S3BUCKET)

			// Assert
			assert.True(tt.expectErr == (err != nil))
			if !tt.expectErr {
				assert.Equal(tt.expected, actual)
			}
		})
	}
}
