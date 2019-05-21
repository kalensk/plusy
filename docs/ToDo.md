# To Do
 
- add acceptance tests

- complete unit testing 

- integrate travis CI

- implement /topstats

- reverse index for friendly chatroom name from chatId

- exponential retry for getUpdates() and when connecting to the database

- how to change log levels? pass a parameter when starting?

- implement adaptive polling to Telegram
That is, if plusy hasn't seen a message in awhile then poll less often. If more messages come then poll more often.

- goroutine per channel

- catch signals to properly shutdown

- daemonize code
  - use systemd
  - See: 
    - https://github.com/golang/go/issues/227
    - https://stackoverflow.com/questions/10067295/how-to-start-a-go-program-as-a-daemon-in-ubuntu   
    - https://en.wikipedia.org/wiki/Process_supervision
    - https://gunes.io/2017/08/25/systemd-vs-supervisor/
    - https://immortal.run/
    - http://supervisord.org/running.html
    - https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/
    - https://blog.questionable.services/article/running-go-applications-in-the-background/
    - https://gist.github.com/elithrar/9539414

- log to syslog?
  - Is this needed if using a process manager such as systemd or supervisorm?
  - https://github.com/sirupsen/logrus/blob/master/hooks/syslog/syslog.go

- write log to rotating file
  - perhaps this should be handled by an external utility. 

- write my own Google Cloud Function or AWS Lambda

- record who gave points to people and who received points from people

- people can change their names so store, userId to firstNanme:UserName

- watch for name change events and update accordingly

- lazily populate list of names? prune when people leave channel esp important for large channels
// TODO: https://github.com/rubenlagus/TelegramBots/issues/297#issuecomment-327614169


- How best to recover from potential panics?
  - use a top level recover() in main?

- Code Coverage
  - See: https://codecov.io/gh/google/wire
  
- use a fuzzer
  - https://github.com/dvyukov/go-fuzz

- look into mutation testing
  - See: https://github.com/zimmski/go-mutesting

- only log getupdates and message reecived and returning telegram offset if they changed at all.

- implement retries when a transaction fails?
  -  EXEC may fail if user changes....so need to retry entire operation, so create
  - implement something like the following function: commitTransaction(num retries, func() )
  - // func watchingMulti(d *redis.Conn, func(d *redis.Conn, tryNum int) bool, keys string...) X  /// where X is type of whatever EXEC returns


# Libraries

- Logging
  - logrus
  - https://github.com/uber-go/zap
  - https://github.com/sirupsen/logrus
  - https://github.com/inconshreveable/log15
  - https://github.com/jeanphorn/log4go

- Options
  - https://github.com/spf13/viper

- CLI
 - https://github.com/spf13/cobra

- Linter
  - https://github.com/alecthomas/gometalinter

- Testing
  - https://github.com/ory/dockertest
  - https://github.com/stretchr/testify
  - https://github.com/golang/mock
  - https://github.com/cweill/gotests
  - https://github.com/corbym/gocrest
  - https://github.com/seanpont/assert
  - https://github.com/pavlo/gosuite

- Router & Dispatcher
  - https://github.com/gorilla/mux



# Issues

- Cant give plus one to plusy since it is not in the user database
```
time="2018-10-07T16:37:40-07:00" level=debug msg="Incrementing count for user: <nil>"
time="2018-10-07T16:37:40-07:00" level=debug msg="Saving spo timestamp for userId: 0"
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x6518b8]
```

- returning the current plusy count after someone gives a +1 always shows 1 instead of the correct current value
```
plusy:
brain has 1 plusy!

Even though Brian currently has more than 1 plusy
```