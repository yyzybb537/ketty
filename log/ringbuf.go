package log

import (
	"io"
	"sync"
	"bytes"
)

type Buffer struct {
	mtx sync.RWMutex
	buf *bytes.Buffer
	skipSize int
	useCount int
}

func NewBuffer(bufSkipSize int) *Buffer {
	return &Buffer{
		buf : bytes.NewBuffer([]byte{}),
		skipSize : bufSkipSize,
	}
}

func (this *Buffer) FlushTo(w io.Writer) (n int, err error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.buf.Len() == 0 {
		return
	}
	n, err = w.Write(this.buf.Bytes())
	this.buf.Reset()
	return
}

func (this *Buffer) commit(b []byte) (n int, err error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	return this.buf.Write(b)
}

func (this *Buffer) inc() {
	this.useCount ++
}

func (this *Buffer) dec() {
	this.useCount --
}

func (this *Buffer) full() bool {
	return this.buf.Len() >= this.skipSize
}

type BufferWrap struct {
	*Buffer
	poolBuf *bytes.Buffer
}

func (this *BufferWrap) Write(b []byte) (n int, err error) {
	return this.poolBuf.Write(b)
}

func (this *BufferWrap) WriteString(s string) (n int, err error) {
	return this.poolBuf.WriteString(s)
}

var gBufPool = sync.Pool{
	New : func() interface{} {
		return bytes.NewBufferString("")
	},
}

func NewBufferWrap(base *Buffer) *BufferWrap {
	return &BufferWrap{base, gBufPool.Get().(*bytes.Buffer)}
}

func (this *BufferWrap) commit() {
	if this.poolBuf.Len() > 0 {
		this.Buffer.commit(this.poolBuf.Bytes())
		this.poolBuf.Reset()
	}
	gBufPool.Put(this.poolBuf)
}

type RingBuf struct {
	bufs []*Buffer
	r int
	w int
	mtx sync.RWMutex
}

func NewRingBuf(nBufs int, bufSkipSize int) *RingBuf {
	rb := &RingBuf{}
	rb.bufs = make([]*Buffer, nBufs, nBufs)
	for i := 0; i < nBufs; i++{
		rb.bufs[i] = NewBuffer(bufSkipSize)
	}
	return rb
}

func (this *RingBuf) GetWriter() *BufferWrap {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	buf := this.bufs[this.w]
	buf.inc()
	return NewBufferWrap(buf)
}

func (this *RingBuf) Put(bufw *BufferWrap) {
	bufw.commit()

	this.mtx.Lock()
	defer this.mtx.Unlock()
	bufw.dec()
	if bufw.full() {
		// skipW
		this.skipW();
	}
}

func (this *RingBuf) NextReader() *Buffer {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	return this.nextReaderWithoutLock()
}

func (this *RingBuf) nextReaderWithoutLock() *Buffer {
	buf := this.bufs[this.r]
	if this.r == this.w {
		if buf.buf.Len() <= 0 {
			return nil
		}

		this.skipW();
	}

	if buf.useCount == 0 {
		this.skipR()
		return buf
	}

	return nil
}

func (this *RingBuf) FlushTo(w io.Writer) (err error) {
	for {
		buf := this.NextReader()
		if buf == nil {
			return
		}

		_, err = buf.FlushTo(w)
		if err != nil {
			return
		}
	}
	return
}

func (this *RingBuf) LockFlushTo(w io.Writer) (err error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for {
		buf := this.nextReaderWithoutLock()
		if buf == nil {
			return
		}

		_, err = buf.FlushTo(w)
		if err != nil {
			return
		}
	}
	return
}

func (this *RingBuf) skipW() {
	newW := (this.w + 1) % len(this.bufs)
	if newW == this.r {
		return
	}
	this.w = newW
}

func (this *RingBuf) skipR() {
	this.r = (this.r + 1) % len(this.bufs)
}
