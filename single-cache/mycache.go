package single_cache

type Getter interface {
	Ger(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

// Get实现Getter接口函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
