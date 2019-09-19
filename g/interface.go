package g

type Metrics interface {
	Collect() error
}
