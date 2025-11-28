package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}

	parts := strings.Split(*video.VideoURL, ",")
	bucket := parts[0]
	fileKey := parts[1]
	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, fileKey, 3600*time.Second)
	if err != nil {
		return video, err
	}

	fmt.Println(presignedURL)

	video.VideoURL = &presignedURL
	return video, nil
}
