package key

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type pathElement struct {
	id     uint32
	harden bool
}

type path []pathElement

func newPath(s string) (path, error) {
	// remove all whitespace
	s = strings.Replace(s, " ", "", -1)
	// treat root specially
	if s == "/" {
		return path{}, nil
	}
	// now validate the path
	// note that other than the pure root marker that we already handled,
	// the numeric part after the slash is not optional
	valpat := regexp.MustCompile("^(/([0-9]+)'?)+$")
	if !valpat.MatchString(s) {
		return nil, errors.New("Not a valid path string")
	}

	parsepat := regexp.MustCompile("/([0-9]+)('?)")
	saa := parsepat.FindAllStringSubmatch(s, -1)
	// saa now has one entry for each path element, and
	// for each entry it has the 0th element as the whole path string,
	// the first as the path ID, and the second as either
	// an apostrophe or an empty string.
	p := make(path, len(saa))
	for i := range saa {
		var err error
		n, err := strconv.ParseUint(saa[i][1], 10, 32)
		if err != nil {
			return nil, err
		}
		p[i].id = uint32(n)
		p[i].harden = saa[i][2] == "'"
	}
	return p, nil
}

func (p path) isParentOf(c path) bool {
	// if the parent is not shorter than the purported child, it can't be a parent
	if len(c) <= len(p) {
		return false
	}
	// everything up to the length of the parent has to be the same
	for i := range p {
		if c[i].id != p[i].id || c[i].harden != p[i].harden {
			return false
		}
	}
	return true
}

// DeriveFrom accepts a parent key and its known path, plus a desired child path
// and derives the child key from the parent according to the path info.
//
// Note that the parent's known path is simply believed -- we have no mechanism to
// check that it's true.
func (k *ExtendedKey) DeriveFrom(parentPath, childPath string) (*ExtendedKey, error) {
	ppath, err := newPath(parentPath)
	if err != nil {
		return nil, err
	}
	cpath, err := newPath(childPath)
	if err != nil {
		return nil, err
	}
	if !ppath.isParentOf(cpath) {
		return nil, errors.New("child is not descended from parent")
	}
	// if we get here we know that ppath is a subset of cpath so we can trim cpath
	cpath = cpath[len(ppath):]

	// now iterate. Note we never assign to *k, so the origin pointer is unchanged.
	for _, e := range cpath {
		if e.harden {
			k, err = k.HardenedChild(e.id)
		} else {
			k, err = k.Child(e.id)
		}
		if err != nil {
			return nil, err
		}
	}
	return k, err
}
