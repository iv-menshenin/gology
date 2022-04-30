package gology

import "time"

func Int(name string, i int) Attr {
	return Attr{name: name, int: int64(i)}
}

func Int16(name string, i int16) Attr {
	return Attr{name: name, int: int64(i)}
}

func Int32(name string, i int32) Attr {
	return Attr{name: name, int: int64(i)}
}

func Int64(name string, i int64) Attr {
	return Attr{name: name, int: i}
}

func UInt(name string, i uint) Attr {
	return Attr{name: name, uint: uint64(i)}
}

func UInt16(name string, i uint16) Attr {
	return Attr{name: name, uint: uint64(i)}
}

func UInt32(name string, i uint32) Attr {
	return Attr{name: name, uint: uint64(i)}
}

func UInt64(name string, i uint64) Attr {
	return Attr{name: name, uint: i}
}

func String(name, s string) Attr {
	return Attr{name: name, str: s}
}

func DateTime(name string, tm time.Time) Attr {
	return Attr{name: name, tm: tm}
}

func Err(err error) Attr {
	return Attr{name: "error", err: err}
}
