package helper

import (
	"errors"
	"fmt"
	"mime/multipart"
	"monitoring-service/dto"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// ValidateUUID checks if the given string is a valid UUID
func ValidateUUID(id string) bool {
	if id == "" {
		return false
	}
	return uuidRegex.MatchString(strings.ToLower(id))
}

// SanitizeString removes potentially dangerous SQL characters
func SanitizeString(input string) string {
	if input == "" {
		return input
	}

	// Remove dangerous SQL characters
	dangerous := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_"}
	result := input

	for _, d := range dangerous {
		result = strings.ReplaceAll(result, d, "")
	}

	// Remove SQL keywords (case insensitive)
	sqlKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"UNION", "EXEC", "EXECUTE", "DECLARE", "CAST", "CONVERT",
	}

	upperResult := strings.ToUpper(result)
	for _, keyword := range sqlKeywords {
		upperResult = strings.ReplaceAll(upperResult, keyword, "")
	}

	return result
}

func ValidateReportType(reportType string) bool {
	validTypes := []string{"WEEKLY_REPORT", "FINAL_REPORT"}
	for _, vt := range validTypes {
		if reportType == vt {
			return true
		}
	}
	return false
}

// SANITIZE REPORT
const (
	MaxFileSize      = 10 * 1024 * 1024 // 10MB
	MaxContentLength = 50 * 1024        // 50KB for text content
)

// Allowed file extensions and MIME types
var allowedExtensions = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".txt":  true,
	".rtf":  true,
}

var allowedMimeTypes = map[string]bool{
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/plain":      true,
	"text/rtf":        true,
	"application/rtf": true,
}

// ValidateReportRequest validates the report request structure
func ValidateReportRequest(req dto.ReportRequest) error {
	// Validate required fields
	if req.ReportScheduleID == "" {
		return errors.New("report schedule ID is required")
	}

	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.ReportType == "" {
		return errors.New("report type is required")
	}

	// Validate UUID format for report_schedule_id
	if !ValidateUUID(req.ReportScheduleID) {
		return errors.New("invalid report schedule ID format")
	}

	// Validate report type
	if !ValidateReportType(req.ReportType) {
		return errors.New("invalid report type")
	}

	// Validate title length and content
	if len(req.Title) > 255 {
		return errors.New("title too long (max 255 characters)")
	}

	if len(req.Content) > MaxContentLength {
		return errors.New("content too long (max 50KB)")
	}

	// Sanitize and validate title
	sanitizedTitle := SanitizeString(req.Title)
	if sanitizedTitle != req.Title {
		return errors.New("title contains invalid characters")
	}

	return nil
}

// ValidateFileUpload validates uploaded file
func ValidateFileUpload(file *multipart.FileHeader) error {
	if file == nil {
		return nil // File is optional
	}

	// Check file size
	if file.Size > MaxFileSize {
		return fmt.Errorf("file too large (max %d MB)", MaxFileSize/(1024*1024))
	}

	if file.Size == 0 {
		return errors.New("file is empty")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("file type not allowed (allowed: %v)", getStringKeys(allowedExtensions))
	}

	// Validate filename
	if len(file.Filename) > 255 {
		return errors.New("filename too long")
	}

	// Check for dangerous filenames
	filename := strings.ToLower(file.Filename)
	dangerousNames := []string{"web.config", ".htaccess", "autorun.inf"}
	for _, dangerous := range dangerousNames {
		if strings.Contains(filename, dangerous) {
			return errors.New("filename not allowed")
		}
	}

	return nil
}

// ValidateMimeType validates file MIME type by reading file header
func ValidateMimeType(file multipart.File) error {
	// Read first 512 bytes to determine MIME type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return errors.New("unable to read file")
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Detect MIME type
	mimeType := http.DetectContentType(buffer)

	// Check against allowed MIME types
	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("file type not allowed (detected: %s)", mimeType)
	}

	return nil
}

// Helper function to get string keys from map
func getStringKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func IsValidTokenFormat(token string) bool {
	// Adjust this based on your token format (JWT, Bearer, etc.)
	if strings.HasPrefix(token, "Bearer ") {
		tokenValue := strings.TrimPrefix(token, "Bearer ")
		// Basic length check for JWT-like tokens
		return len(tokenValue) > 10 && !strings.Contains(tokenValue, " ")
	}

	// For other token formats, add appropriate validation
	return len(token) > 10 && !strings.Contains(token, " ")
}

func ValidatePaginationParams(ctx *gin.Context) (int, int, error) {
	page := 1
	limit := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil || p < 1 {
			return 0, 0, errors.New("invalid page parameter")
		} else {
			page = p
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil || l < 1 || l > 100 {
			return 0, 0, errors.New("invalid limit parameter (max 100)")
		} else {
			limit = l
		}
	}

	return page, limit, nil
}
