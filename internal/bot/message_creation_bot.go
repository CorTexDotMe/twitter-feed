package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"twitter-feed/internal/model"

	lorelai "github.com/UltiRequiem/lorelai/pkg"
)

func SendPostMessageRequests() {
	serverUrl := os.Getenv("SERVER_URL")
	sendingTimeot, err := strconv.Atoi(os.Getenv("SENDING_TIMEOUT_SECONDS"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		message := model.Message{
			Username: lorelai.Word(),
			Message:  lorelai.Sentence(),
		}
		body, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(sendingTimeot) * time.Second)
			continue
		}

		response, err := http.Post(serverUrl, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(sendingTimeot) * time.Second)
			continue
		}
		fmt.Println(response)

		time.Sleep(time.Duration(sendingTimeot) * time.Second)
	}
}
