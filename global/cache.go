package global

var cache = map[string]string{}

func Set(key, val string) {
	cache[key] = val
}


func Get(key string) string {
	return cache[key]
}
