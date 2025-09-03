package etcd

// OptionFunc ...
type OptionFunc func(storage *ValuesStorage)

// WithPath ...
func WithPath(path string) OptionFunc {
	return func(storage *ValuesStorage) {
		storage.path = path
	}
}
