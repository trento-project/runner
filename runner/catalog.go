package runner

const (
	CatalogDestinationFile = "ansible/catalog.json"
)

type Catalog []*CatalogCheck

type CatalogCheck struct {
	ID             string `json:"id,omitempty" binding:"required"`
	Name           string `json:"name,omitempty" binding:"required"`
	Group          string `json:"group" binding:"required"`
	Provider       string `json:"provider" binding:"required"`
	Description    string `json:"description,omitempty"`
	Remediation    string `json:"remediation,omitempty"`
	Implementation string `json:"implementation,omitempty"`
	Labels         string `json:"labels,omitempty"`
	Premium        bool   `json:"premium,omitempty"`
}
