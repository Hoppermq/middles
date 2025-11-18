package pkg

type UUID [16]byte

type UUIDGenerator interface {
	Generate() (UUID, error)
}
