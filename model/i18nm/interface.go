package i18nm

// Model is an interface for model instances.
type Model interface {
	Short_() string
	GetField(string) interface{}
	FromRow(map[string]interface{})
}

// GetRowID returns the ID of a model instance.
func GetRowID(mdl Model) uint64 {
	var id uint64
	switch v := mdl.GetField(`Id`).(type) {
	case uint:
		id = uint64(v)
	case uint8:
		id = uint64(v)
	case uint16:
		id = uint64(v)
	case uint32:
		id = uint64(v)
	case uint64:
		id = v
	default:
	}
	return id
}
