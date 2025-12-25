package types

// Section represents a collapsible information section
type Section struct {
	Name     string
	Expanded bool
	Data     map[string]string
	TreeData []TreeItem // Hierarchical data structure
	LiveData bool       // Whether this section supports live updates
	Order    []string   // Order of keys for display
	UseTree  bool       // Whether to use tree structure for display
}

// TreeItem represents a hierarchical data item
type TreeItem struct {
	Name     string
	Children map[string]string
	Order    []string
}
