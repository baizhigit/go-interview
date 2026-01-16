package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type NumberPattern struct {
	Number  string
	Pattern string
	Score   int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	numbers := readNumbers(filename)
	vipNumbers := findVIPNumbers(numbers)

	// Sort by score (highest first)
	sort.Slice(vipNumbers, func(i, j int) bool {
		return vipNumbers[i].Score > vipNumbers[j].Score
	})

	// Write results to file
	outputFile := "vip_numbers.txt"
	writeResults(vipNumbers, outputFile)

	// Also print to console
	fmt.Printf("Found %d VIP numbers:\n\n", len(vipNumbers))
	for _, vip := range vipNumbers {
		fmt.Printf("%s - %s (Score: %d)\n", vip.Number, vip.Pattern, vip.Score)
	}
	fmt.Printf("\nResults written to %s\n", outputFile)
}

func readNumbers(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var numbers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && len(line) >= 10 {
			numbers = append(numbers, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	return numbers
}

func writeResults(vipNumbers []NumberPattern, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write header
	fmt.Fprintf(writer, "VIP Phone Numbers - Found %d patterns\n", len(vipNumbers))
	fmt.Fprintf(writer, "="+strings.Repeat("=", 60)+"\n\n")

	// Write each VIP number
	for _, vip := range vipNumbers {
		fmt.Fprintf(writer, "%s - %s (Score: %d)\n", vip.Number, vip.Pattern, vip.Score)
	}

	writer.Flush()
}

func findVIPNumbers(numbers []string) []NumberPattern {
	var vipNumbers []NumberPattern

	for _, number := range numbers {
		patterns := analyzeNumber(number)
		if len(patterns) > 0 {
			// Take the best pattern for this number
			bestPattern := patterns[0]
			for _, p := range patterns {
				if p.Score > bestPattern.Score {
					bestPattern = p
				}
			}
			vipNumbers = append(vipNumbers, bestPattern)
		}
	}

	return vipNumbers
}

func analyzeNumber(number string) []NumberPattern {
	var patterns []NumberPattern

	if len(number) < 10 {
		return patterns
	}

	// Get last 7 digits (the part after 706)
	lastDigits := number[len(number)-7:]

	// Easy to remember 3-digit patterns in format XXX-YY-ZZ
	// Split last 7 digits into 3-2-2 format
	if len(lastDigits) == 7 {
		part1 := lastDigits[0:3] // First 3 digits
		part2 := lastDigits[3:5] // Next 2 digits
		part3 := lastDigits[5:7] // Last 2 digits

		score, pattern := analyzeEasyPattern(part1, part2, part3)
		if score > 0 {
			patterns = append(patterns, NumberPattern{number, pattern, score})
		}
	}

	return patterns
}

func analyzeEasyPattern(p1, p2, p3 string) (int, string) {
	// Check for various easy-to-remember patterns

	// Special pattern: XXX has two zeros, YY has one zero, ZZ has one zero
	// Examples: 200-10-01, 100-20-05, 300-05-10
	zeros1 := countZeros(p1)
	zeros2 := countZeros(p2)
	zeros3 := countZeros(p3)

	if zeros1 == 2 && zeros2 == 1 && zeros3 == 1 {
		return 100, fmt.Sprintf("Two-One-One Zeros: %s-%s-%s", p1, p2, p3)
	}

	// All same digits in part1 (e.g., 111, 222, 333)
	if p1[0] == p1[1] && p1[1] == p1[2] {
		return 95, fmt.Sprintf("Triple Start: %s-%s-%s", p1, p2, p3)
	}

	// All same in part2 (e.g., XX-00-XX, XX-11-XX)
	if p2[0] == p2[1] {
		if p3[0] == p3[1] {
			// Both part2 and part3 are doubles
			return 90, fmt.Sprintf("Double Pairs: %s-%s-%s", p1, p2, p3)
		}
		return 85, fmt.Sprintf("Double Middle: %s-%s-%s", p1, p2, p3)
	}

	// All same in part3 (e.g., XX-YY-00, XX-YY-11)
	if p3[0] == p3[1] {
		return 85, fmt.Sprintf("Double End: %s-%s-%s", p1, p2, p3)
	}

	// Sequential in part1 (e.g., 123, 234, 345)
	if hasSequentialAscending(p1) {
		return 80, fmt.Sprintf("Sequential Start: %s-%s-%s", p1, p2, p3)
	}

	// Palindrome in all 3 parts together (e.g., 100-20-01 reads as 1002001)
	fullStr := p1 + p2 + p3
	if isPalindrome(fullStr) {
		return 75, fmt.Sprintf("Palindrome: %s-%s-%s", p1, p2, p3)
	}

	// Round numbers (ends with 00)
	if p3 == "00" {
		return 70, fmt.Sprintf("Round End: %s-%s-%s", p1, p2, p3)
	}

	// Repeated digit pattern (e.g., 121, 131, 141)
	if p1[0] == p1[2] && p1[0] != p1[1] {
		return 65, fmt.Sprintf("ABA Pattern: %s-%s-%s", p1, p2, p3)
	}

	// Sequential pairs (e.g., 12-23-34 or similar)
	if isSequentialPair(p1, p2) || isSequentialPair(p2, p3) {
		return 60, fmt.Sprintf("Sequential Pairs: %s-%s-%s", p1, p2, p3)
	}

	return 0, ""
}

func isPalindrome(s string) bool {
	if len(s) < 3 {
		return false
	}
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

func isSequentialPair(s1, s2 string) bool {
	if len(s1) < 1 || len(s2) < 1 {
		return false
	}
	// Check if last digit of s1 + 1 == first digit of s2
	return int(s1[len(s1)-1])+1 == int(s2[0])
}

func countZeros(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			count++
		}
	}
	return count
}

func hasSequentialAscending(s string) bool {
	count := 0
	for i := 0; i < len(s)-1; i++ {
		diff := int(s[i+1]) - int(s[i])
		if diff == 1 {
			count++
			if count >= 2 {
				return true
			}
		} else {
			count = 0
		}
	}
	return false
}
