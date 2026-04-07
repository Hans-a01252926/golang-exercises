package main

import "testing"

func TestStackPushAndTop(t *testing.T) {
	s := Stack{}
	s.Push(10)
	s.Push(20)
	s.Push(30)

	got := s.Top()
	want := 30

	if got != want {
		t.Fatalf("Top() = %d; want %d", got, want)
	}
}

func TestStackPop(t *testing.T) {
	s := Stack{}
	s.Push(1)
	s.Push(2)
	s.Push(3)

	s.Pop()

	got := s.Top()
	want := 2

	if got != want {
		t.Fatalf("After Pop(), Top() = %d; want %d", got, want)
	}
}

func TestStackPopEmpty(t *testing.T) {
	s := Stack{}
	s.Pop()

	got := s.Top()
	want := 0

	if got != want {
		t.Fatalf("Top() on empty stack = %d; want %d", got, want)
	}
}

func TestQueueEnqueueAndFront(t *testing.T) {
	q := Queue{}
	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)

	got := q.Front()
	want := 10

	if got != want {
		t.Fatalf("Front() = %d; want %d", got, want)
	}
}

func TestQueueDequeue(t *testing.T) {
	q := Queue{}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	q.Dequeue()

	got := q.Front()
	want := 2

	if got != want {
		t.Fatalf("After Dequeue(), Front() = %d; want %d", got, want)
	}
}

func TestQueueDequeueEmpty(t *testing.T) {
	q := Queue{}
	q.Dequeue()

	got := q.Front()
	want := 0

	if got != want {
		t.Fatalf("Front() on empty queue = %d; want %d", got, want)
	}
}

func TestHashTableCreation(t *testing.T) {
	h := hashTable()

	if h["first"] != "1st" {
		t.Fatalf("h[\"first\"] = %s; want %s", h["first"], "1st")
	}
	if h["second"] != "2nd" {
		t.Fatalf("h[\"second\"] = %s; want %s", h["second"], "2nd")
	}
	if h["third"] != "3rd" {
		t.Fatalf("h[\"third\"] = %s; want %s", h["third"], "3rd")
	}
}

func TestHashTableAdd(t *testing.T) {
	h := hashTable()
	add(h, "fourth", "4th")

	got := h["fourth"]
	want := "4th"

	if got != want {
		t.Fatalf("h[\"fourth\"] = %s; want %s", got, want)
	}
}

func TestHashTableRemove(t *testing.T) {
	h := hashTable()
	remove(h, "second")

	_, exists := h["second"]
	if exists {
		t.Fatal("Expected key 'second' to be removed")
	}
}

func TestHashTableUpdate(t *testing.T) {
	h := hashTable()
	add(h, "first", "updated")

	got := h["first"]
	want := "updated"

	if got != want {
		t.Fatalf("h[\"first\"] = %s; want %s", got, want)
	}
}
