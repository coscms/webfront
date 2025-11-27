package i18nm

type Model interface {
	Short_() string
	GetField(string) interface{}
	FromRow(map[string]interface{})
}
