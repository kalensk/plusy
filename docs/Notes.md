# Notes

 
# Style Guide
- log messages are capitalized, but error messages are lower case


## Process Supervision using systemd

Update plusy.service and to see effects
```
sudo cp plusy.service /etc/systemd/system
sudo systemctl daemon-reload 
sudo systemctl enable plusy.service
sudo systemctl start plusy
sudo systemctl status plusy
ps auxwf | grep 

```

## Creating Mocks
- `mockery -name=<MyThingInterface>`
- Rename file to mock_mything.go, and rename the actual struct it created to MockMyThing
- Move it out of the mocks package. 
- See: https://github.com/vektra/mockery


## Creating a Bot
- Talk to the BotFather
- "Bot Settings" > "Group Privacy" > "Turn off"
- "Edit Bot" > "Edit Commands"
```
help - get help
stats - get stats for a specified user, or the whole channel if no user is specified
topstats - get top stats for the channel

```

## Telegram

`curl -H "Content-Type: application/json" https://api.telegram.org/bot<token>/sendMessage -X POST -d '{ "chat_id": "-276219865", "text": "donut" }'`
`curl -H "Content-Type: application/json" https://api.telegram.org/bot<token>/sendSticker -X POST -d '{ "chat_id": "-276219865", "sticker": "CAADAQADlxcAAh1poQfdatXu101EtAI" }'`

sending mp3 to Telegram
`curl -X POST -H "Content-Type: multipart/form-data" -F audio=@ah.mp3  "https://api.telegram.org/bot<token>/sendAudio?chat_id=<chatId>"`


```
Messages, commands and requests sent by users are passed to the software running on your servers. Our intermediary server handles all encryption and communication with the Telegram API for you. You communicate with this server via a simple HTTPS-interface that offers a simplified version of the Telegram API. We call that interface our Bot API.
```
- [Bot API](https://core.telegram.org/bots/api)
- [Long Pooling](https://core.telegram.org/bots/faq#long-polling-gives-me-the-same-updates-again-and-again) 
- [Bot FAQ](https://core.telegram.org/bots/faq#how-do-i-get-updates)
- [Bot Privacy Mode](https://core.telegram.org/bots#privacy-mode)
```
Privacy mode is enabled by default for all bots, except bots that were added to the group as admins (bot admins always receive all messages). It can be disabled, so that the bot receives all messages like an ordinary user. We only recommend doing this in cases where it is absolutely necessary for your bot to work — users can always see a bot's current privacy setting in the group members list. In most cases, using the force reply option for the bot's messages should be more than enough.

...

2. Bot admins and bots with privacy mode disabled will receive all messages except messages sent by other bots.

```
