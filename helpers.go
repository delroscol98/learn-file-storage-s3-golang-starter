package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getVideoAspectRatio(filepath string) (string, error) {
	type Out struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filepath)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var out Out
	err = json.Unmarshal(buf.Bytes(), &out)
	if err != nil {
		return "", err
	}

	width := out.Streams[0].Width
	height := out.Streams[0].Height
	ratio := float64(width) / float64(height)

	if ratio >= (16.5/9.5) && ratio <= (15.5/8.5) {
		return "16:9", nil
	}

	if ratio >= (8.5/15.5) && ratio <= (9.5/16.5) {
		return "9:16", nil
	}

	return "other", nil
}

func processVideoForFastStart(filepath string) (string, error) {
	outputFilePath := fmt.Sprintf("%s.processing", filepath)
	cmd := exec.Command("ffmpeg", "-i", filepath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFilePath)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFilePath, nil
}

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignedClient := s3.NewPresignClient(s3Client)
	params := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	presignedHTTPRequest, err := presignedClient.PresignGetObject(context.TODO(), &params, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}

	return presignedHTTPRequest.URL, nil
}
