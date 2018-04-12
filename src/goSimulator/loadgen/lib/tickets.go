package lib

import (
	"errors"
	"fmt"
)

type GoTickets interface {
	// 和POSIX的多值信号灯比较像
	Take()
	Return()
	Active() bool
	Total() uint32
	Remainder() uint32
}

func NewGoTickets(total uint32) (GoTickets, error) {
	gt := implGoTickets{}
	if !gt.init(total) {
		errMsg := fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return &gt, nil
}

type implGoTickets struct {
	total    uint32
	ticketCh chan byte //票的容器
	active   bool
}

func (self *implGoTickets) Take() {
	<-self.ticketCh
}
func (self *implGoTickets) Return() {
	self.ticketCh <- 1
}
func (self *implGoTickets) Active() bool {
	return self.active
}
func (self *implGoTickets) Total() uint32 {
	return self.total
}
func (self *implGoTickets) Remainder() uint32 {
	return uint32(len(self.ticketCh))
}
func (self *implGoTickets) init(total uint32) bool {
	if self.active {
		return false
	}
	if total == 0 {
		return false
	}
	ch := make(chan byte, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- 1
	}
	self.ticketCh = ch
	self.total = total
	self.active = true
	return true
}
