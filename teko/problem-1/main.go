package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// OnceHour is the number of seconds in an hour
const OnceHour = 3600

// enqueue adds an element to the end of the queue
func enqueue(queue []int64, element int64) []int64 {
	queue = append(queue, element)
	return queue
}

// dequeue removes the first element from the queue
func dequeue(queue []int64) []int64 {
	return queue[1:]
}

// RateLimiteCore is the core of the rate limiter
type RateLimiteCore struct {
	// rate limit (hourly)
	rateLimit  int
	requestLog []int64
}

// SetRateLimit sets the rate limit
func (r *RateLimiteCore) SetRateLimit(rateLimit int) {
	r.rateLimit = rateLimit
}

// StringToUnix converts a string to a Unix timestamp
func (r *RateLimiteCore) StringToUnix(requestedAt string) (int64, error) {
	t, err := time.Parse(time.RFC3339, requestedAt)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// CheckRateLimit checks if the number of requests is less than the rate limit
func (r *RateLimiteCore) CheckRateLimit(requestedAt int64) bool {
	// remove requests older than an hour
	for _, oldRequest := range r.requestLog {
		if requestedAt-oldRequest > OnceHour {
			r.requestLog = dequeue(r.requestLog)
		}
	}

	// check if the number of requests is less than the rate limit
	if len(r.requestLog) < r.rateLimit {
		// add the request to the queue
		r.requestLog = enqueue(r.requestLog, requestedAt)
		return true
	} else {
		return false
	}
}

func main() {
	// open the input file
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	defer file.Close()

	// remove the output file if it exists
	if err := os.Remove("output.txt"); err != nil {
		fmt.Print("no found output file, then creating a new one...")
	}
	// create the output file
	outFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("error creating output file: ", err)
		return
	}
	defer outFile.Close()

	// create a new scanner to read file by line
	scanner := bufio.NewScanner(file)

	// create a new rate limit core
	rateLimitCore := RateLimiteCore{}

	// line number
	line := 1
	// number of requests
	nbRequests := 0
	// rate limit (hourly)
	rateLimit := 0

	// read the file line by line
	for scanner.Scan() {
		if line == 1 {
			// read the first line
			// read the number of requests and rate limit
			fmt.Sscanf(scanner.Text(), "%d %d", &nbRequests, &rateLimit)

			// set the rate limit
			rateLimitCore.SetRateLimit(rateLimit)
		} else {
			// read the rest of the lines
			// check if the number of requests has been reached
			if line > nbRequests+1 {
				break
			}

			// read the timestamp in string format
			requestedAtStr := scanner.Text()
			// convert string to time (ISO 8601 format equivalent to RFC 3339)
			requestedAt, err := rateLimitCore.StringToUnix(strings.Trim(requestedAtStr, " "))
			if err != nil {
				log.Fatalf("error parsing time: %v", err)
			}

			// check if the number of requests is less than the rate limit
			result := rateLimitCore.CheckRateLimit(requestedAt)

			// write the result to the output file
			_, err = outFile.WriteString("\n" + strconv.FormatBool(result))
			if err != nil {
				fmt.Println("error writing to output file: ", err)
				return
			}

		}
		line++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error scanning input file: %v", err)
	}
}
