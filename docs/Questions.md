# Questions

1) 

My database New() Nmethod does the following:
```
buffer := bytes.NewBuffer([]byte{})
encoder := gob.NewEncoder(buffer)
decoder := gob.NewDecoder(buffer)
```
return &Database{conn: connection, timeout: timeout, buffer: buffer, encoder: encoder, decoder: decoder}
Where I am only initializing a NewEconcoder, NewDecorder, etc. once in my constructor a-la Java style. However, I am currently getting the error in GetUser() gob: unknown type id or corrupted data and other such shenanigans.


And my SaveUser() method is using those as follows:
```
d.buffer.Reset()
 if err := d.encoder.Encode(&user); err != nil {
  panic(err)
 }

d.conn.Do("SET", "somekey", d.buffer.Bytes())

And my GetUser() method is similar:
var user messages.User
 d.buffer.Write(reply)
 defer d.buffer.Reset()
 err = d.decoder.Decode(&user)
 if err != nil {
  panic(err)
 }

 return &user
```
If I create a new Decoder, Encoder, and Buffer every time in my SaveUser() and GetUser() it all works great. Halp.


A1) Gob outputs stateful data for efficiency :(

---

2) Why does does it not make sense to have a function take in a poitner to an interface. For example.
`func New(database *db.Database, telegramClient *telegram.TelegramInterface) *Plusy`
vs.
`func New(database *db.Database, telegramClient telegram.TelegramInterface) *Plusy`
 

A2)


 ---
 
 3) Given plusy and telegram, what is the difference between using gob vs. JSON for serializing data to redis?
 Gob ties me to using golang, whereas JSON is more universal.
 
 A3)
 
 ---
 
 4) Whats the proper way to do error handeling?
 
 A4)
 
 ---
 
 5)  p.client.SomeFunc() does not work in the below example, but p.database.SomeFunc() does? Why? It is because client 
 and database are interfaces and using a pointer to an interface does not work because ....? 
 ```
 
type Plusy struct {
	log      *logrus.Logger
	database db.Database
	client   *client.ChatClient
}

func New(log *logrus.Logger, database db.Database, client *client.ChatClient) *Plusy {
	return &Plusy{log: log, database: database, client: client}
}

```

5A)

---

6) 
Why exactly does this work, but not the following? That is, why does it not work with Parse has a receiver?
Im assuming I would need an instance of Options in main in order to do that...

Works:
```
package options
func Parse() *Options {
    // ...
}



package main

options := options.Parse()
	
```

Does not work:
```
package options

func (o *Options) Parse() *Options {
    // ...
}



package main

options := options.Parse()
	
```



6A)

---

7Q)

IncrementCount does a bunch of database calls and such, and I want to wrap that in a redis Multi/Exec to make it transactional. Problem is that for my logging I basically do a database user lookup to print a feindly log message. That is within the Multi/Exec transaction block which is no good, no good

Basically Im not sure how to do a database getUser call for my friendly log message from wihtin a Multi/Exec transaction


7A)
You need two connections to redis (one for the txn, one where inconsistent reads are performed), or to change the scope of your transactional operation.
Or you could defer logging until after the exec succeeds.


---





to READ

http://www.funcmain.com/gob_encoding_an_interface
https://blog.golang.org/laws-of-reflection
https://research.swtch.com/interfaces


https://stackoverflow.com/questions/13511203/why-cant-i-assign-a-struct-to-an-interface/13511853#13511853


https://stackoverflow.com/questions/23148812/whats-the-meaning-of-interface/23148998#23148998