package controls

// Model represents the data structure used for defining visibility, dimensions, and content.
type Model struct {
	IsVisible bool
	ViewWidth int
	Content   string
}

func New(body string) Model {
	return Model{
		IsVisible: true,
		ViewWidth: 80,
		Content:   body,
	}
}
