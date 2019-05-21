# Plusy

plusy is an [inline Telegram bot](https://core.telegram.org/bots/inline) that records and provides statistics when someone in a chat `+1`'s an individual.
It is a simple project to help me better learn golang and redis. 

# Using Plusy
If you do not want to run your own plusy server


# Running Plusy
To run your own plusy server

1. start the database
  - `make database type=redis` other types are postgres or neo4j 
1. start plusy
 - `make install`
 



Command: `@plusy stats`

Sample Output:
```
Points Given: 3
Top Given To: bob 1, alice 2

Received: 2 pts
Top Receivers: Dan 1, bob 1
```

Command: `@plusy help`
Sample Output:
```
@plusy stats 
```


- Database Scheme
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


# Gif of showing it work in action

# Features
- commands
- redis database pooling
- responds to both inline queries and replied messages
- allows for querying over both username and firstname
- daemonize with systemd
- can run plusy in a docker container...
- unit test and acceptance tests using docker

# Overall Design
GetUpdates is a stream of messages from all chat rooms




# Challenges


## Message vs. EditedMessage vs. GenericMessage Interfacce etc.

// incoming result can either be a message or EditedMessage.
// find a better way of accepting a generic "message" that is either type
func (d *Database) SaveLastMessage(message *messages.Message) {

## Logging API token
What are the current thoughts on logging api tokens?

man during a failed http.Get while getting telegram updates. Im currently just panic'ing cuz I haven't really implemented anything else. Anyway I found out that it logs this
panic: Get https://api.telegram.org/bot<token>/getUpdates?timeout=2&offset=28839812: dial tcp: lookup api.telegram.org on 127.0.1.1:53: read udp 127.0.0.1:46463->127.0.1.1:53: i/o timeout

put all your log lines through a regex or equivalent to strip anything that looks like secrets.

## Issue

" level=debug msg="Processing topN stats command for chatId: -1001082930701\n"
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x652cf8]

goroutine 1 [running]:
gj!One should not be allowed to give +1's to plusy bot

fuzzieadmin
Ok problem is that I am giving points to user plusy and such even though "user plusy" is not in my database as a user
On every single incoming message I do saveOrRemoveUser(message). But turns out plusy doesn't see its own messages so it never records itself as a user
le sigh



## Issue Acceptance Testing


Was doing something like the following, but it was jank and not extensible. Wanted to say first call to getUpdates return X
and second call to getUpdates return Y. And for sendMessage then return K. Something readable, extensible, clean, etc.

```go

var docker *testutils.Docker
var database *redisdb.Redis
var log *logrus.Logger

func TestMain(m *testing.M) {
	docker = testutils.NewRedisDocker(m)
	log = docker.Log
	database = docker.SetupTestRedis()
	returnCode := docker.Run()
	docker.TearDown(returnCode)
}

func getFakeResponseBytes(fileName string) []byte {
	_, acceptanceTestFileName, _, _ := runtime.Caller(0)
	responseFile := filepath.Join(filepath.Dir(acceptanceTestFileName), fileName)
	fakeResponseBytes, err := ioutil.ReadFile(responseFile)
	if err != nil {
		panic(err)
	}

	//var update messages.UpdateResponse
	//err = json.Unmarshal(fakeResponseBytes, &update)

	return fakeResponseBytes
}

func Test_Acceptance(t *testing.T) {
	var command *exec.Cmd

	fakeResponse1 := getFakeResponseBytes("response1.json")
	//fakeResponse2 := getFakeResponseBytes("response2.json")

	var count int
	testServer := httptest.NewServer(http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("Inside testServer")
		fmt.Println(request.Method)
		fmt.Println(request.RequestURI)
		fmt.Println(request.Body)
	
		if strings.Contains(request.URL.Path, "getUpdates") {
			respWriter.WriteHeader(200)
			respWriter.Header().Add("Content-Type", "application/json")
			respWriter.Write(fakeResponse1)
			if count > 0 {
				command.Process.Kill()
			}
			count++
		}
	
		if strings.Contains(request.URL.Path, "sendMessage") {
			requestBytes, err := ioutil.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}
	
			fmt.Println(string(requestBytes))
		}
	
	}))
	defer testServer.Close()

	plusyBinary := filepath.Join(os.Getenv("GOPATH"), "bin", "plusy")
	command = exec.Command(plusyBinary, "--token", "123",
		"--telegram-api-url", testServer.URL+"/bot",
		"--database-host-port", docker.GetDatabaseConnectionString())
	// todo: pass in a flag (-f) to write to stdout/stderr instead of plusy.log

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	command.Stdout = &stdOut
	command.Stderr = &stdErr
	err := command.Start()
	duration, err := time.ParseDuration("10m")
	time.Sleep(duration)
	if err != nil {
		log.Println(err)
		log.Println(stdOut.String())
		log.Println(stdErr.String())
	}
}

```
 