package workpool

type Worker struct {
	pool *Pool

	task chan f
}
