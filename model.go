package distance

//Matrix response
type Matrix struct {
	Res []ResourceSet `json:"resourceSets"`
}

//ResourceSet model
type ResourceSet struct {
	Res []map[string]interface{} `json:"resources"`
}
