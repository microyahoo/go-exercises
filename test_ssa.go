package main

import "fmt"

// S is a struct
type S struct {
	b [8]byte
}

func keys(m map[S]struct{}) [][]byte {
	var z [][]byte
	for k := range m {
		z = append(z, k.b[:])
	}
	return z
}

func main() {
	fmt.Println(keys(map[S]struct{}{
		S{b: [8]byte{1}}: struct{}{},
		S{b: [8]byte{2}}: struct{}{},
	}))
}

// func keys2(m map[S]struct{}) [][]byte {
//     var z [][]byte

//     var iter map.iter[S]struct{}
//     runtime.mapiterinit(map[S]struct{}.(type), m, &iter)

//     s := new(S)
//     for iter.key != nil {
//         if iter.key == nil {
//            panic("key is nil")
//         }
//         *s = *iter.key
//         if s.b == nil {
//            panic("no next key")
//         }

//         runtime.mapiternext(&iter)

//         z0 := s.b[:8]
//         if len(z) + 1 > cap(z) {
//             runtime.growslice([]uint8.(type), z, len(z) + 1)
//         }

//         z[len(z)] = z0
//         z = z[:len(z)+1]
//     }

//     return z
// }
