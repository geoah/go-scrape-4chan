package main

import (
	r "gopkg.in/gorethink/gorethink.v3"
)

// Persistence -
type Persistence interface {
	Put(...interface{}) error
	GetThread(id string) (*Thread, error)
}

// RethinkPersistence -
type RethinkPersistence struct {
	session *r.Session
}

// Put -
// TODO(geoah): Do we care about which one erroed out?
func (p *RethinkPersistence) Put(table string, shouldError bool, entries ...interface{}) error {
	conflict := "update"
	if shouldError {
		conflict = "error"
	}
	for _, entry := range entries {
		q := r.Table(table).Insert(entry, r.InsertOpts{
			Conflict: conflict,
		})
		if _, err := q.RunWrite(p.session); err != nil {
			return err
		}
	}
	return nil
}

// GetThread -
func (p *RethinkPersistence) GetThread(id string) (*Thread, error) {
	thread := &Thread{}
	q := r.Table(threadsTable).Get(id)
	res, err := q.Run(p.session)
	if err != nil {
		return nil, err
	}
	if err := res.One(&thread); err != nil {
		return nil, err
	}
	return thread, nil
}
