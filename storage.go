package main

type StorageType int

const (
	Warehouse = iota
	TruckTrailer
)

const warehouseCapacity uint64 = 1000

const (
	truckCapacity  uint64 = 10
	truckTypeLimit int    = 1
)

type Storage struct {
	ID        uint64
	Storage   StorageType
	Dice      map[ItemType]map[int]uint64 // Stores number per face per dice type
	Capacity  uint64                      // How many dice can be held
	Count     uint64                      // How many dice are held?
	TypeLimit int                         // How many types of dice can be held? 0 for infinite
	TypeCount int                         // How many types of dice are there?
}

// NewStorage constructs a new Storage.
// the Dice map is initialised, but containing maps are not.
func (g *Game) NewStorage(storageType StorageType, capacity uint64, typeLimit int) *Storage {
	return &Storage{
		ID:        g.NextID(),
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
		s.Dice[item][face]++
		s.Count++
		return true
	} else if s.TypeLimit == 0 || s.TypeCount <= s.TypeLimit {
		s.Dice[item] = map[int]uint64{face: 1}
		s.Count++
		s.TypeCount++
		return true
	}
	return false
}

// StoreDie adds an amount of dice of given type to storage.
// Returns false if sum is not less than capacity or the type limit is reached.
func (s *Storage) StoreDice(item ItemType, face int, count uint64) bool {
	if s.Count >= s.Capacity {
		return false
	}
	_, exists := s.Dice[item]
	if exists {
		s.Dice[item][face] += count
		s.Count += count
		return true
	} else if s.TypeLimit == 0 || s.TypeCount <= s.TypeLimit {
		s.Dice[item] = map[int]uint64{face: count}
		s.Count += count
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
func (s *Storage) Load(storage *Storage) bool {

	if storage.Count+s.Count > s.Capacity {
		return false
	}

	for diceType, dice := range storage.Dice {
		for face, count := range dice {
			s.StoreDice(diceType, face, count)
		}
	}
	return true
}
