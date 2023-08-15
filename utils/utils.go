package utils

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

type CivilTime time.Time

func (c *CivilTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) // get rid of "
	if value == "" || value == "null" {
		return nil
	}
	timeFoemats := [2]string{
		"2006-01-02",
		// "02-01-2006",
	}
	layout := timeFoemats[0]

	t, err := time.Parse(layout, value) // parse time
	if err != nil {
		return err
	}
	current_time := time.Now()
	if t.Before(current_time) {
		return errors.New("date allready passed")
	}
	*c = CivilTime(t) // set result using the pointer
	return nil
}

func RandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
