/*

How to Run?
- in the same folder, you create config.txt, which cointains how many elements do you want the client to send to the server
- open a terminal inside the folder and write 'go run primeNo.go'
- open 1 or 2 more terminals in that folder also and connect with the server using 'telnet 127.0.0.1 8080'

What happens here?
-> The client(s), those terminals in which you wrote telnet are sending an array of numbers to the server
-> The server is returning the total number of digits from all the prime numbers in that array
-> For example, if the client writes: 23, 17, 15, 3, 18
=> the output will be '5 digits', because numbers 23, 17, 3 are prime numbers and they have a total digits of 5.

*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"strconv"
	"strings"
)

// boolean function to verify whether a number is prime or not
func isPrime(nr int) bool {
	for i := 2; i <= int(math.Floor(float64(nr)/2)); i++ {
		if nr%i == 0 {
			return false
		}
	}
	return true
}

// counts how many digits there are in a number
func countCifre(nr int) int {
	var count = 0
	for nr > 0 {
		nr = nr / 10
		count++
	}
	return count
}

// checks whether the server got any error or not
func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", message)
}

func main() {
	data, err := ioutil.ReadFile("config.txt")
	if err != nil {
		fmt.Println("Error when reading from file", err)
		return
	}
	length := string(data)          // stringify the data written from the config
	lung, _ := strconv.Atoi(length) // converted the array to integer and then print it to the console
	fmt.Println(lung)

	clientCount := 1
	allClients := make(map[net.Conn]int) // mapping the connection on the channel

	ln, err := net.Listen("tcp", ":8080")
	check(err, "Server is ready.")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		allClients[conn] = clientCount
		fmt.Printf("Client %d Connected.\n", allClients[conn])

		clientCount++

		go func() {
			reader := bufio.NewReader(conn)

			for {
			tag:
				incoming, err := reader.ReadString('\n') // message sent by client

				if err != nil {
					fmt.Printf("Client %d Disconnected.\n", allClients[conn])
					break
				}

				fmt.Printf("Client %d made a request with following data: %s", allClients[conn], incoming)
				conn.Write([]byte("Server got the request.\n"))

				incoming = incoming[0 : len(incoming)-2] // removing the \n from the received message

				v := strings.Split(incoming, ",") // spliting the array with numbers using "," and saving each element from the array using the variable 'v'
				var vector []int
				for i := 0; i < len(v); i++ {
					elem, _ := strconv.Atoi(v[i]) // converting each element from the array to integer
					vector = append(vector, elem)
				} // 'vector' contains the "cuvintele" (words) sent by client

				if len(vector) != lung { // 'lung' is the length configured in the file
					conn.Write([]byte("Error! Write a vector with " + length + " elements.\n"))
					fmt.Println(lung)
					goto tag // continue from the "tag"
				}

				conn.Write([]byte("Server is processing the data.\n"))

				countPrime := 0
				for i := 0; i < len(vector); i++ {
					if isPrime(vector[i]) {
						countPrime = countPrime + countCifre(vector[i])
					}
				}
				countPrimeArray := strconv.Itoa(countPrime)

				name := strconv.Itoa(allClients[conn])
				fmt.Printf("Server is sending " + incoming + " => " + countPrimeArray + "\n")
				conn.Write([]byte("Client " + name + " got the answer: " + incoming + " => " + countPrimeArray + " digits \n"))
			}
		}()
	}
}
