#1/usr/bin/env bash


# A,X -> rock       1
# B,Y -> paper      2
# C,Z -> scissors   3
#
# lose -> 0
# draw -> 3
# win  -> 6
#
# A Y  = 2 (paper) + 6 (win)     = 8
# B X  = 1 (rock) + 0 (lose)     = 1
# C Z  = 3 (scissors) + 3 (draw) = 6
# total                          = 15 

gsed \
    -e 's/ /+/g' \
    -e 's/\(A+Y\|B+Z\|C+X\)/6+\1/g' \
    -e 's/\(A+Z\|B+X\|C+Y\)/0+\1/g' \
    -e 's/\(A+X\|B+Y\|C+Z\)/3+\1/g' \
    input.txt \
    | tr 'ABCXYZ' '000123' \
    | bc | paste -sd+ - | bc

