#include "textflag.h"

#define SWAP(x, y, t) MOVQ x, t; MOVQ y, x; MOVQ t, y

// func Swap(a, b int) (int, int)
TEXT ·Swap(SB), $0-32
    MOVQ a+0(FP), AX // AX = a
    MOVQ b+8(FP), BX // BX = b

    SWAP(AX, BX, CX)     // AX, BX = b, a

    MOVQ AX, ret0+16(FP) // return
    MOVQ BX, ret1+24(FP) //
    RET

TEXT ·test(SB), $24-0
	MOVQ $0, a-8*2(SP) // a = 0
	MOVQ $0, b-8*1(SP) // b = 0

	MOVQ $10, AX
	MOVQ AX, a-8*2(SP)
	// MOVQ $10, a-8*2(SP) 

	MOVQ AX, 0(SP)
	CALL ·println(SB)
	// CALL runtime·printint(SB)
	// CALL runtime·printnl(SB)

	RET

TEXT ·printnl_nosplit(SB), NOSPLIT, $8
// TEXT ·printnl_nosplit(SB), $8
	// CALL runtime·printnl(SB)
	MOVQ $10, AX

	MOVQ AX, 0(SP)
	CALL ·println(SB)
	RET
