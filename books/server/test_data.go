package server

import (
	booksTutorial "tutorial/gen/go/proto/books"
)

var (
	testBook1 = &booksTutorial.Book{
		Author: "Sam Newman",
		Title:  "Building microservices",
		Isbn:   "978-1491950357",
	}
	testBook2 = &booksTutorial.Book{
		Author: "Mat Ryer",
		Title:  "Building web applications",
		Isbn:   "978-1787123496",
	}
)
