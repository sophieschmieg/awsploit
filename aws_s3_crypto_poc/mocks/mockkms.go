package mocks

import (
	"crypto/rand"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
)

type mockAWSKMSClient struct {
	kmsiface.KMSAPI
	aead tink.AEAD
}

func (c mockAWSKMSClient) DecryptWithContext(ctx aws.Context, input *kms.DecryptInput, opt ...request.Option) (*kms.DecryptOutput, error) {
	d, err := c.aead.Decrypt(input.CiphertextBlob, nil)
	if err != nil {
		return nil, err
	}
	return &kms.DecryptOutput{Plaintext: d}, nil
}

func (c mockAWSKMSClient) GenerateDataKeyWithContext(ctx aws.Context, in *kms.GenerateDataKeyInput, opt ...request.Option) (*kms.GenerateDataKeyOutput, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	ci, err := c.aead.Encrypt(key, nil)
	if err != nil {
		return nil, err
	}
	return &kms.GenerateDataKeyOutput{
		Plaintext:      key,
		CiphertextBlob: ci,
	}, nil
}

func NewMockAWSKMSClient(sess *session.Session) (kmsiface.KMSAPI, error) {
	h, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		return nil, err
	}
	aead, err := aead.New(h)
	if err != nil {
		return nil, err
	}
	return &mockAWSKMSClient{KMSAPI: kms.New(sess), aead: aead}, nil
}
