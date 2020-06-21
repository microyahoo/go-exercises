#include "textflag.h"

TEXT ·ptrToFunc(SB), NOSPLIT, $0-16
	MOVQ ptr+0(FP), AX
	MOVQ AX, ret+8(FP)
	RET

TEXT ·asmFuncTwiceClosureAddr(SB), NOSPLIT, $0-8
	LEAQ ·asmFuncTwiceClosureBody(SB), AX
	MOVQ AX, ret+0(FP)
	RET

TEXT ·asmFuncTwiceClosureBody(SB), NOSPLIT|NEEDCTXT, $0-8
	MOVQ 8(DX), AX
	ADDQ AX, AX
	MOVQ AX, 8(DX)
	MOVQ AX, ret+0(FP)

	RET

TEXT ·syscallWrite_darwin(SB), NOSPLIT, $0
	MOVQ $(0x2000000+4), AX
	MOVQ fd+0(FP), DI
	MOVQ msg_data+8(FP), SI
	MOVQ msg_len+16(FP), DX
	SYSCALL
	MOVQ AX, ret+0(FP)
	RET

