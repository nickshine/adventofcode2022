#1/usr/bin/env bash


# A -> rock       1
# B -> paper      2
# C -> scissors   3
#
# X -> lose -> 0
# Y -> draw -> 3
# Z -> win  -> 6
#
# A Y  = draw, rock draws to rock ->  A -> 3 (draw) + 1 (rock) = 4
# B X  = lose rock loses to paper ->  A -> 0 (lose) + 1 (rock) = 1
# C Z  = win, rock beats scissors ->  A -> 6 (win)  + 1 (rock) = 7
# total                                                        = 12 
#
# lose cases: AC, BA, CB
# draw cases: AA, BB, CC
# win cases:  AB, BC, CA

gsed \
    -e 's/A X/0+C/g' \
    -e 's/B X/0+A/g' \
    -e 's/C X/0+B/g' \
    -e 's/A Y/3+A/g' \
    -e 's/B Y/3+B/g' \
    -e 's/C Y/3+C/g' \
    -e 's/A Z/6+B/g' \
    -e 's/B Z/6+C/g' \
    -e 's/C Z/6+A/g' \
    input.txt \
    | tr 'ABC' '123' | bc | paste -sd+ - | bc

