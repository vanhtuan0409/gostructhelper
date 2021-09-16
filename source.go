package gostructhelper

import "io"

type source struct {
	r    io.Reader
	path string
}

type SourceOpt func(s *source)

func SourceWithReader(r io.Reader) SourceOpt {
	return func(s *source) {
		if r != nil {
			s.r = r
		}
	}
}

func SourceWithPath(path string) SourceOpt {
	return func(s *source) {
		if path != "" {
			s.path = path
		}
	}
}

func NewSource(opts ...SourceOpt) *source {
	s := &source{
		r:    nil,
		path: "",
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
