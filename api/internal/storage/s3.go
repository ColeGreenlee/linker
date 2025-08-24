package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	
	"linker/internal/config"
	"linker/internal/utils"
)

type S3Client struct {
	session    *session.Session
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	config     *config.S3Config
}

type UploadResult struct {
	Key      string
	URL      string
	Size     int64
	MimeType string
}

func NewS3Client(cfg *config.S3Config) (*S3Client, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("S3 storage is not enabled")
	}

	// Configure AWS session
	awsConfig := &aws.Config{
		Region:           aws.String(cfg.Region),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		DisableSSL:       aws.Bool(!cfg.UseSSL),
		S3ForcePathStyle: aws.Bool(true), // Required for MinIO
	}

	if cfg.Endpoint != "" {
		awsConfig.Endpoint = aws.String(fmt.Sprintf("http%s://%s", 
			func() string { if cfg.UseSSL { return "s" }; return "" }(), 
			cfg.Endpoint))
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	s3Client := s3.New(sess)
	
	// Test connection and create bucket if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err = s3Client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.BucketName),
	})
	if err != nil {
		// Try to create bucket
		_, createErr := s3Client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.BucketName),
		})
		if createErr != nil {
			return nil, fmt.Errorf("bucket %s doesn't exist and failed to create: %w", cfg.BucketName, createErr)
		}
	}

	return &S3Client{
		session:    sess,
		s3Client:   s3Client,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		config:     cfg,
	}, nil
}

func (s *S3Client) Upload(ctx context.Context, filename string, content io.Reader, mimeType string) (*UploadResult, error) {
	// Generate unique S3 key
	s3Key := s.generateS3Key(filename)
	
	// Read content to get size
	contentBytes, err := io.ReadAll(content)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}
	
	contentReader := bytes.NewReader(contentBytes)
	size := int64(len(contentBytes))
	
	// Check file size limit
	if s.config.MaxFileSize > 0 && size > s.config.MaxFileSize*1024*1024 {
		return nil, fmt.Errorf("file size %d bytes exceeds limit of %d MB", 
			size, s.config.MaxFileSize)
	}
	
	// Check MIME type
	if !s.isMimeTypeAllowed(mimeType) {
		return nil, fmt.Errorf("MIME type %s is not allowed", mimeType)
	}
	
	// Upload to S3
	uploadInput := &s3manager.UploadInput{
		Bucket:      aws.String(s.config.BucketName),
		Key:         aws.String(s3Key),
		Body:        contentReader,
		ContentType: aws.String(mimeType),
		Metadata: map[string]*string{
			"original-filename": aws.String(filename),
			"upload-time":       aws.String(time.Now().Format(time.RFC3339)),
		},
	}
	
	result, err := s.uploader.UploadWithContext(ctx, uploadInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	
	return &UploadResult{
		Key:      s3Key,
		URL:      result.Location,
		Size:     size,
		MimeType: mimeType,
	}, nil
}

func (s *S3Client) Download(ctx context.Context, s3Key string) (io.ReadCloser, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(s3Key),
	}
	
	result, err := s.s3Client.GetObjectWithContext(ctx, getObjectInput)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	
	return result.Body, nil
}

func (s *S3Client) GetDownloadURL(ctx context.Context, s3Key string, duration time.Duration) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(s3Key),
	})
	
	url, err := req.Presign(duration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	
	return url, nil
}

func (s *S3Client) Delete(ctx context.Context, s3Key string) error {
	_, err := s.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

func (s *S3Client) GetFileInfo(ctx context.Context, s3Key string) (*s3.HeadObjectOutput, error) {
	result, err := s.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	
	return result, nil
}

func (s *S3Client) generateS3Key(filename string) string {
	// Generate a unique key with timestamp and UUID
	uuid := utils.GenerateUUID()
	timestamp := time.Now().Format("2006/01/02")
	ext := filepath.Ext(filename)
	
	return fmt.Sprintf("%s/%s%s", timestamp, uuid, ext)
}

func (s *S3Client) isMimeTypeAllowed(mimeType string) bool {
	if len(s.config.AllowedMimeTypes) == 0 {
		return true // No restrictions
	}
	
	for _, allowed := range s.config.AllowedMimeTypes {
		if allowed == mimeType {
			return true
		}
		// Support wildcard matching (e.g., "image/*")
		if strings.HasSuffix(allowed, "/*") {
			prefix := strings.TrimSuffix(allowed, "/*")
			if strings.HasPrefix(mimeType, prefix+"/") {
				return true
			}
		}
	}
	
	return false
}

// Helper functions for file operations
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func GetMimeTypeFromExtension(ext string) string {
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".zip":  "application/zip",
		".json": "application/json",
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
	}
	
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	
	return "application/octet-stream"
}

func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}