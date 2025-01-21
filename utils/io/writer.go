package io

import "io"

type MultiWriter []io.Writer

func NewMultiWriter(ws ...io.Writer) *MultiWriter {
	return (*MultiWriter)(&ws)
}

func (w *MultiWriter) Write(p []byte) (n int, err error) {
	for _, writer := range *w {
		n, err := writer.Write(p)
		if err != nil {
			return n, err
		}
	}

	return len(p), nil
}
