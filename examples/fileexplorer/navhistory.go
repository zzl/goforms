package main

import "bytes"

type NavHistory struct {
	items []*ShellItem
	index int
}

func (this *NavHistory) Push(item *ShellItem) {
	if this.index < len(this.items) &&
		bytes.Equal(this.items[this.index].id, item.id) {
		return
	}
	if this.items == nil {
		this.items = []*ShellItem{item}
	} else {
		this.items = append(this.items[:this.index+1], item)
	}
	this.index = len(this.items) - 1
}

func (this *NavHistory) CanBack() bool {
	return this.index > 0
}

func (this *NavHistory) Back() *ShellItem {
	this.index -= 1
	return this.items[this.index]
}

func (this *NavHistory) CanForward() bool {
	return this.index < len(this.items)-1
}

func (this *NavHistory) Forward() *ShellItem {
	this.index += 1
	return this.items[this.index]
}
