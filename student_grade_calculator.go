package main

import (
	"fmt"
)

func calculateGrade(score float64) string {
	if score >= 90 {
		return "A"
	} else if score >= 80 {
		return "B"
	} else if score >= 70 {
		return "C"
	} else if score >= 60 {
		return "D"
	} else {
		return "F"
	}
}

func calculateAvarage(scores []float64) float64 {
	if len(scores) == 0 {
		return 0
	}
	sum := 0.0
	for _, score := range scores {
		sum += score
	}
	return sum / float64(len(scores))
}

type Subject struct {
	Name  string
	Score float64
}

func inputSubAndScore(numSubjects int) ([]string, []float64) {
	var scores []float64
	var subjects []string
	for i := 0; i < numSubjects; i++ {
		var subject Subject
		fmt.Printf("Enter the Subject #%d name and score (e.g. Math 85): ", i+1)
		_, err := fmt.Scanf("%s %f", &subject.Name, &subject.Score)

		if err != nil || subject.Score < 0 || subject.Score > 100 {
			fmt.Println("Invalid input. Please enter a valid score between 0 and 100.")
			i-- // Decrement i to repeat this iteration
			continue
		}
		// check if the subject already exists
		exists := false
		for _, s := range subjects {
			if subject.Name == s {
				exists = true
				break
			}
		}
		if exists {
			fmt.Println("Subject already exists. Please enter a different subject name.")
			i-- // Decrement i to repeat this iteration
			continue
		}
		
		scores = append(scores, subject.Score)
		subjects = append(subjects, subject.Name)
		fmt.Printf("Subject: %s, Score: %.2f\n", subject.Name, subject.Score)
	}

	return subjects, scores
}

// function for inputing student name and number of subjects
func inputStudentName() string {
	var fname string
	var lname string
	for {
		fmt.Print("Enter student name (e.g. John Doe): ")
		_, err := fmt.Scanf("%s %s", &fname, &lname)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid name.")
			continue
		}
		break
	}
	return fname + " " + lname
}

// function for inputing number of subjects
func inputNumSubjects() int {
	var numSubjects int
	for {
		fmt.Print("Enter number of subjects: ")
		_, err := fmt.Scanf("%d", &numSubjects)
		if err != nil || numSubjects <= 0 {
			fmt.Println("Invalid input. Please enter a positive integer for the number of subjects.")
			continue
		}
		break
	}
	return numSubjects
}

func main() {
	var scores []float64
	var subjects []string

	var name string
	var numSubjects int

	fmt.Println("Welcome to the Student Grade Calculator!")
	name = inputStudentName()
	numSubjects = inputNumSubjects()

	subjects, scores = inputSubAndScore(numSubjects)

	average := calculateAvarage(scores)
	// calculate the individual grades'
	fmt.Println("Calculating individual grade...")
	for i, score := range scores {
		grade := calculateGrade(score)
		fmt.Printf("Subject: %s, Score: %.2f, Grade: %s\n", subjects[i], score, grade)
	}

	fmt.Printf("Average score for %s: %.2f\n", name, average)
	totalGrade := calculateGrade(average)
	fmt.Printf("Total average grade for %s: %s\n", name, totalGrade)
}
