package out

type ResourceResponse struct {
	ResourceID  uint   `json:"resource_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
