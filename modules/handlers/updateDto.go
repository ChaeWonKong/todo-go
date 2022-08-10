package handlers

type UpdateDto struct {
	Title   string `json:"title" validate:"omitempty"`
	Checked *bool  `json:"checked" validate:"omitempty"`
}
