package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
)

type (
	Report struct {
		Meta
		DataInsert    Variants
		DataGet       Variants
		DataGetAll    Variants
		DataGetHeader Variants
	}

	Meta struct {
		Nums     []uint
		Variants []string
	}

	Variants = map[string]Series
	Series   = map[uint]Stamp // size per stamp

	Stamp struct {
		Elapse float64
		Bytes  uint64
	}
)

func (r *Report) Fill(f io.Reader) error {

	var rd = bufio.NewScanner(f)

	r.DataInsert = make(Variants)
	r.DataGet = make(Variants)
	r.DataGetAll = make(Variants)
	r.DataGetHeader = make(Variants)

	for i := 0; i < 3; i++ {
		switch {
		case !rd.Scan():
			const D = "failed to skip first 3 lines, line: %d, err: %w"
			return fmt.Errorf(D, i+1, rd.Err())

		case rd.Err() != nil:
			const D = "skipping first 3 lines is OK, but got err, line: %d, err: %w"
			return fmt.Errorf(D, i+1, rd.Err())

		case strings.HasPrefix(rd.Text(), "Benchmark"):
			const D = "some of first 3 line contain data, not header, line: %d"
			return fmt.Errorf(D, i+1)
		}
	}

	var readLines int
	for rd.Scan() {
		const Prefix = "BenchmarkUkvs/"

		var s = rd.Text()
		if !strings.HasPrefix(s, Prefix) {
			break
		}

		var err = r.parseLine(s)
		if err != nil {
			return fmt.Errorf("[%d line]: %w", readLines, err)
		}
		readLines++
	}

	log.Printf("READ %d LINES\n", readLines)
	return nil
}

func (r *Report) Aggregate() error {

	var anyKey string
	for anyKey = range r.DataInsert {
		break
	}

	var series = r.DataInsert[anyKey]
	r.Nums = make([]uint, 0, len(series))
	for num := range series {
		r.Nums = append(r.Nums, num)
	}

	r.Variants = make([]string, 0, len(r.DataInsert))
	for key := range r.DataInsert {
		r.Variants = append(r.Variants, key)
	}

	slices.Sort(r.Nums)
	slices.Sort(r.Variants)

	return nil
}

func (r *Report) parseLine(s string) error {
	const Prefix = "BenchmarkUkvs/"

	if !strings.HasPrefix(s, Prefix) {
		return fmt.Errorf("benchmark prefix is not found (%s)", s)
	}

	s = s[len(Prefix):]
	var dest Variants

	const PrefixActionInsert = "Insert/"
	const PrefixActionGet = "Get/"
	const PrefixActionGetAll = "GetAll/"
	const PrefixActionGetHeader = "GetHeader/"

	switch {
	case strings.HasPrefix(s, PrefixActionInsert):
		dest = r.DataInsert
		s = s[len(PrefixActionInsert):]

	case strings.HasPrefix(s, PrefixActionGetAll):
		dest = r.DataGetAll
		s = s[len(PrefixActionGetAll):]

	case strings.HasPrefix(s, PrefixActionGetHeader):
		dest = r.DataGetHeader
		s = s[len(PrefixActionGetHeader):]

	case strings.HasPrefix(s, PrefixActionGet):
		dest = r.DataGet
		s = s[len(PrefixActionGet):]

	default:
		return fmt.Errorf("action prefix is not found (%s)", s)
	}

	var idx = strings.IndexByte(s, '/')
	if idx == -1 {
		return fmt.Errorf("action variant is not found (%s)", s)
	}

	var variantStr = s[:idx]
	s = s[idx+1:]

	var variantDest Series
	if variantDest = dest[variantStr]; variantDest == nil {
		dest[variantStr] = make(Series)
		variantDest = dest[variantStr]
	}

	for idx = 0; idx < len(s); idx++ {
		if s[idx] < '0' || s[idx] > '9' {
			break
		}
	}

	if idx == 0 {
		return fmt.Errorf("invalid number of series (%s)", s)
	}

	var num, err = strconv.ParseUint(s[:idx], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid number of series (%s): %w", s, err)
	}

	// Skip until spaces

	s = s[idx:]
	for idx = 0; idx < len(s); idx++ {
		if s[idx] <= ' ' {
			break
		}
	}

	// Skip until number

	for ; idx < len(s); idx++ {
		if s[idx] >= '0' && s[idx] <= '9' {
			break
		}
	}

	// Now it's number of iterations.
	// Don't need for now. Also skip.

	for ; idx < len(s); idx++ {
		if s[idx] < '0' || s[idx] > '9' {
			break
		}
	}

	for ; idx < len(s); idx++ {
		if s[idx] >= '0' && s[idx] <= '9' {
			break
		}
	}

	// It's elapse of one iteration now (NS).

	s = s[idx:]
	for idx = 0; idx < len(s); idx++ {
		if !(s[idx] >= '0' && s[idx] <= '9' || s[idx] == '.') {
			break
		}
	}

	var elapse float64
	if elapse, err = strconv.ParseFloat(s[:idx], 32); err != nil {
		return fmt.Errorf("invalid number of elapse time (%s): %w", s, err)
	}

	// Skip until number

	s = s[idx:]
	for idx = 0; idx < len(s); idx++ {
		if s[idx] >= '0' && s[idx] <= '9' {
			break
		}
	}

	s = s[idx:]
	for idx = 0; idx < len(s); idx++ {
		if !(s[idx] >= '0' && s[idx] <= '9') {
			break
		}
	}

	var bytes uint64
	if bytes, err = strconv.ParseUint(s[:idx], 10, 32); err != nil {
		return fmt.Errorf("invalid number of alloc bytes (%s): %w", s, err)
	}

	variantDest[uint(num)] = Stamp{Elapse: elapse, Bytes: bytes}
	return nil
}
