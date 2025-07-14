package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"library_management/models"
	"library_management/services"
)

type LibraryController struct {
	service services.LibraryService
}

func NewLibraryController(service services.LibraryService) *LibraryController {
	return &LibraryController{service: service}
}

func (lc *LibraryController) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\n--- Library Management System ---")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books by Member")
		fmt.Println("0. Exit")
		fmt.Print("Enter your choice: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			lc.addBook(scanner)
		case "2":
			lc.removeBook(scanner)
		case "3":
			lc.borrowBook(scanner)
		case "4":
			lc.returnBook(scanner)
		case "5":
			lc.listAvailableBooks()
		case "6":
			lc.listBorrowedBooks(scanner)
		case "0":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func (lc *LibraryController) addBook(scanner *bufio.Scanner) {
	fmt.Print("Enter Book ID: ")
	scanner.Scan()
	id, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter Book Title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Enter Book Author: ")
	scanner.Scan()
	author := scanner.Text()

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	}
	lc.service.AddBook(book)
	fmt.Println("Book added successfully.")
}

func (lc *LibraryController) removeBook(scanner *bufio.Scanner) {
	fmt.Print("Enter Book ID to remove: ")
	scanner.Scan()
	id, _ := strconv.Atoi(scanner.Text())
	lc.service.RemoveBook(id)
}

func (lc *LibraryController) borrowBook(scanner *bufio.Scanner) {
	fmt.Print("Enter Book ID to borrow: ")
	scanner.Scan()
	bookID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter Member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())

	err := lc.service.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed successfully.")
	}
}

func (lc *LibraryController) returnBook(scanner *bufio.Scanner) {
	fmt.Print("Enter Book ID to return: ")
	scanner.Scan()
	bookID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter Member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())

	err := lc.service.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned successfully.")
	}
}

func (lc *LibraryController) listAvailableBooks() {
	books := lc.service.ListAvailableBooks()
	fmt.Println("\nAvailable Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func (lc *LibraryController) listBorrowedBooks(scanner *bufio.Scanner) {
	fmt.Print("Enter Member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())
	books := lc.service.ListBorrowedBooks(memberID)
	fmt.Printf("\nBooks borrowed by Member %d:\n", memberID)
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}
