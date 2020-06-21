#include "textflag.h"

DATA world<>+0(SB)/8, $"hello wo"
DATA world<>+8(SB)/4, $"rld "

GLOBL world<>+0(SB), RODATA, $12

TEXT ·hello(SB),$88-0
	SUBQ	$88, SP
	MOVQ	BP, 80(SP)
	LEAQ	80(SP), BP

	LEAQ	world<>+0(SB), AX 
	MOVQ	AX, my_string+48(SP)        
	MOVQ	$11, my_string+56(SP)
	MOVQ	$0, autotmp_0+64(SP)
	MOVQ	$0, autotmp_0+72(SP)
	LEAQ	type·string(SB), AX
	MOVQ	AX, (SP)
	LEAQ	my_string+48(SP), AX        
	MOVQ	AX, 8(SP)

	CALL	runtime·convT2E(SB)           
	MOVQ	24(SP), AX
	MOVQ	16(SP), CX                    
	MOVQ	CX, autotmp_0+64(SP)        
	MOVQ	AX, autotmp_0+72(SP)
	LEAQ	autotmp_0+64(SP), AX        
	MOVQ	AX, (SP)                      
	MOVQ	$1, 8(SP)                      
	MOVQ	$1, 16(SP)

	CALL	fmt·Println(SB)

	MOVQ 80(SP), BP
	ADDQ $88, SP
	RET
