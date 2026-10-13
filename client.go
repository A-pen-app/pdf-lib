package pdf

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/A-pen-app/logging"
)

type Client interface {
	GenerateResumePDF(ctx context.Context, template ResumeTemplate, data map[string]interface{}) (*GenerateResult, error)
	GenerateShareImage(ctx context.Context, template SharePostTemplate, data SharePostData) (*GenerateResult, error)
	VideoThumbnail(ctx context.Context, videoURL string) (*GenerateResult, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) Client {
	return &client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *client) GenerateResumePDF(ctx context.Context, template ResumeTemplate, data map[string]interface{}) (*GenerateResult, error) {
	requestData := GenerateRequest{
		OutputType: OutputTypePDF,
		Template:   string(template),
		Data:       data,
	}

	return c.generate(ctx, requestData)
}

func (c *client) GenerateShareImage(ctx context.Context, template SharePostTemplate, data SharePostData) (*GenerateResult, error) {
	requestData := GenerateRequest{
		OutputType: OutputTypeImage,
		Template:   string(template),
		Format:     FormatPNG,
		Data:       data,
	}

	return c.generate(ctx, requestData)
}

// VideoThumbnail returns the first frame of the video at videoURL as a JPEG;
// pdf-service fetches videoURL itself.
func (c *client) VideoThumbnail(ctx context.Context, videoURL string) (*GenerateResult, error) {
	if c.baseURL == "" {
		return nil, errors.New("PDF service unavailable")
	}

	jsonData, err := json.Marshal(VideoThumbnailRequest{URL: videoURL})
	if err != nil {
		return nil, fmt.Errorf("invalid request data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/video-thumbnail", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logging.Errorw(ctx, "call pdf service failed", "err", err)
		return nil, fmt.Errorf("PDF service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logging.Errorw(ctx, "pdf service returned error", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("PDF service error (status %d): %s", resp.StatusCode, string(body))
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("video thumbnail failed: %w", err)
	}

	return &GenerateResult{
		Data:      result,
		MimeType:  "image/jpeg",
		Extension: ".jpg",
	}, nil
}

func (c *client) generate(ctx context.Context, request GenerateRequest) (*GenerateResult, error) {
	if c.baseURL == "" {
		return nil, errors.New("PDF service unavailable")
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		logging.Errorw(ctx, "marshal pdf data failed", "err", err)
		return nil, fmt.Errorf("invalid request data: %v", err)
	}

	logging.Infow(ctx, "sending data to PDF service", "jsonData", string(jsonData))

	resp, err := http.Post(
		c.baseURL+"/api/v1/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		logging.Errorw(ctx, "call pdf service failed", "err", err)
		return nil, fmt.Errorf("PDF service unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logging.Errorw(ctx, "pdf service returned error", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("PDF service error (status %d): %s", resp.StatusCode, string(body))
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Errorw(ctx, "read pdf response failed", "err", err)
		return nil, fmt.Errorf("PDF generation failed: failed to read response")
	}

	extension, mimeType := request.GetExtensionAndMimeType()

	return &GenerateResult{
		Data:      result,
		MimeType:  mimeType,
		Extension: extension,
	}, nil
}
