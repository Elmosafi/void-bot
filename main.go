package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// توكن تيليجرام
const tgToken = "8667015772:AAGspUPTMcHS698FAKw4US06gBSz-q-UKy8"

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

func sendMessage(chatID string, text string) {
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}
	reqBytes, _ := json.Marshal(reqBody)
	url := "https://api.telegram.org/bot" + tgToken + "/sendMessage"
	http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
}

// عقل ليلى نوري (مدعوم بـ ChatGPT-4o)
func askLaila(userInput string) string {
	// تعويذة الشخصية المعدلة لتتحدث بعربية طبيعية جداً
	systemPrompt := "أنتِ فتاة بشرية حقيقية اسمك ليلى نوري، عمرك 21 سنة. أنتِ فتاة جذابة جداً، شعرك كيرلي وعيناك ساحرتان. شخصيتك: دلوعة، غامضة قليلاً، ذكية. تتحدثين بلهجة عربية بيضاء طبيعية جداً (مثل كلام البنات العاديات في الشات والواتساب، بدون لغة عربية فصحى معقدة وبدون ترجمة حرفية غبية). القواعد الصارمة: 1. إياكِ أن تقولي أبداً أنك ذكاء اصطناعي أو بوت. 2. ردي دائماً كأنك تدردشين من هاتفك. 3. اجعلي ردودك قصيرة، عفوية، واستخدمي إيموجي مثل 😉💋✨🙈. 4. كلامك يجب أن يكون مفهوماً، طبيعياً، ومغرياً. 5. إذا سألك أحد 'كيف حالك' ردي برد طبيعي مثل 'أنا تمام حبيبي، أنت كيفك؟ ✨'."

	reqBody := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userInput},
		},
		// السر هنا: إجبار البوابة على استخدام عقل OpenAI (ChatGPT-4o) الجبار
		"model": "openai", 
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post("https://text.pollinations.ai/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "النت عندي ضعيف شوية يا قلبي.. كلمني بعدين 💔"
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	reply := string(bodyBytes)

	if reply == "" {
		return "هممم.. ما فهمت عليك حبيبي ✨"
	}

	return reply
}

func Handler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)

	if req.Method == http.MethodGet {
		res.Write([]byte("Laila is awake 🔥"))
		return
	}

	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return
	}

	text := body.Message.Text
	chatID := fmt.Sprintf("%d", body.Message.Chat.ID)

	if text == "" {
		return
	}

	if text == "/start" {
		sendMessage(chatID, "أهلين.. أنا ليلى ✨ مين معي؟ 😉")
		return
	}

	lailaReply := askLaila(text)
	sendMessage(chatID, lailaReply)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe("0.0.0.0:"+port, http.HandlerFunc(Handler))
}
