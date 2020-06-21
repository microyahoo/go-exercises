#include "textflag.h"

// TEXT ·getg(SB), ABIInternal, $32-16
// TEXT	"".main.func1(SB), ABIInternal, $88-8
TEXT ·getg(SB), NOSPLIT|ABIInternal, $32-16
	MOVQ (TLS), AX //get runtime.g
	MOVQ $type·runtime·g(SB), BX // get runtime.g type
	MOVQ AX, 8(SP)
	MOVQ BX, 0(SP)
	CALL runtime·convT2E(SB) // convert (*g) to interface{}
	MOVQ 16(SP), AX
	MOVQ 24(SP), BX

	// return interface{}
	MOVQ AX, ret+0(FP)
	MOVQ BX, ret+8(FP)
	RET
