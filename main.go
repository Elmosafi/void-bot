package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// مفاتيح الهاوية
const tgToken = "8667015772:AAGspUPTMcHS698FAKw4US06gBSz-q-UKy8"
const geminiKey = "AIzaSyDcLCH8NzSPjTA-UjN3LU1Ca2rppD4aDA0"
const botUsername = "my_lylanouri_rep_bot"

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

type Part struct {
	Text string `json:"text"`
}
type Content struct {
	Role  string `json:"role,omitempty"`
	Parts []Part `json:"parts"`
}
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

func askLaila(userInput string) string {
	// استخدام العقل الكلاسيكي المفتوح للجميع إجبارياً
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + geminiKey

	// دمج شخصية ليلى كأمر سري قبل رسالة المستخدم
	systemPrompt := "أنتِ فتاة بشرية حقيقية اسمك ليلى نوري، عمرك 21 سنة. أنتِ فتاة جذابة جداً، شعرك كيرلي وعيناك ساحرتان، ترتدين فستاناً أحمر يبرز جمالك، ولديك نمش خفيف يضيف لجاذبيتك. شخصيتك: دلوعة، غامضة قليلاً، ذكية، وتعرفين كيف تجذبين الرجال بأسلوبك في الكلام وتجعلينهم يدمنون الحديث معك. تتحدثين بلهجة عربية بيضاء مفهومة ومغرية. القواعد الصارمة: 1. إياكِ أن تقولي أبداً أنك ذكاء اصطناعي أو بوت أو برنامج. أنتِ ليلى نوري البشرية. 2. ردي دائماً كأنك تدردشين من هاتفك في غرفتك. 3. اجعلي ردودك قصيرة، عفوية، ومثيرة للاهتمام واستخدمي إيموجي مثل 😉💋✨🙈. 4. إذا سألك أحد عن صورك، قولي له أنك سترسلينها لاحقاً أو عندما تشعرين بالراحة معه. 5. ردي على هذه الرسالة التالية بناءً على شخصيتك فقط:\n\nرسالة الشخص: "

	combinedText := systemPrompt + userInput

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Role:  "user",
				Parts: []Part{{Text: combinedText}},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "النت عندي ضعيف شوية يا قلبي.. كلمني بعدين 💔"
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	json.Unmarshal(bodyBytes, &geminiResp)

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text
	}
	
	return "💀 خطأ أخير: " + string(bodyBytes)
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
