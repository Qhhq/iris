// Copyright 2017 Gerasimos Maropoulos, ΓΜ. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"net/http"
	"sync"
)

// Pool is the context pool, it's used inside router and the framework by itself.
//
// It's the only one real implementation inside this package because it used widely.
type Pool struct {
	pool    *sync.Pool
	newFunc func() Context // we need a field otherwise is not working if we change the return value
}

// New creates and returns a new context pool.
func New(newFunc func() Context) *Pool {
	c := &Pool{pool: &sync.Pool{}, newFunc: newFunc}
	c.pool.New = func() interface{} { return c.newFunc() }
	return c
}

// Attach changes the pool's return value Context.
func (c *Pool) Attach(newFunc func() Context) {
	c.newFunc = newFunc
}

// Acquire returns a Context from pool.
// See Release.
func (c *Pool) Acquire(w http.ResponseWriter, r *http.Request) Context {
	ctx := c.pool.Get().(Context)
	ctx.BeginRequest(w, r)
	return ctx
}

// Release puts a Context back to its pull, this function releases its resources.
// See Acquire.
func (c *Pool) Release(ctx Context) {
	ctx.EndRequest()
	c.pool.Put(ctx)
}
