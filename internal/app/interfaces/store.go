package interfaces

type Store interface {
	Entity() EntityRepo
	MemData() MemDataRepo
}
