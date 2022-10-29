package server

import (
	"pcbook/pb"
	"sync"

	"github.com/jinzhu/copier"
)

type LaptopStore interface {
	Save(*pb.Laptop) error
}

type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.data[laptop.Id]; ok {
		return ErrorLaptopAlreadyExists
	}

	// deepcopy to prevent unwanted changes in the records
	// couldn't just use pb.Laptop instead of doing this deepcopy on *pb.Laptop?
	copied, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[copied.Id] = copied
	return nil
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, ErrorCanNotCopyLaptop
	}

	return other, nil
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		mutex: sync.RWMutex{},
		data:  make(map[string]*pb.Laptop),
	}
}
