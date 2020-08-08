package mocks

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3crypto"
)

type MockAWSS3Client struct {
	objects map[string]string
	headers map[string]http.Header
}

func NewMockAWSS3Client() *MockAWSS3Client {
	return &MockAWSS3Client{
		objects: make(map[string]string),
		headers: make(map[string]http.Header),
	}
}

func (c MockAWSS3Client) GetObjectDirect(bucket string, key string) ([]byte, http.Header, error) {
	bk := encodeBucketKey(bucket, key)
	enc, ok := c.objects[bk]
	if !ok {
		return nil, nil, fmt.Errorf("Bucket %q not found", bk)
	}
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, nil, err
	}
	header, ok := c.headers[bk]
	if !ok {
		return nil, nil, fmt.Errorf("Bucket %q not found", bk)
	}
	return data, header, nil

}

func (c *MockAWSS3Client) PutObjectDirect(bucket string, key string, data []byte, header http.Header) error {
	bk := encodeBucketKey(bucket, key)
	c.objects[bk] = base64.StdEncoding.EncodeToString(data)
	c.headers[bk] = header
	return nil
}

func (c *MockAWSS3Client) MockPutObjectRequest(req *request.Request) error {
	req.Handlers.Send.Clear()
	req.Handlers.Send.PushBack(func(r *request.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			r.Error = err
			return
		}
		in, ok := req.Params.(*s3.PutObjectInput)
		if !ok {
			r.Error = fmt.Errorf("could not cast params as PutObjectInput, got %T", req.Params)
			return
		}
		header := http.Header{}
		for key, value := range in.Metadata {
			header[http.CanonicalHeaderKey("x-amz-meta-"+strings.ToLower(key))] = []string{*value}
		}
		header[http.CanonicalHeaderKey("X-Amzn-Requestid")] = []string{"1"}
		c.PutObjectDirect(*in.Bucket, *in.Key, data, header)
		r.HTTPResponse = &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte{})),
		}
	})
	return nil
}

func (c MockAWSS3Client) MockGetObjectRequest(req *request.Request, out *s3.GetObjectOutput) error {
	in, ok := req.Params.(*s3.GetObjectInput)
	if !ok {
		return fmt.Errorf("could not cast params as GetObjectInput, got %T", req.Params)
	}
	req.Handlers.Send.Clear()
	req.Handlers.Send.PushBack(func(r *request.Request) {
		data, header, err := c.GetObjectDirect(*in.Bucket, *in.Key)
		if err != nil {
			r.Error = err
			return
		}
		r.HTTPResponse = &http.Response{
			StatusCode: 200,
			Header:     header,
			Body:       ioutil.NopCloser(bytes.NewBuffer(data)),
		}
		out.Metadata = make(map[string]*string)
		out.Metadata["x-amz-wrap-alg"] = aws.String(s3crypto.KMSWrap)
	})
	return nil
}

func encodeBucketKey(bucket, key string) string {
	return bucket + "|" + key
}
