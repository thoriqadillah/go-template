package seeder

import "github.com/stephenafamo/bob"

type Seeder interface {
	Seed(db *bob.DB) error
}

// INFO: register all the seeder here
var seeders = []Seeder{}
