package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {
    log.SetPrefix("greetings: ")
    log.SetFlags(0)

    names := []string{"Gladys", "Samantha", "Darrin"}

    // Get a greeting message and print it.
    messages, err := greetings.Hellos(names)
    
    if err != nil {
        log.Fatal(err)
    }
    
    // If no error was returned, print the returned map of
    // messages to the console.
    fmt.Println(messages)
}