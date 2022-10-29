package server

import (
	"pcbook/pb"
	"sync"

	"github.com/jinzhu/copier"
)

type LaptopStore interface {
	// saves an instance of pb.Laptop
	Save(*pb.Laptop) error

	// potentially returns an instance of pb.Laptop and an error
	// given the ID string of the laptop
	Find(string) (*pb.Laptop, error)
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

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
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
