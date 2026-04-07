package main

import (
	"fmt"
)

// Funcion para verificar si stacks  queues esta vacia
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Parte de Stacks (LIFO)
type Stack struct {
	items []int
}

func (s *Stack) Push(item int) { // Agrega item a la pila
	fmt.Println("Push: ", item)
	s.items = append(s.items, item)
}

func (s *Stack) Pop() { // Borra item que entro a la pila
	if s.IsEmpty() {
		return
	}
	fmt.Println("Pop: ", s.items[len(s.items)-1])
	s.items = s.items[:len(s.items)-1]
}

func (s *Stack) Top() int { // Retorna item que acaba de entrar
	if s.IsEmpty() {
		return 0
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) Print() []int {
	for _, item := range s.items {
		fmt.Println(item)
	}
	return s.items
}

// Parte de Queues (FIFO)
type Queue struct {
	items []int
}

func (q *Queue) Enqueue(item int) {
	fmt.Println("Enqueue: ", item)
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() { // Borra item que entro
	if len(q.items) == 0 {
		return
	}
	fmt.Println("Dequeue: ", q.items[0])
	q.items = q.items[1:]

}

func (q *Queue) Front() int { // Retorna item que acaba de entrar
	if len(q.items) == 0 {
		return 0
	}
	return q.items[0]
}

func (q *Queue) Print() []int {
	for _, item := range q.items {
		fmt.Println(item)
	}
	return q.items
}

func orderedMap() map[string]string {
	ordered := make(map[string]string)
	ordered["first"] = "1st"
	ordered["second"] = "2nd"
	ordered["third"] = "3rd"
	return ordered
}

func hashTable() map[string]string {
	hashtable := make(map[string]string)
	hashtable["first"] = "1st"
	hashtable["second"] = "2nd"
	hashtable["third"] = "3rd"
	return hashtable
}

func add(m map[string]string, key string, value string) {
	m[key] = value
}

func remove(m map[string]string, key string) {
	delete(m, key)
}

func printHash(m map[string]string) {
	for key, value := range m {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func main() {
	fmt.Println("Mi implementacion de Stacks:")
	ejemploStack := Stack{}
	ejemploStack.Push(1)
	ejemploStack.Push(2)
	ejemploStack.Push(3)
	ejemploStack.Push(4)
	ejemploStack.Pop()

	fmt.Println("Top de Pila: ")
	fmt.Println(ejemploStack.Top())

	fmt.Println("Pila completa: ")
	ejemploStack.Print()

	fmt.Println("\nMi implementacion de Queues: ")
	ejemploQueue := Queue{}
	ejemploQueue.Enqueue(1)
	ejemploQueue.Enqueue(2)
	ejemploQueue.Enqueue(3)
	ejemploQueue.Enqueue(4)
	ejemploQueue.Dequeue()

	fmt.Println("Front de Queue: ")
	fmt.Println(ejemploQueue.Front())

	fmt.Println("Queue completa: ")
	ejemploQueue.Print()

	fmt.Println("\nMi implementacion de HashTable: ")
	ejemploHash := hashTable()
	printHash(ejemploHash)

	add(ejemploHash, "fourth", "4th")
	fmt.Println("\nHashTable despues de agregar un nuevo elemento: ")
	printHash(ejemploHash)

	remove(ejemploHash, "second")
	fmt.Println("\nHashTable despues de eliminar un elemento: ")
	printHash(ejemploHash)
}
