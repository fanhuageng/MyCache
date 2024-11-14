package single_cache

// ByteView保存字节的不可变视图。
type ByteView struct {
	b []byte
}

func (bv ByteView) Len() int {
	return len(bv.b)
}

func (bv ByteView) String() string {
	return string(bv.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}
