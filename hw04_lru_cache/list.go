package hw04lrucache

import (
	"fmt"
	"strings"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	count int
	head  *ListItem
	tail  *ListItem
}

func (l *list) String() string {
	s := strings.Builder{}
	item := l.head
	for item != nil {
		var prevVal, nextVal any
		if item.Prev != nil {
			prevVal = item.Prev.Value
		}
		if item.Next != nil {
			nextVal = item.Next.Value
		}
		s.WriteString(fmt.Sprintf("prev: %v, value: %v, next %v\n", prevVal, item.Value, nextVal))
		item = item.Next
	}
	return s.String()
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFront := &ListItem{
		Prev:  nil,
		Value: v,
		Next:  l.head,
	}
	if l.head == nil {
		l.head = newFront
		l.tail = newFront
	} else {
		l.head.Prev = newFront
		l.head = newFront
	}
	l.count++
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{
		Prev:  l.tail,
		Value: v,
		Next:  nil,
	}
	if l.tail == nil {
		l.head = newBack
		l.tail = newBack
	} else {
		l.tail.Next = newBack
		l.tail = newBack
	}
	l.count++
	return l.head
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
		i.Next = l.head
		l.head.Prev = i
		l.head = i
		i.Prev = nil
	}
}

func NewList() List {
	return new(list)
}
