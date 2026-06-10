package colors

import "encoding/binary"

func MapIterationsToUint16Colors(iterations uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, iterations)
	return buf
}
