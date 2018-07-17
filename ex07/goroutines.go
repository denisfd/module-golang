package goroutines

func close_channel(wait_flag chan byte, channel chan string) {
	<-wait_flag
	close(channel)
}

func append_brackets(finished_flag chan byte, input chan string, output chan string) {
	for str := range input {
		str = "(" + str + ")"
		output <- str
	}
	finished_flag <- 1
}

func Process(input chan string) chan string {
	output := make(chan string)
	wait_flag := make(chan byte)

	go close_channel(wait_flag, output)
	go append_brackets(wait_flag, input, output)

	return output
}
