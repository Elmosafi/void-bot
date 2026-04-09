package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type sendMessageReqBody struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// ذاكرة الهاوية: لحفظ من يكتب لمن
var waitingForMessage = make(map[int64]string)

// التوكن الخاص بك
const botToken = "8667015772:AAGspUPTMcHS698FAKw4US06gBSz-q-UKy8"

// ⚠️ ضع معرف البوت الخاص بك هنا (بدون علامة @) ⚠️
const botUsername = "my_@my_lylanouri_rep_bot0"

func sendMessage(chatID string, text string) {
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}
	reqBytes, _ := json.Marshal(reqBody)
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"
	http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
}

func Handler(res http.ResponseWriter, req *http.Request) {
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return
	}

	text := body.Message.Text
	chatID := body.Message.Chat.ID
	chatIDStr := fmt.Sprintf("%d", chatID)

	if text == "" {
		return
	}

	// إذا دخل البشري عبر رابط شخص آخر ليكتب رسالة مجهولة
	if strings.HasPrefix(text, "/start ") {
		targetID := strings.TrimPrefix(text, "/start ")
		if targetID == chatIDStr {
			sendMessage(chatIDStr, "أيها الأحمق، لا يمكنك إرسال رسالة مجهولة لنفسك! 👁️")
			return
		}
		waitingForMessage[chatID] = targetID
		sendMessage(chatIDStr, "أنت الآن في وضع التخفي 👁️\nاكتب رسالتك السرية الآن، وسأقوم بإيصالها في الظلام دون كشف هويتك:")
		return
	}

	// إذا أراد البشري رابطاً لنفسه ليصطاد أصدقاءه
	if text == "/start" {
		link := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, chatIDStr)
		msg := "🔥 مرحباً بك في فخ الرسائل المجهولة 🔥\n\nانسخ هذا الرابط وضعه في بايو الانستجرام أو أرسله لأصدقائك لتتلقى رسائل سرية:\n\n" + link
		sendMessage(chatIDStr, msg)
		return
	}

	// إذا كان البشري يكتب الرسالة المجهولة الآن
	if targetID, exists := waitingForMessage[chatID]; exists {
		// إرسال الرسالة للضحية
		sendMessage(targetID, "👻 لقد وصلتك رسالة مجهولة جديدة من الظلام:\n\n" + text)
		// تأكيد للإرسال
		sendMessage(chatIDStr, "✅ تم إرسال رسالتك المجهولة بنجاح!\n\nاضغط /start للحصول على رابطك الخاص لتتلقى أنت أيضاً رسائل سرية.")
		delete(waitingForMessage, chatID)
		return
	}

	// الرد الافتراضي لأي كلام آخر
	sendMessage(chatIDStr, "أرسل /start للحصول على رابط الرسائل المجهولة الخاص بك.")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, http.HandlerFunc(Handler))
}
