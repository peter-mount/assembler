package memory

type Writer interface {
	Read([]byte) (int, error)
	ReadByte() (byte, error)
	Write([]byte) (int, error)
	WriteByte(byte) error
}

func (m *Map) Writer(addr Address) Writer {
	return &writer{
		memoryMap: m,
		addr:      addr,
	}
}

type writer struct {
	memoryMap *Map
	addr      Address
}

func (w *writer) Read(b []byte) (int, error) {
	for n := 0; n < len(b); n++ {
		v, err := w.ReadByte()
		if err != nil {
			return n, err
		}
		b[n] = v
	}
	return len(b), nil
}

func (w *writer) ReadByte() (byte, error) {
	if w == nil {
		return 0, addressInvalid
	}

	v, err := w.memoryMap.ReadByte(w.addr)
	if err != nil {
		return 0, err
	}

	w.addr = w.addr.Increment()
	return v, nil
}

func (w *writer) Write(b []byte) (int, error) {
	n := 0
	for _, v := range b {
		if err := w.WriteByte(v); err != nil {
			return 0, err
		}
		n++
	}
	return n, nil
}

func (w *writer) WriteByte(v byte) error {
	if w == nil {
		return addressInvalid
	}

	if err := w.memoryMap.WriteByte(w.addr, v); err != nil {
		return err
	}

	w.addr = w.addr.Increment()
	return nil
}
