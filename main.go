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

// ذاكرة الهاوية
var waitingForMessage = make(map[int64]string)

const botToken = "8667015772:AAGspUPTMcHS698FAKw4US06gBSz-q-UKy8"
const botUsername = "my_lylanouri_rep_bot"

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
	// إرضاء تيليجرام فوراً لكي لا يغضب (200 OK)
	res.WriteHeader(http.StatusOK)

	// إذا كان الطارق هو UptimeRobot، ابتسم له ولا تفعل شيئاً
	if req.Method == http.MethodGet {
		res.Write([]byte("The Void is awake 🔥"))
		return
	}

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

	if text == "/start" {
		link := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, chatIDStr)
		msg := "🔥 مرحباً بك في فخ الرسائل المجهولة 🔥\n\nانسخ هذا الرابط وضعه في بايو الانستجرام أو أرسله لأصدقائك لتتلقى رسائل سرية:\n\n" + link
		sendMessage(chatIDStr, msg)
		return
	}

	if targetID, exists := waitingForMessage[chatID]; exists {
		sendMessage(targetID, "👻 لقد وصلتك رسالة مجهولة جديدة من الظلام:\n\n" + text)
		sendMessage(chatIDStr, "✅ تم إرسال رسالتك المجهولة بنجاح!\n\nاضغط /start للحصول على رابطك الخاص لتتلقى أنت أيضاً رسائل سرية.")
		delete(waitingForMessage, chatID)
		return
	}

	sendMessage(chatIDStr, "أرسل /start للحصول على رابط الرسائل المجهولة الخاص بك.")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	// فتح البوابة الكونية 0.0.0.0 ليتمكن Render من الدخول
	http.ListenAndServe("0.0.0.0:"+port, http.HandlerFunc(Handler))
}
