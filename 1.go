package main

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdData struct {
	FileID  string
	Caption string
	IsVideo bool
}

var (
	botToken       = "8467228808:AAFECOp5yEOtryP8X5Gk2codpJAWtCE0dp0"
	adminID  int64 = 7518992824
	channels       = make(map[int64]string)

	userRequests  = make(map[int64]map[int64]bool)
	adminState    = make(map[int64]string)
	tempChannelID = make(map[int64]int64)
	userAdData    = make(map[int64]*AdData)
)

func hasSentAllRequests(userID int64) bool {
	if len(channels) == 0 {
		return true
	}
	for cID := range channels {
		if userRequests[userID] == nil || !userRequests[userID][cID] {
			return false
		}
	}
	return true
}

func getMainMenu(userID int64) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📣 Reklama tayyorlash")),
	)
	if userID == adminID {
		keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("➕ Kanal qo'shish")))
	}
	return keyboard
}

func getCancelMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("❌ Bekor qilish")),
	)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Println("Bot ishga tushdi...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// 1. Kanalga qo'shilish so'rovi
		if update.ChatJoinRequest != nil {
			uID := update.ChatJoinRequest.From.ID
			cID := update.ChatJoinRequest.Chat.ID
			if userRequests[uID] == nil {
				userRequests[uID] = make(map[int64]bool)
			}
			userRequests[uID][cID] = true
			continue
		}

		// 2. Callback (Tekshirish tugmasi)
		if update.CallbackQuery != nil {
			cb := update.CallbackQuery
			if cb.Data == "check_sub" {
				if hasSentAllRequests(cb.From.ID) {
					bot.Send(tgbotapi.NewDeleteMessage(cb.Message.Chat.ID, cb.Message.MessageID))
					bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "✅ So'rovlar aniqlandi! Endi /start bosing."))
				} else {
					bot.Request(tgbotapi.NewCallbackWithAlert(cb.ID, "❌ Siz hali barcha kanallarga so'rov yubormagansiz!"))
				}
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		msgText := update.Message.Text
		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		// 3. Majburiy obuna tekshiruvi (faqat oddiy foydalanuvchilar uchun)
		if userID != adminID && !hasSentAllRequests(userID) {
			var rows [][]tgbotapi.InlineKeyboardButton
			for _, link := range channels {
				btn := tgbotapi.NewInlineKeyboardButtonURL("📩 So'rov yuborish", link)
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔄 Tekshirish", "check_sub")))
			msg := tgbotapi.NewMessage(chatID, "🤖 Botdan foydalanish uchun kanallarga so'rov yuboring!")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
			bot.Send(msg)
			continue
		}

		// 4. Bekor qilish buyrug'i
		if msgText == "❌ Bekor qilish" {
			delete(adminState, userID)
			delete(userAdData, userID)
			msg := tgbotapi.NewMessage(chatID, "🚫 Jarayon bekor qilindi.")
			msg.ReplyMarkup = getMainMenu(userID)
			bot.Send(msg)
			continue
		}

		// 5. Admin State Mantiqi
		if state, ok := adminState[userID]; ok {
			switch state {
			case "wait_id":
				id, _ := strconv.ParseInt(msgText, 10, 64)
				tempChannelID[userID] = id
				adminState[userID] = "wait_link"
				bot.Send(tgbotapi.NewMessage(chatID, "🔗 So'rov linkini yuboring:"))
				continue

			case "wait_link":
				channels[tempChannelID[userID]] = msgText
				delete(adminState, userID)
				msg := tgbotapi.NewMessage(chatID, "✅ Kanal saqlandi!")
				msg.ReplyMarkup = getMainMenu(userID)
				bot.Send(msg)
				continue

			case "wait_media":
				if update.Message.Photo != nil {
					photos := update.Message.Photo
					userAdData[userID] = &AdData{FileID: photos[len(photos)-1].FileID, IsVideo: false}
					adminState[userID] = "wait_text"
					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Reklama matnini (izoh) kiriting:"))
				} else if update.Message.Video != nil {
					userAdData[userID] = &AdData{FileID: update.Message.Video.FileID, IsVideo: true}
					adminState[userID] = "wait_text"
					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Reklama matnini (izoh) kiriting:"))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Iltimos, rasm yoki video yuboring!"))
				}
				continue

			case "wait_text":
				if msgText == "" {
					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Iltimos, reklama matnini yozing:"))
					continue
				}
				userAdData[userID].Caption = msgText
				adminState[userID] = "wait_ad_link"
				bot.Send(tgbotapi.NewMessage(chatID, "🔗 Tugma uchun havola (link) yuboring:"))
				continue

			case "wait_ad_link":
				if msgText == "" {
					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Iltimos, link yuboring:"))
					continue
				}
				adLink := msgText
				data := userAdData[userID]

				// REKLAMANI YUBORISH
				btn := tgbotapi.NewInlineKeyboardButtonURL("Tomosha qilish", adLink)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

				if data.IsVideo {
					v := tgbotapi.NewVideo(chatID, tgbotapi.FileID(data.FileID))
					v.Caption = data.Caption
					v.ReplyMarkup = keyboard
					bot.Send(v)
				} else {
					p := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(data.FileID))
					p.Caption = data.Caption
					p.ReplyMarkup = keyboard
					bot.Send(p)
				}

				delete(adminState, userID)
				delete(userAdData, userID) // Ma'lumotni tozalash
				msg := tgbotapi.NewMessage(chatID, "✅ Reklamangiz tayyor!")
				msg.ReplyMarkup = getMainMenu(userID)
				bot.Send(msg)
				continue
			}
		}

		// 6. Asosiy buyruqlar
		switch msgText {
		case "/start":
			msg := tgbotapi.NewMessage(chatID, "🚀 Xush kelibsiz!")
			msg.ReplyMarkup = getMainMenu(userID)
			bot.Send(msg)

		case "📣 Reklama tayyorlash":
			adminState[userID] = "wait_media"
			msg := tgbotapi.NewMessage(chatID, "📸 Rasm yoki 🎬 video yuboring:")
			msg.ReplyMarkup = getCancelMenu()
			bot.Send(msg)

		case "➕ Kanal qo'shish":
			if userID == adminID {
				adminState[userID] = "wait_id"
				msg := tgbotapi.NewMessage(chatID, "🆔 Kanal ID raqamini kiriting:")
				msg.ReplyMarkup = getCancelMenu()
				bot.Send(msg)
			}
		}
	}
}
