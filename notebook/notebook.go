package notebook

type Block interface {
	GetType() int64
	GetId() string
}
