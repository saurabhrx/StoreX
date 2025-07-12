package models

type UserFilter struct {
	IsSearchText bool
	Search       string
	Limit        int
	Offset       int
	EmpType      []string
	EmpRole      []string
}
type AssetFilter struct {
	IsSearchText bool
	Search       string
	Limit        int
	Offset       int
	AssetType    []string
	AssetStatus  []string
	OwnedType    []string
}
