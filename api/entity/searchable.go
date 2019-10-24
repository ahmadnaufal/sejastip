package entity

type DynamicFilter map[string]string

// Searchable is an interface that enables an entity columns in a repository to be searched
type Searchable interface {
	GetFilterQuery(filter DynamicFilter) string
}
