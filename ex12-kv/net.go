package main

import (
	"math"
	"net"
)

/*
 * req - op 1 byte, key 16 bytes, value 8 bytes
 * resp - code 1 byte value 8 bytes
 */

type ErrorSize struct{}

func (e *ErrorSize) Error() string {
	return "Wrong Size"
}

type Request struct {
	Op    string
	Key   string
	Value float64
}

type Response struct {
	Code  string
	Value string
}

func num(s byte) byte {
	switch {
	case s >= '0' && s <= '9':
		return s - '0'
	case s >= 'a' && s <= 'f':
		return 10 + s - 'a'
	}

	return s
}

func (r *Request) FormRequest() []byte {
	req := make([]byte, 25)

	req[0] = byte(r.Op[0])
	for i := 0; i < 16; i++ {
		req[1+i] = 16 * num(r.Key[2*i])
		req[1+i] += num(r.Key[2*i+1])
	}

	n := math.Float64bits(r.Value)
	req[17] = byte(n >> 56)
	req[18] = byte(n >> 48)
	req[19] = byte(n >> 40)
	req[20] = byte(n >> 32)
	req[21] = byte(n >> 24)
	req[22] = byte(n >> 16)
	req[23] = byte(n >> 8)
	req[24] = byte(n)

	return req
}

func (r *Response) FormResponse() []byte {
	resp := make([]byte, 9)

	resp[0] = byte(r.Code[0])

	for i := 0; i < 8; i++ {
		resp[1+i] = byte(r.Value[i])
	}

	return resp
}

func alpha(i byte) byte {
	switch {
	case i >= 0 && i <= 9:
		return '0' + i
	case i <= 15 && i >= 10:
		return 'a' + (i - 10)
	}
	return i
}

func (r *Request) ParseRequest(req []byte) error {
	if len(req) != 25 {
		return &ErrorSize{}
	}

	r.Op = string(req[0])

	r.Key = ""
	key := req[1:17]
	for _, v := range key {
		p1 := v / 16
		p2 := v % 16
		r.Key += string(alpha(p1))
		r.Key += string(alpha(p2))
	}

	var bits uint64 = 0
	for _, v := range req[17:25] {
		bits *= 256
		bits += uint64(v)
	}

	r.Value = math.Float64frombits(bits)

	return nil
}

func (r *Response) ParseResponse(resp []byte) error {
	if len(resp) != 9 {
		return &ErrorSize{}
	}
	r.Code = string(resp[0])
	r.Value = string(resp[1:])

	return nil
}

func (r *Request) Send(conn net.Conn) {
	conn.Write(r.FormRequest())
}

func (r *Request) Get(conn net.Conn) error {
	buf := make([]byte, 25)

	_, err := conn.Read(buf)

	if err != nil {
		return err
	}

	err = r.ParseRequest(buf)

	if err != nil {
		return err
	}

	return nil
}

func (r *Response) Send(conn net.Conn) {
	conn.Write(r.FormResponse())
}

func (r *Response) Get(conn net.Conn) error {
	buf := make([]byte, 9)

	_, err := conn.Read(buf)

	if err != nil {
		return err
	}

	err = r.ParseResponse(buf)

	if err != nil {
		return err
	}

	return nil
}
