; ********************************************************************************
; This is a simple Hello World example for the BBC micro
; ********************************************************************************
;   Machine "bbc-tape"
    CPU "6502"
    ORG 0x2000

.start
    LDY #0
.l1 LDA (text),Y
    BEQ l2
    JSR &ffee
    INY
    BNE l1
l2  RTS

.text
    EQUS "Hello world!", 13, 10, 0

brktest1
    BRK           ; Software break
    EQUB 1        ; Error code 0
    EQUS "Silly"  ; This is a real error message in BBC BASIC, try: AUTO10,1000
    EQUB 0        ; End of message marker

brktest2
    BRK  2       ; Software break Error code 2
    EQUS "Silly"  ; This is a real error message in BBC BASIC, try: AUTO10,1000
    EQUB 0        ; End of message marker

brktest2
    BRK  #3       ; Software break Error code 3
    EQUS "Silly"  ; This is a real error message in BBC BASIC, try: AUTO10,1000
    EQUB 0        ; End of message marker
