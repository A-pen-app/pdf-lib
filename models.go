package pdf

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
	OutputType string      `json:"outputType"`
	Template   string      `json:"template"`
	Format     string      `json:"format"`
	Data       interface{} `json:"data"`
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
