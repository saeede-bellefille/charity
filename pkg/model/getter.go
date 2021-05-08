package model

type Getter interface {
	Get(key string) string
}

type MapGetter map[string]string

func (mg MapGetter) Get(key string) string {
	return map[string]string(mg)[key]
}
