package services

import (
	"fmt"
	"library_management/models"
)

type LibraryService interface {
	AddBook(book models.Book)
    RemoveBook(bookID int)
    BorrowBook(bookID int, memberID int) error
    ReturnBook(bookID int, memberID int) error
    ListAvailableBooks() []models.Book
    ListBorrowedBooks(memberID int) []models.Book
}

type libraryService struct {
    books       []models.Book
    members     []models.Member
	borrowedBooks map[int][]int // Maps memberID to a list of borrowed book IDs
}

func NewLibraryService() LibraryService {
	return &libraryService{
		books:       []models.Book{},
		members:     []models.Member{},
		borrowedBooks: make(map[int][]int),
	}
}

// Implement LibraryService interface methods
func (ls *libraryService) AddBook(book models.Book) {
	// Check if the book already exists
	for _, b := range ls.books {
		if b.ID == book.ID {
			fmt.Printf("Book with ID %d already exists.\n", book.ID)
			return
		}
	}
	// Add the book to the library
	if book.Status == "" {
		book.Status = "Available"
	}

	ls.books = append(ls.books, book)
}

func (ls *libraryService) RemoveBook(bookID int) {
	if len(ls.books) == 0 {
		fmt.Println("No books available to remove.")
		return
	}

	found := false
	for i, book := range ls.books {
		if book.ID == bookID {
			fmt.Printf("Removing book with ID %d.\n", bookID)
			ls.books = append(ls.books[:i], ls.books[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Book with ID %d not found.\n", bookID)
	}
}


func (ls *libraryService) BorrowBook(bookID int, memberID int) error {
	bookIdx := -1
	for i, book := range ls.books {
	if bookID == book.ID {
		bookIdx = i
		break
	}
	}
	
	if bookIdx == -1 {
		return fmt.Errorf("book with id %d not found", bookIdx)
	}

	if ls.books[bookIdx].Status == "Borrowed" {
		return fmt.Errorf("book with id %d is currently borrowed", bookIdx)
	}

	memberIdx := -1
	for i, member := range ls.members {
		if member.ID == memberID {
		memberIdx = i
		break
		}
	}
	if memberIdx == -1 {
		return fmt.Errorf("member with id %d not found", memberID)
	}
	ls.books[bookIdx].Status = "Borrowed" 
	ls.members[memberIdx].BorrowedBooks = append(ls.members[memberIdx].BorrowedBooks, ls.books[bookIdx])
	ls.borrowedBooks[memberID] = append(ls.borrowedBooks[memberID], bookID)
	return nil
}

func (ls *libraryService) ReturnBook(bookID int, memberID int) error {
	borrowed, ok := ls.borrowedBooks[memberID]
	if !ok {
		return fmt.Errorf("member with id %d has not borrowed any books", memberID)
	}
	
	for i, id := range borrowed {
		if id == bookID {
			ls.borrowedBooks[memberID] = append(borrowed[:i], borrowed[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("book with id %d not found in member's borrowed books", bookID)
}

func (ls *libraryService) ListAvailableBooks() []models.Book {
	borrowed := make(map[int]bool)
	for _, books := range ls.borrowedBooks {
		for _, bookID := range books {
			borrowed[bookID] = true
		}
	}
	var availableBooks []models.Book
	for _, book := range ls.books {
		if !borrowed[book.ID] {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (ls *libraryService) ListBorrowedBooks(memberID int) []models.Book {

	var borrowedBooks []models.Book
	bookIDs := ls.borrowedBooks[memberID]
	for _, id := range bookIDs {
		for _, book := range ls.books {
			if book.ID == id {
				borrowedBooks = append(borrowedBooks, book)
			}
		}
	}
	return borrowedBooks
}

