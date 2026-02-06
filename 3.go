package main

//
//import (
//	"log"
//	"strings"
//
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//)
//
//type AdData struct {
//	FileID     string
//	Caption    string
//	IsVideo    bool
//	ButtonText string
//	AdLink     string
//}
//
//var (
//	botToken = "8467228808:AAEpZ2_Qt825dnOPNIYBoqAy8OGfeK2Kyrc"
//	adminID  = int64(7518992824)
//
//	adminState = make(map[int64]string)
//	userAdData = make(map[int64]*AdData)
//	channels   = make(map[int64]string)
//)
//
//func main() {
//	bot, err := tgbotapi.NewBotAPI(botToken)
//	if err != nil {
//		log.Panic(err)
//	}
//	bot.Debug = true
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//	updates := bot.GetUpdatesChan(u)
//
//	for update := range updates {
//		if update.CallbackQuery != nil {
//			cb := update.CallbackQuery
//			if cb.Data == "start_sending" {
//				adminState[cb.From.ID] = "wait_target_channel"
//				msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "🔗 Reklama yubormoqchi bo'lgan kanal linkini yuboring (masalan: @kanal_nomi):")
//				bot.Send(msg)
//				bot.Request(tgbotapi.NewCallback(cb.ID, ""))
//			}
//			continue
//		}
//
//		if update.Message == nil {
//			continue
//		}
//
//		msgText := update.Message.Text
//		chatID := update.Message.Chat.ID
//		userID := update.Message.From.ID
//
//		if msgText == "❌ Bekor qilish" {
//			delete(adminState, userID)
//			msg := tgbotapi.NewMessage(chatID, "🚫 Bekor qilindi.")
//			msg.ReplyMarkup = getMainMenu(userID)
//			bot.Send(msg)
//			continue
//		}
//
//		if state, ok := adminState[userID]; ok {
//			switch state {
//			case "wait_media":
//				if update.Message.Photo != nil {
//					photos := update.Message.Photo
//					userAdData[userID] = &AdData{FileID: photos[len(photos)-1].FileID, IsVideo: false}
//					adminState[userID] = "wait_text"
//					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Matnni kiriting:"))
//				} else if update.Message.Video != nil {
//					userAdData[userID] = &AdData{FileID: update.Message.Video.FileID, IsVideo: true}
//					adminState[userID] = "wait_text"
//					bot.Send(tgbotapi.NewMessage(chatID, "✍️ Matnni kiriting:"))
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
//				bot.Send(tgbotapi.NewMessage(chatID, "🔗 Tugma linkini yuboring:"))
//				continue
//
//			case "wait_ad_link":
//				userAdData[userID].AdLink = msgText
//				data := userAdData[userID]
//				btn := tgbotapi.NewInlineKeyboardButtonData("📤 Uzatish", "start_sending")
//				keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))
//
//				bot.Send(tgbotapi.NewMessage(chatID, "👀 Reklama tayyor! Uni yuborish uchun tugmani bosing:"))
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
//				delete(adminState, userID)
//				continue
//
//			case "wait_target_channel":
//				targetChat := msgText
//				if !strings.HasPrefix(targetChat, "@") {
//					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Kanal linki @ bilan boshlanishi kerak!"))
//					continue
//				}
//
//				botMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
//					ChatConfigWithUser: tgbotapi.ChatConfigWithUser{SuperGroupUsername: targetChat, UserID: bot.Self.ID},
//				})
//				if err != nil || (!botMember.IsAdministrator() && !botMember.IsCreator()) {
//					bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Men u kanalda admin emasman! Avval meni admin qiling."))
//					continue
//				}
//
//				userMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
//					ChatConfigWithUser: tgbotapi.ChatConfigWithUser{SuperGroupUsername: targetChat, UserID: userID},
//				})
//				if err != nil || (!userMember.IsAdministrator() && !userMember.IsCreator()) {
//					bot.Send(tgbotapi.NewMessage(chatID, "❌ Siz bu kanalda admin emassiz! Birovning kanaliga reklama yuborish jinoyat! 🛑"))
//					delete(adminState, userID)
//					continue
//				}
//
//				data := userAdData[userID]
//				btn := tgbotapi.NewInlineKeyboardButtonURL(data.ButtonText, data.AdLink)
//				keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))
//
//				if data.IsVideo {
//					v := tgbotapi.NewVideo(0, tgbotapi.FileID(data.FileID))
//					v.BaseChat.ChannelUsername = targetChat
//					v.Caption = data.Caption
//					v.ReplyMarkup = keyboard
//					bot.Send(v)
//				} else {
//					p := tgbotapi.NewPhoto(0, tgbotapi.FileID(data.FileID))
//					p.BaseChat.ChannelUsername = targetChat
//					p.Caption = data.Caption
//					p.ReplyMarkup = keyboard
//					bot.Send(p)
//				}
//
//				bot.Send(tgbotapi.NewMessage(chatID, "🚀 Reklama kanalga yuborildi!"))
//				delete(adminState, userID)
//				delete(userAdData, userID)
//				continue
//			}
//		}
//
//		switch msgText {
//		case "/start":
//			msg := tgbotapi.NewMessage(chatID, "Xush kelibsiz!")
//			msg.ReplyMarkup = getMainMenu(userID)
//			bot.Send(msg)
//		case "📣 Reklama tayyorlash":
//			adminState[userID] = "wait_media"
//			msg := tgbotapi.NewMessage(chatID, "📸 Rasm yoki video yuboring:")
//			msg.ReplyMarkup = getCancelMenu()
//			bot.Send(msg)
//		}
//	}
//}
//
//func getMainMenu(userID int64) tgbotapi.ReplyKeyboardMarkup {
//	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📣 Reklama tayyorlash")))
//}
//
//func getCancelMenu() tgbotapi.ReplyKeyboardMarkup {
//	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("❌ Bekor qilish")))
//}
