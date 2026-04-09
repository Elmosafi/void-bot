package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
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
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func Handler(res http.ResponseWriter, req *http.Request) {
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return
	}

	// إذا لم تكن هناك رسالة نصية، لا تفعل شيئاً
	if body.Message.Text == "" {
		return
	}

	// الرد الشيطاني الذي سيرسله البوت
	replyText := "أنا مستيقظ يا سيدي... لقد استدعيتني من العدم! 🔥\nلقد قلت لي: " + body.Message.Text

	reqBody := &sendMessageReqBody{
		ChatID: body.Message.Chat.ID,
		Text:   replyText,
	}

	reqBytes, _ := json.Marshal(reqBody)

	// التوكن الخاص بك مزروع هنا في قلب الظلام
	token := "8667015772:AAGspUPTMcHS698FAKw4US06gBSz-q-UKy8"
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
}

func main() {
	// إجبار الخادم على فتح البوابة الصحيحة
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, http.HandlerFunc(Handler))
}
