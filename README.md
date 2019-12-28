## craps


### dice rolls
 2 -> 1
 3 -> 2
 4 -> 3
 5 -> 4
 6 -> 5
 7 -> 6
 8 -> 5
 9 -> 4
 10 -> 3
 11 -> 2
 12 -> 1

### win on pass
7, 11
6/36 + 2/36 = 8/36 = 0.222

### lose on pass
2, 3, 12
1/36 + 2/36 + 1/36 = 4/36 = 0.111

### calculating probability of x rolling before y
 pr(x before y) = [pr(x) / pr(x) + pr(y)]

### pass or come odds
 pays 1:1
 comeout win = pr(7) + pr(11) = 8/36
 after comeout = pr(4)×pr(4 before 7) + pr(5)×pr(5 before 7) + pr(6)×pr(6 before 7) + pr(8)×pr(8 before 7) + pr(9)×pr(9 before 7) + pr(10)×pr(10 before 7)
               = (3/36 * 3/9) + (4/36 * 4/10) + 5/36 * 5/11 + 5/36 * 5/11 + 4/36 * 4/10 + 3/36 * 3/9
 the overall probability of winning is 8/36 + 9648/35640 = 17568/35640 = 244/495
 the probability of losing is obviously 1-(244/495) = 251/495
 the player's edge is thus (244/495)×(+1) + (251/495)×(-1) = -7/495 ≈ -1.414%. (house edge)

244/495 * x - 251/495 = 0
244/495 * x = 251/495
244 * x = 251
x = 251/244 (payout if 0% edge)

### buying the odds
 0% house edge
 pays:
  - 2/1 on 4, 10
  - 3/2 on 5, 9
  - 6/5 on 6, 8

 edge = pr(4 before 7) * payout - pr(7 before 4)

 4 and 10: 3/9 × 2/1 + -6/9 = 0.000%
 5 and 9:  4/10 × 3/2 + -6/10 = 0.000%
 6 and 8: 5/11 × 6/5 + -6/11 = 0.000%

### place bet
 place 6 and 8: pays 7/6

 target: 80-100 rolls per hour.
 start at 6pm -> play until 2-3am: 6hours

 6 * 100 = 600 rolls ->
   i want to play 6h with $300
   $300/600 = **50c per roll loss**

### plapayouts

pass/come, 1.41% house edge
  - 1
    
taking the odds, 0% house edge
  - 2/1 on 4, 10
  - 3/2 on 5, 9
  - 6/5 on 6, 8

place bets:
  - 7/6 on 6, 8 with 1.52% house edge


### strategies

- always place 6/8
- have 4 numbers working
- use $30 per shooter on a table min. of $15
