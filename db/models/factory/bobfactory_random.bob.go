// Code generated by BobGen psql v0.30.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

import (
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jaswdr/faker/v2"
)

var defaultFaker = faker.New()

func random_string(f *faker.Faker) string {
	if f == nil {
		f = &defaultFaker
	}

	return strings.Join(f.Lorem().Words(f.IntBetween(1, 5)), " ")
}

func random_time_Time(f *faker.Faker) time.Time {
	if f == nil {
		f = &defaultFaker
	}

	year := time.Hour * 24 * 365
	min := time.Now().Add(-year)
	max := time.Now().Add(year)
	return f.Time().TimeBetween(min, max)
}

func random_uuid_UUID(f *faker.Faker) uuid.UUID {
	if f == nil {
		f = &defaultFaker
	}

	return uuid.Must(uuid.NewV4())
}
