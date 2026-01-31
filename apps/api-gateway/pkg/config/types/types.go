package types

type Client interface {
	GetStringSlice(key string, defaultValues []string) []string
	GetInt(key string, defaultValue int) int
	GetFloat64(key string, defaultValue float64) float64
	GetJSONPropertyAndUnmarshal(key string, structType interface{}) error
}
