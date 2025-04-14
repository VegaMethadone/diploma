package notebook

type BlockInterface interface {
	GetType() string
	GetId() int
	GetBody() any
	GetComments()
}
