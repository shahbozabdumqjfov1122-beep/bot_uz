package main

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdData struct {
	FileID     string
	Caption    string
	IsVideo    bool
	ButtonText string
	AdLink     string
}

var (
	botToken = "8467228808:AAFECOp5yEOtryP8X5Gk2codpJAWtCE0dp0"
	adminID  = int64(7518992824)

	channels      = make(map[int64]string) // Majburiy obuna kanallari
	userRequests  = make(map[int64]map[int64]bool)
	adminState    = make(map[int64]string)
	tempChannelID = make(map[int64]int64)
	userAdData    = make(map[int64]*AdData)
)

// Majburiy obunani tekshirish
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
	row1 := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📣 Reklama tayyorlash"))
	keyboard := tgbotapi.NewReplyKeyboard(row1)
	if userID == adminID {
		row2 := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("➕ Kanal qo'shish"))
		keyboard.Keyboard = append(keyboard.Keyboard, row2)
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
		// Kanalga qo'shilish so'rovini ushlash
		if update.ChatJoinRequest != nil {
			uID := update.ChatJoinRequest.From.ID
			cID := update.ChatJoinRequest.Chat.ID
			if userRequests[uID] == nil {
				userRequests[uID] = make(map[int64]bool)
			}
			userRequests[uID][cID] = true
			continue
		}

		// Callback tugmalar
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

			if cb.Data == "start_sending" {
				adminState[cb.From.ID] = "wait_target_channel"
				bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🔗 Reklama yubormoqchi bo'lgan kanal linkini yuboring (@kanal_nomi):"))
				bot.Request(tgbotapi.NewCallback(cb.ID, ""))
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		msgText := update.Message.Text
		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		// Majburiy obuna tekshiruvi (Admin uchun istisno)
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

		if msgText == "❌ Bekor qilish" {
			delete(adminState, userID)
			delete(userAdData, userID)
			msg := tgbotapi.NewMessage(chatID, "🚫 Bekor qilindi.")
			msg.ReplyMarkup = getMainMenu(userID)
			bot.Send(msg)
			continue
		}

		// State Machine
		if state, ok := adminState[userID]; ok {
			switch state {
			case "wait_id":
				id, _ := strconv.ParseInt(msgText, 10, 64)
				tempChannelID[userID] = id
				adminState[userID] = "wait_link"
				bot.Send(tgbotapi.NewMessage(chatID, "🔗 Kanal so'rov linkini yuboring:"))
				continue

			case "wait_link":
				channels[tempChannelID[userID]] = msgText
				delete(adminState, userID)
				msg := tgbotapi.NewMessage(chatID, "✅ Kanal majburiy obunaga saqlandi!")
				msg.ReplyMarkup = getMainMenu(userID)
				bot.Send(msg)
				continue

			case "wait_media":
				if update.Message.Photo != nil {
					photos := update.Message.Photo
					userAdData[userID] = &AdData{FileID: photos[len(photos)-1].FileID, IsVideo: false}
					adminState[userID] = "wait_text"
					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Matnni kiriting:"))
				} else if update.Message.Video != nil {
					userAdData[userID] = &AdData{FileID: update.Message.Video.FileID, IsVideo: true}
					adminState[userID] = "wait_text"
					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Matnni kiriting:"))
				}
				continue

			case "wait_text":
				userAdData[userID].Caption = msgText
				adminState[userID] = "wait_btn_text"
				bot.Send(tgbotapi.NewMessage(chatID, "🔘 Tugma matnini kiriting:"))
				continue

			case "wait_btn_text":
				userAdData[userID].ButtonText = msgText
				adminState[userID] = "wait_ad_link"
				bot.Send(tgbotapi.NewMessage(chatID, "🔗 Tugma linkini yuboring:"))
				continue

			case "wait_ad_link":
				userAdData[userID].AdLink = msgText
				data := userAdData[userID]
				btn := tgbotapi.NewInlineKeyboardButtonData("📤 Uzatish", "start_sending")
				keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

				bot.Send(tgbotapi.NewMessage(chatID, "👀 Reklama tayyor! Uni yuborish uchun tugmani bosing:"))
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
				continue

			case "wait_target_channel":
				targetChat := msgText
				if !strings.HasPrefix(targetChat, "@") {
					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Kanal linki @ bilan boshlanishi kerak!"))
					continue
				}

				// Bot adminligini tekshirish
				botMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
					ChatConfigWithUser: tgbotapi.ChatConfigWithUser{SuperGroupUsername: targetChat, UserID: bot.Self.ID},
				})
				if err != nil || (!botMember.IsAdministrator() && !botMember.IsCreator()) {
					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Men u kanalda admin emasman! Avval meni admin qiling."))
					continue
				}

				// Foydalanuvchi adminligini tekshirish
				userMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
					ChatConfigWithUser: tgbotapi.ChatConfigWithUser{SuperGroupUsername: targetChat, UserID: userID},
				})
				if err != nil || (!userMember.IsAdministrator() && !userMember.IsCreator()) {
					bot.Send(tgbotapi.NewMessage(chatID, "❌ Siz bu kanalda admin emassiz! Birovning kanaliga reklama yuborish jinoyat! 🛑"))
					delete(adminState, userID)
					continue
				}

				data := userAdData[userID]
				btn := tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)
				keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

				if data.IsVideo {
					v := tgbotapi.NewVideo(0, tgbotapi.FileID(data.FileID))
					v.BaseChat.ChannelUsername = targetChat
					v.Caption = data.Caption
					v.ReplyMarkup = keyboard
					bot.Send(v)
				} else {
					p := tgbotapi.NewPhoto(0, tgbotapi.FileID(data.FileID))
					p.BaseChat.ChannelUsername = targetChat
					p.Caption = data.Caption
					p.ReplyMarkup = keyboard
					bot.Send(p)
				}

				bot.Send(tgbotapi.NewMessage(chatID, "🚀 Reklama kanalga yuborildi!"))
				delete(adminState, userID)
				delete(userAdData, userID)
				continue
			}
		}

		switch msgText {
		case "/start":
			msg := tgbotapi.NewMessage(chatID, "🚀 Xush kelibsiz!")
			msg.ReplyMarkup = getMainMenu(userID)
			bot.Send(msg)

		case "📣 Reklama tayyorlash":
			adminState[userID] = "wait_media"
			msg := tgbotapi.NewMessage(chatID, "📸 Rasm yoki video yuboring:")
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
