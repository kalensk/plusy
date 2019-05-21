# Interestings

## 
- how to add plusy as a user intially since people can +1 plusy...Or don't allow plus one on plusy



## Allowing Inline Plus One's
Two scenarios to consider: 1) quoting someone, 2) an inline reply

1) quoting
```
bob
|alice:
| i like turtles
+1
```

2) an inline reply
```
alice: i like turtles
bob: +1
```
However, this may result in the following issue if another person says something between the time someone says something interesting and someone gives them a plus one.
```
alice: i like turtles
dan: what did you do yesterday?
bob: +1
```



And the plusy table would look something like:  
giver |  receiver | msg_text | msg_datetime | 
---|---|--- 
bob | alice | "i like turtles" 




## Queries
telegram only guarnatees userId as being unique
Problem

// Inverted Index
// chatid => [userid-DougC, userId-DougB, userId-Brian, userID-Tokie]
// firstname:doug => [userid-DougC, userid-DougB]
// return scores for both DougC and DougB since they are in => [userid-DougC, userid-DougB] and its an ambigous search term
//

in a set storing:
user:firstname:<firstname> => [userid, userid]
user:username:<username> => [userid, userid]

Hash:
user:id:<userid> => firstname: <firstname> username: <username> , etc.
// TODO: this is not unique across all of telegram. There can be multiple users with the same username/firstName, so we need a more unique way to store it.


One cool thing about composite indexes is that they are handy in order to represent graphs, using a data structure which is called Hexastore.

The hexastore provides a representation for relations between objects, formed by a subject, a predicate and an object. A simple relation between objects could be:

