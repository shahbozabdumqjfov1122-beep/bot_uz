package main

//
//import (
//	"fmt"
//	"log"
//	"strconv"
//	"strings"
//
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//)
//
//// AdData reklama ma'lumotlarini saqlash uchun
//type AdData struct {
//	FileID     string
//	Caption    string
//	IsVideo    bool
//	HasMedia   bool
//	ButtonText string
//	AdLink     string
//}
//
//var (
//	// TOKENNI YANGILANG!
//	botToken   = "8534860816:AAEH3QSbf9bj5vr4ARG7tbusvC70WpZgdqY"
//	adminState = make(map[int64]string)
//	userAdData = make(map[int64]*AdData)
//	// Kanallar ro'yxati (Aslida DB ishlatish kerak, hozircha mapda)
//	autoApproveChannels = make(map[int64]bool)
//)
//
//func main() {
//	bot, err := tgbotapi.NewBotAPI(botToken)
//	if err != nil {
//		log.Panic("Botni ishga tushirishda xato:", err)
//	}
//	bot.Debug = false // Loglar juda ko'payib ketmasligi uchun false qildim
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//	updates := bot.GetUpdatesChan(u)
//
//	log.Println("✅ Bot muvaffaqiyatli ishga tushdi!")
//
//	for update := range updates {
//		// 1. Avto-qabul so'rovlarini ushlash
//		if update.ChatJoinRequest != nil {
//			HandleAutoApprove(bot, update.ChatJoinRequest)
//			continue
//		}
//
//		// 2. Tugmalarni ushlash (Callback)
//		if update.CallbackQuery != nil {
//			handleCallback(bot, update)
//			continue
//		}
//
//		// 3. Xabarlarni ushlash
//		if update.Message == nil {
//			continue
//		}
//
//		handleMessage(bot, update)
//	}
//}
//
//func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//	msg := update.Message
//	chatID := msg.Chat.ID
//	userID := msg.From.ID
//	text := msg.Text
//
//	// Har doim "Bekor qilish"ni tekshirish
//	if text == "❌ Bekor qilish" {
//		resetUserState(bot, chatID, userID)
//		return
//	}
//
//	// Admin holatiga qarab ishlash
//	if state, ok := adminState[userID]; ok {
//		switch state {
//		case "wait_accept_channel":
//			SetupJoinRequest(bot, update)
//			return
//		case "wait_media":
//			handleMediaInput(bot, update)
//			return
//		case "wait_text":
//			if userAdData[userID] == nil {
//				userAdData[userID] = &AdData{}
//			}
//			userAdData[userID].Caption = text
//			adminState[userID] = "wait_btn_text"
//			bot.Send(tgbotapi.NewMessage(chatID, "⚙️ Tugma matnini kiriting:"))
//			return
//		case "wait_btn_text":
//			userAdData[userID].ButtonText = text
//			adminState[userID] = "wait_ad_link"
//			bot.Send(tgbotapi.NewMessage(chatID, "🔗 Tugma linkini yuboring (https://...):"))
//			return
//		case "wait_ad_link":
//			if !strings.HasPrefix(text, "http") {
//				bot.Send(tgbotapi.NewMessage(chatID, "⚠️ To'g'ri link yuboring (http bilan boshlansin):"))
//				return
//			}
//			userAdData[userID].AdLink = text
//			sendPreview(bot, chatID, userID)
//			return
//		case "wait_target_channel":
//			FinalizeAdSending(bot, update)
//			return
//		}
//	}
//
//	// Asosiy buyruqlar
//	switch text {
//	case "/start":
//		sendWelcomeMessage(bot, chatID)
//	case "📣 Reklama tayyorlash":
//		startAdCreation(bot, chatID, userID)
//	case "🔄 Avto-qabulni sozlash":
//		startAutoApproveSetup(bot, chatID, userID)
//	}
//}
//
//func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//	cb := update.CallbackQuery
//	data := cb.Data
//	userID := cb.From.ID
//	chatID := cb.Message.Chat.ID
//
//	switch {
//	case data == "start_sending":
//		adminState[userID] = "wait_target_channel"
//		bot.Send(tgbotapi.NewMessage(chatID, "🔗 Reklama yubormoqchi bo'lgan kanal ID yoki @username yuboring:"))
//		bot.Request(tgbotapi.NewCallback(cb.ID, ""))
//	case data == "cancel_accept":
//		delete(adminState, userID)
//		bot.Send(tgbotapi.NewEditMessageText(chatID, cb.Message.MessageID, "🚫 Bekor qilindi."))
//	case strings.HasPrefix(data, "start_accept_"):
//		channelIDStr := strings.TrimPrefix(data, "start_accept_")
//		channelID, _ := strconv.ParseInt(channelIDStr, 10, 64)
//		autoApproveChannels[channelID] = true // Xotiraga saqlash
//		bot.Send(tgbotapi.NewEditMessageText(chatID, cb.Message.MessageID, fmt.Sprintf("✅ Kanal (%d) uchun avto-qabul yoqildi!", channelID)))
//		delete(adminState, userID)
//	}
//}
//
//func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Xush kelibsiz! Bot orqali reklama tayyorlash va kanallarga avto-qabulni sozlash mumkin.")
//	msg.ReplyMarkup = getMainMenu()
//	bot.Send(msg)
//}
//
//func startAdCreation(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
//	adminState[userID] = "wait_media"
//	msg := tgbotapi.NewMessage(chatID, "📸 Rasm yoki 📹 video yuboring:")
//	msg.ReplyMarkup = getMediaMenu()
//	bot.Send(msg)
//}
//
//func startAutoApproveSetup(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
//	adminState[userID] = "wait_accept_channel"
//	text := "🔄 **Avto-qabulni sozlash:**\n\n1. Botni kanalingizga admin qiling.\n2. Kanaldan biror xabarni shu yerga forward qiling."
//	msg := tgbotapi.NewMessage(chatID, text)
//	msg.ParseMode = "Markdown"
//	msg.ReplyMarkup = getCancelMenu()
//	bot.Send(msg)
//}
//
//func handleMediaInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//	msg := update.Message
//	userID := msg.From.ID
//	chatID := msg.Chat.ID
//
//	userAdData[userID] = &AdData{}
//	if msg.Photo != nil {
//		userAdData[userID].FileID = msg.Photo[len(msg.Photo)-1].FileID
//		userAdData[userID].HasMedia = true
//		userAdData[userID].IsVideo = false
//	} else if msg.Video != nil {
//		userAdData[userID].FileID = msg.Video.FileID
//		userAdData[userID].HasMedia = true
//		userAdData[userID].IsVideo = true
//	} else if msg.Text == "⏭ Tashlab ketish" {
//		userAdData[userID].HasMedia = false
//	} else {
//		bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Iltimos, rasm yoki video yuboring."))
//		return
//	}
//
//	adminState[userID] = "wait_text"
//	resp := tgbotapi.NewMessage(chatID, "✍️ Reklama matnini kiriting:")
//	resp.ReplyMarkup = getCancelMenu()
//	bot.Send(resp)
//}
//
//func SetupJoinRequest(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//	msg := update.Message
//	var targetChatID int64
//
//	if msg.ForwardFromChat != nil {
//		targetChatID = msg.ForwardFromChat.ID
//	} else {
//		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "⚠️ Iltimos, kanaldan xabarni forward qiling!"))
//		return
//	}
//
//	keyboard := tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("🚀 Yoqish", fmt.Sprintf("start_accept_%d", targetChatID)),
//			tgbotapi.NewInlineKeyboardButtonData("❌ Bekor qilish", "cancel_accept"),
//		),
//	)
//
//	resp := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("📡 Kanal aniqlandi: %s\nID: %d\nAvto-qabulni yoqamizmi?", msg.ForwardFromChat.Title, targetChatID))
//	resp.ReplyMarkup = keyboard
//	bot.Send(resp)
//}
//
//func FinalizeAdSending(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//	msg := update.Message
//	targetChat := msg.Text // Foydalanuvchi yuborgan @kanal_nomi yoki ChatID
//	userID := msg.From.ID
//
//	data := userAdData[userID]
//	btn := tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)
//	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))
//
//	var config tgbotapi.Chattable
//
//	if !data.HasMedia {
//		// Oddiy matnli xabar
//		m := tgbotapi.NewMessageToChannel(targetChat, data.Caption)
//		m.ReplyMarkup = keyboard
//		config = m
//	} else if data.IsVideo {
//		// Video yuborish (Kanalga yuborish uchun NewVideo va BaseChat ishlatiladi)
//		v := tgbotapi.NewVideo(0, tgbotapi.FileID(data.FileID))
//		v.BaseChat.ChannelUsername = targetChat    // @username bo'lsa
//		if strings.HasPrefix(targetChat, "-100") { // Agar ID bo'lsa
//			id, _ := strconv.ParseInt(targetChat, 10, 64)
//			v.BaseChat.ChatID = id
//			v.BaseChat.ChannelUsername = "" // ID ishlatilganda username bo'sh bo'lishi kerak
//		}
//		v.Caption = data.Caption
//		v.ReplyMarkup = keyboard
//		config = v
//	} else {
//		// Rasm yuborish
//		p := tgbotapi.NewPhoto(0, tgbotapi.FileID(data.FileID))
//		p.BaseChat.ChannelUsername = targetChat
//		if strings.HasPrefix(targetChat, "-100") {
//			id, _ := strconv.ParseInt(targetChat, 10, 64)
//			p.BaseChat.ChatID = id
//			p.BaseChat.ChannelUsername = ""
//		}
//		p.Caption = data.Caption
//		p.ReplyMarkup = keyboard
//		config = p
//	}
//
//	_, err := bot.Send(config)
//	if err != nil {
//		log.Printf("Yuborishda xato: %v", err)
//		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Xato: Bot kanalda adminmi yoki kanal linki to'g'rimi?"))
//	} else {
//		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "🚀 Reklama muvaffaqiyatli yuborildi!"))
//	}
//	resetUserState(bot, msg.Chat.ID, userID)
//}
//func sendPreview(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
//	data := userAdData[userID]
//	keyboard := tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📤 Kanalga yuborish", "start_sending")),
//	)
//
//	bot.Send(tgbotapi.NewMessage(chatID, "👀 Reklama ko'rinishi:"))
//	if !data.HasMedia {
//		m := tgbotapi.NewMessage(chatID, data.Caption)
//		m.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)))
//		bot.Send(m)
//	} else if data.IsVideo {
//		v := tgbotapi.NewVideo(chatID, tgbotapi.FileID(data.FileID))
//		v.Caption = data.Caption
//		v.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)))
//		bot.Send(v)
//	} else {
//		p := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(data.FileID))
//		p.Caption = data.Caption
//		p.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)))
//		bot.Send(p)
//	}
//
//	finalMsg := tgbotapi.NewMessage(chatID, "Tayyor bo'lsa yuborish tugmasini bosing:")
//	finalMsg.ReplyMarkup = keyboard
//	bot.Send(finalMsg)
//}
//
//func HandleAutoApprove(bot *tgbotapi.BotAPI, request *tgbotapi.ChatJoinRequest) {
//	channelID := request.Chat.ID
//	userID := request.From.ID
//
//	// Faqat sozlangan kanallar uchun qabul qilish
//	if !autoApproveChannels[channelID] {
//		log.Printf("⚠️ Kanal (%d) uchun avto-qabul yoqilmagan", channelID)
//		return
//	}
//
//	approve := tgbotapi.ApproveChatJoinRequestConfig{
//		ChatConfig: tgbotapi.ChatConfig{ChatID: channelID},
//		UserID:     userID,
//	}
//
//	_, err := bot.Request(approve)
//	if err != nil {
//		log.Printf("❌ Qabul qilishda xato: %v", err)
//	} else {
//		log.Printf("✅ Tasdiqlandi: UserID %d, Kanal %d", userID, channelID)
//	}
//}
//
//func resetUserState(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
//	delete(adminState, userID)
//	delete(userAdData, userID)
//	msg := tgbotapi.NewMessage(chatID, "Bosh menyu")
//	msg.ReplyMarkup = getMainMenu()
//	bot.Send(msg)
//}
//
//func getMainMenu() tgbotapi.ReplyKeyboardMarkup {
//	markup := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📣 Reklama tayyorlash")),
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("🔄 Avto-qabulni sozlash")),
//	)
//	markup.ResizeKeyboard = true
//	return markup
//}
//
//func getCancelMenu() tgbotapi.ReplyKeyboardMarkup {
//	markup := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("❌ Bekor qilish")),
//	)
//	markup.ResizeKeyboard = true
//	return markup
//}
//
//func getMediaMenu() tgbotapi.ReplyKeyboardMarkup {
//	markup := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("⏭ Tashlab ketish")),
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("❌ Bekor qilish")),
//	)
//	markup.ResizeKeyboard = true
//	return markup
//}
