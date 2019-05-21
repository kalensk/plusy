package messages

import "time"

type UpdateResponse struct {
	Ok               bool     `json:"ok"`
	ErrorCode        int      `json:error_code,omitempty"`
	ErrorDescription string   `json:"description,omitempty"`
	Result           []Result `json:"result"`
}

type Result struct {
	UpdateID      int64    `json:"update_id"`
	Message       *Message `json:"message,omitempty"`
	EditedMessage *Message `json:"edited_message,omitempty"`
}

type Message struct {
	MessageID      int64           `json:"message_id"`
	From           *User           `json:"from"`
	Date           int64           `json:"date"`
	Chat           Chat            `json:"chat"`
	Text           *string         `json:"text,omitempty"`
	Entities       []MessageEntity `json:"entities,omitempty"`
	ReplyToMessage *ReplyToMessage `json:"reply_to_message,omitempty"`
	NewChatMember  *User           `json:"new_chat_member,omitempty"`
	LeftChatMember *User           `json:"left_chat_member,omitempty"`
	EditDate       *int64          `json:"edit_date,omitempty"` // part of EditMessage
}

type ReplyToMessage struct {
	Message
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID                          int64  `json:"id"`
	Type                        string `json:"type"`
	Title                       string `json:"title,omitempty"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

type MessageEntity struct {
	Type   string  `json:"type"`
	Offset int     `json:"offset"`
	Length int     `json:"length"`
	Url    *string `json:"url,omitempty"`
	User   *User   `json:"user,omitempty"`
}

///////////////////////////////////

type SendMessageRequest struct {
	ChatId    int64  `json:"chat_id"`
	Text      string `json:"text,omitempty"`
	Sticker   string `json:"sticker,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
	Audio     string `json:"audio,omitempty"`
}

type SendMessageResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
	Result      struct {
		Message
	} `json:"result"`
}

///////////////////////////////

func HasText(message *Message) bool {
	return message != nil && message.Text != nil && *message.Text != ""
}

func ConvertMessageToMapForLogging(message *Message) map[string]interface{} {
	chat := map[string]interface{}{
		"title": message.Chat.Title,
		"id":    message.Chat.ID,
	}

	user := map[string]interface{}{
		"id":         message.From.ID,
		"first_name": message.From.FirstName,
		"user_name":  message.From.FirstName,
	}

	loggedMessage := map[string]interface{}{
		"chat": chat,
		"user": user,
		"text": message.Text,
		"date": time.Unix(message.Date, 0),
	}

	return loggedMessage
}

///////////////////////////////////

// New and Left Chat Participants

/*
{
  "ok": true,
  "result": [
    {
      "update_id": 733841785,
      "message": {
        "message_id": 553,
        "from": {
          "id": 260754952,
          "is_bot": false,
          "first_name": "fuzzie",
          "username": "fuzzie_wuzzie",
          "language_code": "en-US"
        },
        "chat": {
          "id": -276219865,
          "title": "Testy Plusy",
          "type": "group",
          "all_members_are_administrators": true
        },
        "date": 1536646115,
        "new_chat_participant": {
          "id": 308188500,
          "is_bot": false,
          "first_name": "brian",
          "last_name": "feaster"
        },
        "new_chat_member": {
          "id": 308188500,
          "is_bot": false,
          "first_name": "brian",
          "last_name": "feaster"
        },
        "new_chat_members": [
          {
            "id": 308188500,
            "is_bot": false,
            "first_name": "brian",
            "last_name": "feaster"
          }
        ]
      }
    }
  ]
}



///////////////////////////////////////

{
  "ok": true,
  "result": [
    {
      "update_id": 733841784,
      "message": {
        "message_id": 552,
        "from": {
          "id": 260754952,
          "is_bot": false,
          "first_name": "fuzzie",
          "username": "fuzzie_wuzzie",
          "language_code": "en-US"
        },
        "chat": {
          "id": -276219865,
          "title": "Testy Plusy",
          "type": "group",
          "all_members_are_administrators": true
        },
        "date": 1536646089,
        "left_chat_participant": {
          "id": 308188500,
          "is_bot": false,
          "first_name": "brian",
          "last_name": "feaster"
        },
        "left_chat_member": {
          "id": 308188500,
          "is_bot": false,
          "first_name": "brian",
          "last_name": "feaster"
        }
      }
    }
  ]
}

*/
