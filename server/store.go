package server

import (
	"context"
	"log"
	"pcstore/pb"
	"sync"
)

type LaptopStore interface {
	// saves an instance of pb.Laptop
	Save(*pb.Laptop) error

	// potentially returns an instance of pb.Laptop and an error
	// given the ID string of the laptop
	Find(string) (*pb.Laptop, error)

	// search for valid laptops given a specific filter
	Search(context.Context, *pb.Filter, func(laptop *pb.Laptop) error) error
}

type JsonLaptopStore struct {
	mutex sync.RWMutex
	path  string
}

func (store *JsonLaptopStore) Save(laptop *pb.Laptop) error {
	// append the laptop to the existing .json file
	// TODO implement this feature
	return nil
}

func (store *JsonLaptopStore) Find(id string) (*pb.Laptop, error) {
	// returns an instance of pb.Laptop given an id
	// or potentially an error
	// TODO implement this feature
	return nil, nil
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

func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		err := ctx.Err()
		if err == context.Canceled || err == context.DeadlineExceeded {
			log.Println("context is canceled")
			return err // is that better to return nil?
		}

		if isQualifiedLaptop(laptop, filter) {
			copied, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(copied)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		mutex: sync.RWMutex{},
		data:  make(map[string]*pb.Laptop),
	}
}
