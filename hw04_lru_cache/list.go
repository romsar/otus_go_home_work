package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v any) *ListItem {
	node := NewListItem()
	node.Value = v

	if l.front == nil {
		l.front = node
		l.back = node
	} else {
		node.Next = l.front
		l.front.Prev = node
		l.front = node
	}

	l.len++

	return node
}

func (l *list) PushBack(v any) *ListItem {
	node := NewListItem()
	node.Value = v

	if l.front == nil {
		l.front = node
		l.back = node
	} else {
		l.back.Next = node
		node.Prev = l.back
		l.back = node
	}

	l.len++

	return node
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	switch {
	case i.Prev != nil && i.Next != nil:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	case i.Prev == nil && i.Next != nil:
		l.front = i.Next.Prev
	case i.Prev != nil && i.Next == nil:
		l.back = i.Prev
		l.back.Next = nil
	case i.Prev == nil && i.Next == nil:
		l.front = nil
		l.back = nil
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return new(list)
}

func NewListItem() *ListItem {
	return new(ListItem)
}
