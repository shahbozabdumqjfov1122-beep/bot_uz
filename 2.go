package main

//
//import (
//	"log"
//	"strconv"
//	"strings"
//
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//)
//
//// Reklama ma'lumotlari uchun struktura
//type AdData struct {
//	FileID     string
//	Caption    string
//	IsVideo    bool
//	ButtonText string
//	AdLink     string // Tugma havolasini saqlash uchun
//}
//
//var (
//	botToken = "8534860816:AAHybGqTACVQ48gFG5fKBkxEhBtDHBSRid0"
//	adminID  = int64(7518992824)
//	channels = make(map[int64]string)
//
//	userRequests  = make(map[int64]map[int64]bool)
//	adminState    = make(map[int64]string)
//	tempChannelID = make(map[int64]int64)
//	userAdData    = make(map[int64]*AdData)
//	allUsers      = make(map[int64]bool) // Bot foydalanuvchilarini saqlash uchun
//)
//
//func hasSentAllRequests(userID int64) bool {
//	if len(channels) == 0 {
//		return true
//	}
//	for cID := range channels {
//		if userRequests[userID] == nil || !userRequests[userID][cID] {
//			return false
//		}
//	}
//	return true
//}
//
//func getMainMenu(userID int64) tgbotapi.ReplyKeyboardMarkup {
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📣 Reklama tayyorlash")),
//	)
//	if userID == adminID {
//		keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("➕ Kanal qo'shish")))
//	}
//	return keyboard
//}
//
//func getCancelMenu() tgbotapi.ReplyKeyboardMarkup {
//	return tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("❌ Bekor qilish")),
//	)
//}
//
//func main() {
//	bot, err := tgbotapi.NewBotAPI(botToken)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//	log.Println("Bot ishga tushdi...")
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//	updates := bot.GetUpdatesChan(u)
//
//	for update := range updates {
//		// Har qanday yangi foydalanuvchini bazaga (mapga) qo'shish
//		if update.Message != nil {
//			allUsers[update.Message.From.ID] = true
//		}
//
//		// Kanalga qo'shilish so'rovlarini ushlash
//		if update.ChatJoinRequest != nil {
//			uID := update.ChatJoinRequest.From.ID
//			cID := update.ChatJoinRequest.Chat.ID
//			if userRequests[uID] == nil {
//				userRequests[uID] = make(map[int64]bool)
//			}
//			userRequests[uID][cID] = true
//			continue
//		}
//
//		// CallbackQuery: Tekshirish va Reklama tarqatish
//		if update.CallbackQuery != nil {
//			cb := update.CallbackQuery
//
//			// 1. Kanal tanlash tugmasi bosilganda
//			if cb.Data == "choose_channel" {
//				if len(channels) == 0 {
//					bot.Request(tgbotapi.NewCallbackWithAlert(cb.ID, "❌ Hali kanallar qo'shilmagan!"))
//					continue
//				}
//
//				var rows [][]tgbotapi.InlineKeyboardButton
//				for id, _ := range channels {
//					// Har bir kanal uchun alohida tugma
//					btn := tgbotapi.NewInlineKeyboardButtonData("Kanal ID: "+strconv.FormatInt(id, 10), "send_to_"+strconv.FormatInt(id, 10))
//					rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
//				}
//
//				msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "Qaysi kanalga yubormoqchisiz?")
//				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
//				bot.Send(msg)
//			}
//
//			// 2. Tanlangan kanalga reklamani yuborish
//			if strings.HasPrefix(cb.Data, "send_to_") {
//				targetChannelStr := strings.TrimPrefix(cb.Data, "send_to_")
//				targetChannelID, _ := strconv.ParseInt(targetChannelStr, 10, 64)
//
//				data := userAdData[cb.From.ID]
//				if data == nil {
//					bot.Request(tgbotapi.NewCallbackWithAlert(cb.ID, "❌ Reklama topilmadi! Qaytadan tayyorlang."))
//					continue
//				}
//
//				btn := tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)
//				keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))
//
//				if data.IsVideo {
//					v := tgbotapi.NewVideo(targetChannelID, tgbotapi.FileID(data.FileID))
//					v.Caption = data.Caption
//					v.ReplyMarkup = keyboard
//					bot.Send(v)
//				} else {
//					p := tgbotapi.NewPhoto(targetChannelID, tgbotapi.FileID(data.FileID))
//					p.Caption = data.Caption
//					p.ReplyMarkup = keyboard
//					bot.Send(p)
//				}
//
//				bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "🚀 Reklama kanalga muvaffaqiyatli yuborildi!"))
//				bot.Request(tgbotapi.NewCallback(cb.ID, "Yuborildi"))
//			}
//		}
//		if update.Message == nil {
//			continue
//		}
//
//		msgText := update.Message.Text
//		chatID := update.Message.Chat.ID
//		userID := update.Message.From.ID
//
//		// Majburiy obuna (Admin uchun istisno)
//		if userID != adminID && !hasSentAllRequests(userID) {
//			var rows [][]tgbotapi.InlineKeyboardButton
//			for _, link := range channels {
//				btn := tgbotapi.NewInlineKeyboardButtonURL("📩 So'rov yuborish", link)
//				rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
//			}
//			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔄 Tekshirish", "check_sub")))
//			msg := tgbotapi.NewMessage(chatID, "🤖 Botdan foydalanish uchun kanallarga so'rov yuboring!")
//			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
//			bot.Send(msg)
//			continue
//		}
//
//		if msgText == "❌ Bekor qilish" {
//			delete(adminState, userID)
//			delete(userAdData, userID)
//			msg := tgbotapi.NewMessage(chatID, "🚫 Jarayon bekor qilindi.")
//			msg.ReplyMarkup = getMainMenu(userID)
//			bot.Send(msg)
//			continue
//		}
//
//		// Admin Bosqichlari (State Machine)
//		if state, ok := adminState[userID]; ok {
//			switch state {
//			case "wait_id":
//				id, _ := strconv.ParseInt(msgText, 10, 64)
//				tempChannelID[userID] = id
//				adminState[userID] = "wait_link"
//				bot.Send(tgbotapi.NewMessage(chatID, "🔗 So'rov linkini yuboring:"))
//				continue
//
//			case "wait_link":
//				channels[tempChannelID[userID]] = msgText
//				delete(adminState, userID)
//				msg := tgbotapi.NewMessage(chatID, "✅ Kanal saqlandi!")
//				msg.ReplyMarkup = getMainMenu(userID)
//				bot.Send(msg)
//				continue
//
//			case "wait_media":
//				if update.Message.Photo != nil {
//					photos := update.Message.Photo
//					userAdData[userID] = &AdData{FileID: photos[len(photos)-1].FileID, IsVideo: false}
//					adminState[userID] = "wait_text"
//					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Reklama matnini (izoh) kiriting:"))
//				} else if update.Message.Video != nil {
//					userAdData[userID] = &AdData{FileID: update.Message.Video.FileID, IsVideo: true}
//					adminState[userID] = "wait_text"
//					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Reklama matnini (izoh) kiriting:"))
//				} else {
//					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Iltimos, rasm yoki video yuboring!"))
//				}
//				continue
//
//			case "wait_text":
//				userAdData[userID].Caption = msgText
//				adminState[userID] = "wait_btn_text"
//				bot.Send(tgbotapi.NewMessage(chatID, "🔘 Tugma matnini kiriting:"))
//				continue
//
//			case "wait_btn_text":
//				userAdData[userID].ButtonText = msgText
//				adminState[userID] = "wait_ad_link"
//				bot.Send(tgbotapi.NewMessage(chatID, "🔗 Tugma uchun link (URL) yuboring:"))
//				continue
//
//			case "wait_ad_link":
//				userAdData[userID].AdLink = msgText
//				data := userAdData[userID]
//
//				// "Uzatish" tugmasini qo'shish
//				btnLink := tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)
//				btnShare := tgbotapi.NewInlineKeyboardButtonData("📤 Uzatish (Kanalni tanlash)", "choose_channel")
//
//				keyboard := tgbotapi.NewInlineKeyboardMarkup(
//					tgbotapi.NewInlineKeyboardRow(btnLink),
//					tgbotapi.NewInlineKeyboardRow(btnShare),
//				)
//
//				bot.Send(tgbotapi.NewMessage(chatID, "✅ Reklama tayyor! Qaysi kanalga yuborishni tanlash uchun 'Uzatish' tugmasini bosing."))
//
//				if data.IsVideo {
//					v := tgbotapi.NewVideo(chatID, tgbotapi.FileID(data.FileID))
//					v.Caption = data.Caption
//					v.ReplyMarkup = keyboard
//					bot.Send(v)
//				} else {
//					p := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(data.FileID))
//					p.Caption = data.Caption
//					p.ReplyMarkup = keyboard
//					bot.Send(p)
//				}
//				continue
//			}
//		}
//
//		// Menyu buyruqlari
//		switch msgText {
//		case "/start":
//			msg := tgbotapi.NewMessage(chatID, "🚀 Xush kelibsiz!")
//			msg.ReplyMarkup = getMainMenu(userID)
//			bot.Send(msg)
//
//		case "📣 Reklama tayyorlash":
//			adminState[userID] = "wait_media"
//			msg := tgbotapi.NewMessage(chatID, "📸 Rasm yoki 🎬 video yuboring:")
//			msg.ReplyMarkup = getCancelMenu()
//			bot.Send(msg)
//
//		case "➕ Kanal qo'shish":
//			if userID == adminID {
//				adminState[userID] = "wait_id"
//				msg := tgbotapi.NewMessage(chatID, "🆔 Kanal ID raqamini kiriting:")
//				msg.ReplyMarkup = getCancelMenu()
//				bot.Send(msg)
//			}
//		}
//	}
//}
