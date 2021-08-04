package main

import (
	"adventure"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 3000, "port for the adventure webapp")
	filename := flag.String("file", "gopher.json", "the file storing our adventure")
	path := strings.Join([]string{"resources", *filename}, string(os.PathSeparator))
	flag.Parse()

	fmt.Printf("Using the story file at %s\n", path)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	story, err := adventure.JsonStory(f)
	if err != nil {
		panic(err)
	}

	// pass custom template
	//t := template.Must(template.New("default").Parse("bla"))
	//h := adventure.NewHandler(story, adventure.WithTemplate(t))
	h := adventure.NewHandler(story)
	fmt.Printf("Starting the server on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
