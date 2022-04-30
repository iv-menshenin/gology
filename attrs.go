package gology

func Int(name string, i int) Attr {
	return Attr{name: name, int: int64(i)}
}

func String(name, s string) Attr {
	return Attr{name: name, str: s}
}
