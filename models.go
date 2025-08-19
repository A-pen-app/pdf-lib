package pdf

// OutputType represents the output type for generation
type OutputType string

const (
	OutputTypePDF   OutputType = "pdf"
	OutputTypeImage OutputType = "image"
)

// Format represents the format for image output
type Format string

const (
	FormatPNG Format = "png"
)

// ResumeTemplate represents the template type for PDF generation
type ResumeTemplate string

const (
	TemplateResumeApen  ResumeTemplate = "apen_resume"
	TemplateResumePhar  ResumeTemplate = "phar_resume"
	TemplateResumeNurse ResumeTemplate = "nurse_resume"
)

type SharePostTemplate string

const (
	TemplateSharePostApen  SharePostTemplate = "apen_post_sharing"
	TemplateSharePostPhar  SharePostTemplate = "phar_post_sharing"
	TemplateSharePostNurse SharePostTemplate = "nurse_post_sharing"
)

type GenerateRequest struct {
	OutputType OutputType  `json:"outputType"`
	Template   string      `json:"template"`
	Format     Format      `json:"format"`
	Data       interface{} `json:"data"`
}

// GetExtensionAndMimeType returns the file extension and MIME type based on output type and format
func (r *GenerateRequest) GetExtensionAndMimeType() (string, string) {
	switch r.OutputType {
	case OutputTypePDF:
		return ".pdf", "application/pdf"
	case OutputTypeImage:
		switch r.Format {
		case FormatPNG:
			return ".png", "image/png"
		}
	}
	return "", ""
}

// SharePostImageData represents the specific data fields for share post image generation
type SharePostData struct {
	Gender         string `json:"gender"`
	Username       string `json:"username"`
	Picture        string `json:"picture"`
	Position       string `json:"position"`
	Specialty      string `json:"specialty"`
	SpecialtyBadge string `json:"specialty_badge"`
	Category       string `json:"category"`
	Title          string `json:"title"`
	Content        string `json:"content"`
}

type GenerateResult struct {
	Data      []byte
	MimeType  string
	Extension string
}
