; ********************************************************************************
; This is a simple Hello World example for the BBC micro
; ********************************************************************************
    Machine "bbc-tape"
    CPU "6502"
    ORG 0x2000

.start
    LDY #0
.l1 LDA (text),Y
    BEQ l2
    JSR &ffee
    INY
    BNE l1
.l2 RTS

.text
    EQUS "Hello world!", 13, 10, 0
