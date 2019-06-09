package network

import "context"

type messageQueue struct {
	context context.Context
	cancel  context.CancelFunc
	action  chan func()
}

func newMessageQueue() *messageQueue {
	ctx, cancel := context.WithCancel(context.Background())

	messageQueue := &messageQueue{
		context: ctx,
		cancel:  cancel,
		action:  make(chan func(), 0),
	}

	go func() {
		for {
			select {
			case <-messageQueue.context.Done():
				return
			case action, ok := <-messageQueue.action:
				if ok {
					action()
				}
			}
		}
	}()

	return messageQueue
}

func (messageQueue *messageQueue) GetWriteChan() chan<- func() {
	return messageQueue.action
}

func (messageQueue *messageQueue) Close() {
	messageQueue.cancel()
}
