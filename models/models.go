package models

type SiteData struct {
	Meta       Meta            `json:"meta"`
	Pages      map[string]Page `json:"pages"` // Flexible page names
	Header     Page
	Components map[string]Page
}

type Meta struct {
	Title string `json:"title"`
}

type Page struct {
	Elements []PageElement `json:"elements"`
}

type PageElement struct {
	Type       string        `json:"type"` // The "type" for each element (e.g., "Div")
	Attributes Attributes    `json:"attributes"`
	Elements   []PageElement `json:"elements"` // Nested elements like "H1"
	Text       string        `json:"text"`     // Text content for elements like "H1"
	Link       string        `json:"link"`
	Import     string        `json:"import"`
}

type Message struct {
	Message string
	Type    string
}

/* type Attributes struct {
	ID      string                 `json:"id"`
	Class   string                 `json:"class"`
	OnClick string                 `json:"onclick"`
	Href    string                 `json:"href"`
	Src     string                 `json:"src"`
	Style   map[string]interface{} `json:"style"` // Flexible styling keys
} */

type Attributes struct {
	Props map[string]string      `json:"props"` // Flexible key-value pairs for any attribute
	Style map[string]interface{} `json:"style"` // For CSS styling properties
}

type Auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Auth0RegisterResponse struct {
	Email   string `json:"email"`
	Success bool   `json:"success"`
}
