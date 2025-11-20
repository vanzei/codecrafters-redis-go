package builtin

import "time"

type blockRequest struct {
	clientID string
	lists    []string
	result   chan string
	timeout  *time.Timer
}

var waiters = map[string][]*blockRequest{}

func addWaiter(req *blockRequest) {
	for _, list := range req.lists {
		waiters[list] = append(waiters[list], req)
	}
}

func removeWaiter(req *blockRequest) {
	if req.timeout != nil {
		req.timeout.Stop()
	}
	for _, list := range req.lists {
		queue := waiters[list]
		for i, w := range queue {
			if w == req {
				waiters[list] = append(queue[:i], queue[i+1:]...)
				break
			}
		}
		if len(waiters[list]) == 0 {
			delete(waiters, list)
		}

	}
}
