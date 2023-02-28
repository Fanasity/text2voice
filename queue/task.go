package queue

type TaskHandler func()

type Task interface {
	Handle() error
}
