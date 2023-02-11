package main

type StorageType int

const (
	Warehouse = iota
)

type Storage struct {
	Storage   StorageType
	Dice      map[ItemType]map[int]uint64 // Stores number per face per dice type
	Capacity  uint64                      // How many dice can be held
	Count     uint64                      // How many dice are held?
	TypeLimit int                         // How many types of dice can be held? 0 for infinite
	TypeCount int                         // How many types of dice are there?
}

func NewStorage(storageType StorageType, capacity uint64, typeLimit int) *Storage {
	return &Storage{
		Dice:      map[ItemType]map[int]uint64{},
		Storage:   storageType,
		Capacity:  capacity,
		TypeLimit: typeLimit,
	}
}

// StoreDie adds a die of given type to storage.
// Returns false if sum is not less than capacity or the type limit is reached.
func (s *Storage) StoreDie(item ItemType, face int) bool {
	if s.Count >= s.Capacity {
		return false
	}
	_, exists := s.Dice[item]
	if exists {
		_, exists = s.Dice[item][face]
		if exists {
			s.Dice[item][face]++
			s.Count++
			return true
		} else {
			s.Dice[item][face] = 1
			s.Count++
			return true
		}
	} else if s.TypeLimit == 0 || s.TypeCount <= s.TypeLimit {
		s.Dice[item] = map[int]uint64{face: 1}
		s.Count++
		s.TypeCount++
		return true
	}
	return false
}

// RemoveDie removes a die from a storage. Should not be used by trucks,
// as it will not reduce the TypeCount
func (s *Storage) RemoveDie(item ItemType, face int) bool {
	if s.Count <= 0 {
		return false
	}

	dice, exists := s.Dice[item]
	if exists {
		_, exists = dice[face]
		if exists {
			if dice[face] > 0 {
				s.Dice[item][face]--
				s.Count--
				return true
			}
		}
	}
	return false
}

// Load adds all the dice in given storage to self.
// Returns false if the added dice would exceed the capacity or type limit.
func (s *Storage) Load(storage Storage) bool {
	return false
}
