package channels

func OK(channel chan bool) bool {
	select {
	case ok := <-channel:
		if ok {
			return ok
		}
	}
	return false
}
