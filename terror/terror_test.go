// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package terror

import (
	"testing"

	"github.com/juju/errors"
	. "github.com/pingcap/check"
)

func TestT(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testTErrorSuite{})

type testTErrorSuite struct {
}

func (s *testTErrorSuite) TestTError(c *C) {
	c.Assert(Parser.String(), Not(Equals), "")
	c.Assert(Optimizer.String(), Not(Equals), "")
	c.Assert(KV.String(), Not(Equals), "")
	c.Assert(Server.String(), Not(Equals), "")

	parserErr := Parser.New(ErrCode(1), "error 1")
	c.Assert(parserErr.Error(), Not(Equals), "")
	c.Assert(Parser.EqualClass(parserErr), IsTrue)
	c.Assert(Parser.NotEqualClass(parserErr), IsFalse)
	c.Assert(Parser.Equal(parserErr, ErrCode(1)), IsTrue)
	c.Assert(Parser.NotEqual(parserErr, ErrCode(1)), IsFalse)

	c.Assert(Parser.Equal(parserErr, ErrCode(2)), IsFalse)
	c.Assert(Optimizer.EqualClass(parserErr), IsFalse)

	optimizerErr := Optimizer.New(ErrCode(2), "abc %s", "def")
	c.Assert(Optimizer.Equal(optimizerErr, ErrCode(2)), IsTrue)

	c.Assert(Optimizer.Equal(errors.New("abc"), ErrCode(3)), IsFalse)
	c.Assert(Optimizer.Equal(nil, ErrCode(3)), IsFalse)
	c.Assert(Optimizer.EqualClass(errors.New("abc")), IsFalse)
	c.Assert(Optimizer.EqualClass(nil), IsFalse)
}

func (s *testTErrorSuite) TestErrorEqual(c *C) {
	e1 := errors.New("test error")
	c.Assert(e1, NotNil)

	e2 := errors.Trace(e1)
	c.Assert(e2, NotNil)

	e3 := errors.Trace(e2)
	c.Assert(e3, NotNil)

	c.Assert(errors.Cause(e2), Equals, e1)
	c.Assert(errors.Cause(e3), Equals, e1)
	c.Assert(errors.Cause(e2), Equals, errors.Cause(e3))

	e4 := errors.New("test error")
	c.Assert(errors.Cause(e4), Not(Equals), e1)

	e5 := errors.Errorf("test error")
	c.Assert(errors.Cause(e5), Not(Equals), e1)

	c.Assert(ErrorEqual(e1, e2), IsTrue)
	c.Assert(ErrorEqual(e1, e3), IsTrue)
	c.Assert(ErrorEqual(e1, e4), IsTrue)
	c.Assert(ErrorEqual(e1, e5), IsTrue)

	var e6 error

	c.Assert(ErrorEqual(nil, nil), IsTrue)
	c.Assert(ErrorNotEqual(e1, e6), IsTrue)
	code1 := ErrCode(1)
	code2 := ErrCode(2)
	te1 := Parser.New(code1, "abc")
	te2 := Parser.New(code1, "def")
	te3 := KV.New(code1, "abc")
	te4 := KV.New(code2, "abc")
	c.Assert(ErrorEqual(te1, te2), IsTrue)
	c.Assert(ErrorEqual(te1, te3), IsFalse)
	c.Assert(ErrorEqual(te3, te4), IsFalse)
}
