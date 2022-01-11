init:
	go mod tidy
run_publisher:
	go run publisher/publisher.go
run_subscriber_all:
	go run subscriber/subscriber.go "#"
run_subscriber_fatal_all:
	go run subscriber/subscriber.go "*.fatal"
run_subscriber_warn_kernel:
	go run subscriber/subscriber.go "kern.warning"