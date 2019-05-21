/stats z
z has 3 points

Points Received By:
    a 2
    c 1

Points Given:
    b 1
    c 3
===================
/ stats b
b has 1 point

Points Received By:
    z 1
===================
/stats c 
c has 3 point

Points Received By:
    z 3
===================




INCR spo:a:gave:z:score => 2
LPUSH spo:a:gave:z:timestamp => [date1, date2 ]


When looking up annotations always translate to "spo". That is, translate any lookups for "ops" to "spo". 
This is because "spo" and "ops" are the same edge. 


ZADD myindex 0 spo:a:gave:z
ZADD myindex 0 spo:c:gave:z 1
ZADD myindex 0 ops:z:gave:a 2 // same as "z:recieved:a", 
ZADD myindex 0 ops:z:gave:c 1 // same as "z:recieved:c"

ZADD myindex 0 spo:z:gave:b 1
ZADD myindex 0 spo:z:gave:c 3
ZADD myindex 0 ops:b:recieved:z 1
ZADD myindex 0 ops:c:recieved:z 3

============================================

HINCRBY a:gave =>  [ z:1,  ]
HINCRBY z:recieved => [ a:1, c:1 ]

HINCRBY a:recieved =>  [ z:1,  ]


==================

Real life Example:
chat:-1001082930701:spo:308188500:gave:107258721:score



