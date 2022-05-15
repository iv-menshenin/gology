package gology

import "time"

func Int(name string, i int) Attr {
	return Attr{name: name, typ: attrInt, int: int64(i)}
}

func Int16(name string, i int16) Attr {
	return Attr{name: name, typ: attrInt, int: int64(i)}
}

func Int32(name string, i int32) Attr {
	return Attr{name: name, typ: attrInt, int: int64(i)}
}

func Int64(name string, i int64) Attr {
	return Attr{name: name, typ: attrInt, int: i}
}

func UInt(name string, i uint) Attr {
	return Attr{name: name, typ: attrUint, uint: uint64(i)}
}

func UInt16(name string, i uint16) Attr {
	return Attr{name: name, typ: attrUint, uint: uint64(i)}
}

func UInt32(name string, i uint32) Attr {
	return Attr{name: name, typ: attrUint, uint: uint64(i)}
}

func UInt64(name string, i uint64) Attr {
	return Attr{name: name, typ: attrUint, uint: i}
}

func String(name, s string) Attr {
	return Attr{name: name, typ: attrString, str: s}
}

func DateTime(name string, tm time.Time) Attr {
	return Attr{name: name, typ: attrTime, tm: tm}
}

func Err(err error) Attr {
	return Attr{name: "error", typ: attrError, err: err}
}

func Float64(name string, f float64) Attr {
	return Attr{name: name, typ: attrFloat, flt: f}
}

func Float32(name string, f float32) Attr {
	return Attr{name: name, typ: attrFloat, flt: float64(f)}
}
