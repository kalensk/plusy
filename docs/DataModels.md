# Data Models

Below is a description on how the data is stored for the various types of databases: in-memory, graph, and relational. 

The following example with users A, B, and C is used to explain how plusy data is stored.

This can be illustrated in the following table, where the first line can be read as user A received two points from user B, and received one point from C.
User A gave B one point, and three points to C.

| User | Points Received (User Points) | Points Given (User Points) |
|------|-------------------------------|--------------------------------|
|A     | B 2, C 1                      | B 1, C 3                       |
|B     | A 1                           | A 2                            |
|C     | A 3                           | A 1                            |


```
/stats A
A has 3 points

Points Received By:
    B 2
    C 1

Points Given:
    B 1
    C 3

```

```
/stats B
B has 1 point

Points Received By:
    A 1

Points Given:
    A 2 
```

```
/stats C 
C has 3 point

Points Received By:
    A 3
    
Points Given:
    A 1
```

---

## In-Memory

Redis is used for the in-memory database. Using Lexicographical indexes as introduced in the [Redis documentation](https://redis.io/topics/indexes).

We use the notation `spo` and `ops` which stand for subject:predicate:object and object:predicate:subject, respectively.
We also prepend each key with a "namespace" such as chat, user, and plusy.
 

| Key | Value | Definition |
|-----|-------|------------|
| plusy:offset | telegram offset number as int64 | the last offset of the messages that plusy received from Telegram. |
| chat:lastMsg:<chatId> | telegram message struct | last Telegram message received for a particular Telegram chat. This is used to find the previous user to give a plus one to when only one update message was received while  parsing an inline message. |
| chat:<chatId>:spo:<giver_userId>:gave:<receiver_userId>:score | points given | number of points giver_userId gave receiver_userId |
| chat:<chatId>:ops:<receiver_userId>:gave:<giver_userId>:score | points recieved | number of points receiver_userId gave giver_userId | 
| chat:<chatId>:spo:<giver_userId>:gave:<receiver_userId>:ts | list of timestamps up to a maximum limit | Used to query timestamps of points given. Used to store timestamp annotations on the node. |
| chat:<chatId>:ops:<receiver_userId>:gave:<giver_userId>:ts | list of timestamps up to a maximum limit | Used to query timestamps of points given. Used to store timestamp annotations on the node. |
| chat:keys:<chatId> | sorted set of lexicographical keys spo:<giver_userId>:gave:<receiver_userId> | Used for the lexicographical lookup of who gave points to whom |
| chat:keys:<chatId> | sorted set of lexicographical keys ops:<receiver_userId>:gave:<giver_userId> | Used for the lexicographical lookup of who received points from whom | |
| chat:top:<chatId> | sorted set of top 10 users by points  | Used to quickly returned top 10 users in a chat |
| user:id:<userId> | user struct | a user struct includes a the userId, firstname, lastname, username and if they are a bot |
| user:un:<username> | sorted set of userId's | Used to query stats for a user by username. An inverted index of username to all possible userId's since usernames may not be unique for a given userId. |
| user:fn:<firstname> | sorted set of userId's | Used to query stats for a user by firstname. An inverted index of firstnames to all possible userId's since firstnames may not be unique for a given userId. |


---

## Graph

Neo4j is used as the graph database.


---

## Relational

Postgres is used as the graph database.


