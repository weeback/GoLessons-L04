package Models

type Product struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Image       string   `json:"image"`
	Id          int      `json:"id"`
}
