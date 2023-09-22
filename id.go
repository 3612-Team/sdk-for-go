package appwrite

type ID struct{}

func (i ID) Custom(id string) string {
	return id
}

func (i ID) Unique() string {
	return "unique()"
}
