package sender

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	UWDChatID = -1001094145433
	Message   = 0
	Sticker   = 1
)

type Sender struct {
	bot *tgbotapi.BotAPI
}

func (s *Sender) Init(bot *tgbotapi.BotAPI) {
	s.bot = bot
}

func (s *Sender) SendMessageToUWDChat(message string) {
	var reply tgbotapi.MessageConfig
	reply = tgbotapi.NewMessage(
		UWDChatID,
		message,
	)

	_, err := s.bot.Send(reply)

	if err != nil {
		log.Println(err)
	}
}

func (s *Sender) SendSticker(msg *tgbotapi.Message, stickerID string) {
	sticker := tgbotapi.NewStickerShare(
		msg.Chat.ID,
		stickerID,
	)

	_, err := s.bot.Send(sticker)

	if err != nil {
		log.Println(err)
	}
}

func (s Sender) SendReply(msg *tgbotapi.Message, text string) {
	var reply tgbotapi.MessageConfig
	reply = tgbotapi.NewMessage(
		msg.Chat.ID,
		text,
	)

	_, err := s.bot.Send(reply)

	if err != nil {
		log.Println(err)
	}
}

func (s Sender) Send(msgConfig *tgbotapi.MessageConfig) *tgbotapi.Message {
	msg, err := s.bot.Send(msgConfig)
	if err != nil {
		log.Println(err)
	}
	return &msg
}

func (s Sender) SendReplyToMessage(msg *tgbotapi.Message, text string) {
	var reply tgbotapi.MessageConfig
	reply = tgbotapi.NewMessage(
		msg.Chat.ID,
		text,
	)
	reply.ReplyToMessageID = msg.MessageID

	_, err := s.bot.Send(reply)

	if err != nil {
		log.Println(err)
	}
}

func (s Sender) SendMarkdownReply(msg *tgbotapi.Message, text string) {
	var reply tgbotapi.MessageConfig
	reply = tgbotapi.NewMessage(
		msg.Chat.ID,
		text,
	)

	reply.ParseMode = "markdown"
	reply.ReplyToMessageID = msg.MessageID

	_, err := s.bot.Send(reply)

	if err != nil {
		log.Println("Send markdown error", err)
	}
}

func (s Sender) SendInlineKeyboardReply(CallbackQuery *tgbotapi.CallbackQuery, text string) {
	s.bot.AnswerCallbackQuery(tgbotapi.NewCallback(CallbackQuery.ID, text))
}

func (s Sender) EditMessageMarkup(msg *tgbotapi.Message, markup *tgbotapi.InlineKeyboardMarkup) tgbotapi.Message {
	edit := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      msg.Chat.ID,
			MessageID:   msg.MessageID,
			ReplyMarkup: markup,
		},
	}

	message, err := s.bot.Send(edit)

	if err != nil {
		log.Println(err)
	}
	return message
}

func (s *Sender) DeleteMessage(msg *tgbotapi.Message) {
	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    msg.Chat.ID,
		MessageID: msg.MessageID,
	}
	s.bot.DeleteMessage(deleteMsg)
}

func (s Sender) EditMessageText(msg *tgbotapi.Message, text string, parsemode string) tgbotapi.Message {
	edit := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    msg.Chat.ID,
			MessageID: msg.MessageID,
		},
		Text:      text,
		ParseMode: parsemode,
	}

	message, err := s.bot.Send(edit)

	if err != nil {
		log.Println(err)
	}
	return message
}

func (s Sender) SendStickerOrText(msg *tgbotapi.Message, chance int, sending string) {
	switch chance {
	case Sticker:
		s.SendSticker(
			msg,
			sending,
		)
	case Message:
		s.SendReply(
			msg,
			sending,
		)
	}
}
