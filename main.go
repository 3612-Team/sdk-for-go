package appwrite

// NewClient initializes a new Appwrite client
func NewClient() Client {
	return Client{
		headers: make(map[string]string),
	}
}
