package main

import (
	"fmt"
	"time"
)

func printFunction() {
	fmt.Println("hi")
}

func main() {
	// bby adding go we specify that we are executing the concurently
	//	go printFunction()
	//	go printFunction()
	//	go printFunction()

	//go routines fork of the main function, and they need to reenter through the join point of the main func
	//when executing concurent action we need to make sure that we await the function to finish execution.
	// time.Sleep is just an example, we await 2 seconds to make sure the goRoutine finished
	//time.Sleep(time.Second * 2)

	///CHANELS - FIFO QUEUES

	myChannel := make(chan string)
	//we fork off the main go routine
	go func() {
		myChannel <- "data"
	}()

	/// blocking line of code
	msg := <-myChannel

	fmt.Println(msg)
	//The main go routine is waiting for the channel to close, or a mesagge to be received from the channel
	// This is a good example of a join point

	//SELECT - main function waits for mesages from multiple channels

	anotherChannel := make(chan string)

	go func() {
		anotherChannel <- "something else"
	}()

	//will choose one case and main func wont be closed untill we get one of the messages below
	select {
	case messageFromMyChannel := <-myChannel:
		fmt.Println(messageFromMyChannel)
	case messageFromAnotherChannel := <-anotherChannel:
		fmt.Println(messageFromAnotherChannel)
	}

	charChannel := make(chan string, 3) //limited capacity of 3

	//an unbuffered channel provides a guarantee that an exchange of go routines is performed at the instance the send and receive take place

	chars := []string{"a", "b", "c"}

	for _, s := range chars {
		select {
		case charChannel <- s:
		}
	}

	close(charChannel)
	//even though the channel was closed, we can still read the data
	for result := range charChannel {
		fmt.Println(result)
	}

	//infinite running go routine, or untill we stop it
	go func() {
		for {
			select {
			default:
				fmt.Println("Doing something")
			}
		}
	}()

	//need to avoid leaks so it doesnt work in the background , using memory

	//our main function need to cancel the go routine

	done := make(chan bool)

	go doWork(done)

	time.Sleep(time.Second * 3)

	close(done)

	//pipelines
	nums := []int{2, 3, 4, 5, 6}

	//stage 1
	dataChannel := sliceToChannel(nums)
	//these are 2 unbuffered channels, so they run syncronously
	//stage2
	finalChannel := sq(dataChannel)

	//stage 3
	for n := range finalChannel {
		fmt.Println(n)
	}
}

func sliceToChannel(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}

		close(out)
	}()

	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()

	return out

}

// can olny read
func doWork(done <-chan bool) {
	for {
		select {
		//it will stop when it receives the done case, from the parent (main func)
		case <-done:
			return
		default:
			fmt.Println("Doing something")
		}
	}
}
