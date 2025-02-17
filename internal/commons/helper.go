package commons

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-scratch/internal/commons"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
)

var location *time.Location

func GenerateCustomUID() string {
	currentTime := time.Now()
	return currentTime.Format("20060102150405")
}

func GetLocalTime() time.Time {
	location, _ := time.LoadLocation("Asia/Makassar")
	return time.Now().In(location)
}

func GenerateQRCode(text string) (string, error) {
	// Generate QR code
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Get PNG bytes
	png, err := qr.PNG(256)
	if err != nil {
		return "", fmt.Errorf("failed to encode QR code as PNG: %w", err)
	}

	// Convert to base64
	base64Str := base64.StdEncoding.EncodeToString(png)

	return base64Str, nil
}

func VerifyQRCode(ctx context.Context, key string, c redis.UniversalClient) (bool, error) {
	redisKey := fmt.Sprintf("%s%s", commons.QR_KEY, key)
	res, err := c.Get(ctx, redisKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}

	type QRCode struct {
		Secret string
		Exp    time.Time
	}

	var qrCode QRCode
	err = json.Unmarshal([]byte(res), &qrCode)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal QR code: %w", err)
	}

	if qrCode.Secret != key {
		return false, nil
	}

	if qrCode.Exp.Before(time.Now()) {
		return false, nil
	}

	return true, nil

}

// SaveQRCodeToFile generates a QR code and saves it to a file
func SaveQRCodeToFile(text, outputPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate and save QR code
	err := qrcode.WriteFile(text, qrcode.Medium, 256, outputPath)
	if err != nil {
		return fmt.Errorf("failed to save QR code: %w", err)
	}

	return nil
}
