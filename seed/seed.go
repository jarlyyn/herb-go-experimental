package seed

type Seed interface {
	HarvestSeed() ([]byte, error)
}

type SeedStorer interface {
	StoreSeed(Seed) error
}

type SeedLoader interface {
	LoadeSeed() (Seed, error)
}
