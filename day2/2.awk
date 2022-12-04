#!/usr/bin/awk -f

/Y/ { s += 3 }
/Z/ { s += 6 }

/A Y|B X|C Z/ { s += 1 }
/B Y|C X|A Z/ { s += 2 }
/C Y|A X|B Z/ { s += 3 }

END { print s }

