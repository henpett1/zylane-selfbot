//go:generate goversioninfo -icon="./icon.ico"
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	strconv "strconv"
	"strings"
	"syscall"
	"time"

	"github.com/admin100/util/console"
	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
)

func config(thing_from_config string) interface{} {
	bytevalue, _ := ioutil.ReadFile("./config.json")
	var result map[string]interface{}
	json.Unmarshal(bytevalue, &result)
	return result[thing_from_config]
}

func get_delay(thing_from_file string) time.Duration {
	bytevalue, _ := ioutil.ReadFile("./config.json")
	var result map[string]time.Duration
	json.Unmarshal(bytevalue, &result)
	return result[thing_from_file]
}

func logo() {
	color.HEX("#65EAB5").Println("███████╗██╗   ██╗██╗      █████╗ ███╗   ██╗███████╗")
	color.HEX("#66E2B9").Println("╚══███╔╝╚██╗ ██╔╝██║     ██╔══██╗████╗  ██║██╔════╝")
	color.HEX("#67D9BC").Println("  ███╔╝  ╚████╔╝ ██║     ███████║██╔██╗ ██║█████╗  ")
	color.HEX("#67D1C0").Println(" ███╔╝    ╚██╔╝  ██║     ██╔══██║██║╚██╗██║██╔══╝  ")
	color.HEX("#68C8C3").Println("███████╗   ██║   ███████╗██║  ██║██║ ╚████║███████╗")
	color.HEX("#69C0C7").Println("╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝")
}

func createconfig() {

	f, err := os.Create("config.json")

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{\n	\"token\": \"\",\n	\"prefix\": \"!\", \"delay:\" 15\n}")

	if err2 != nil {
		fmt.Println(err2)
	}

	fmt.Println("Made config. Fill it in and start the bot")
}

func disclog() {
	token := config("token").(string)
	discord, err := discordgo.New(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Print("\033[H\033[2J")
	fmt.Println(" ")
	nbServers := len(discord.State.Guilds)
	strconv.Itoa(nbServers)
	nbFriends := len(discord.State.Relationships)
	strconv.Itoa(nbFriends)
	logo()
	fmt.Println("[-] Zylane selfbot is running!")
	fmt.Println("[-] Logged in as: ", discord.State.User.Username, "#", discord.State.User.Discriminator)
	console.SetConsoleTitle(fmt.Sprintf("Zylane selfbot | Logged in as %s#%s", discord.State.User.Username, discord.State.User.Discriminator))
	fmt.Println("[-] Guilds:", nbServers)
	fmt.Println("[-] Friends:", nbFriends)
	fmt.Println(" ")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	discord.Close()
}

func main() {
	fmt.Print("\033[H\033[2J")
	logo()
	file, err := os.OpenFile("config.json", os.O_RDWR, 0644)
	if errors.Is(err, os.ErrNotExist) {
		createconfig()
		defer file.Close()
	} else {
		defer file.Close()
		fmt.Println("config.json exists, starting bot")
		disclog()
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID != s.State.User.ID {
		return
	}
	delaytime := get_delay("delay")
	prefix := config("prefix").(string)
	argslic := m.Content[len(prefix):]
	args := strings.Split(argslic, " ")
	command := strings.ToLower(args[0])
	switch command {
	case "help":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		msg, err := s.ChannelMessageSend(m.ChannelID, "```yaml\n"+
			"Zylane selfbot\n"+
			"\n"+
			"Commands:\n"+
			"\n"+
			prefix+"help - Shows this message\n"+
			prefix+"ping - Pong\n"+
			prefix+"restart - Restarts the bot\n"+
			prefix+"shutdown - Shuts down the bot\n"+
			prefix+"dog - Sends a random image of a dog\n"+
			prefix+"cat - Sends a random image of a cat\n"+
			prefix+"raccon - Sends a random image of a racoon\n"+
			prefix+"redpanda - Sends a random image of a red panda\n"+
			prefix+"fox - Sends a random image of a fox\n"+
			prefix+"msgspam <ammount> <text> - Spams a message\n"+
			"```\n")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("[-] Used command: help")
		time.Sleep(delaytime * time.Second)
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
	case "ping":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		msg, err := s.ChannelMessageSend(m.ChannelID, "```yaml\n"+
			"Zylane selfbot\n"+
			"Pong!\n"+
			"```\n")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("[-] Used command: ping")
		time.Sleep(delaytime * time.Second)
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
	case "restart":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, "```yaml\n"+
			"Zylane selfbot\n"+
			"Restarting...\n"+
			"```\n")
		fmt.Println("[-] Used command: restart")
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		cmd := exec.Command(fmt.Sprintf("start %s/Zylane.exe", exPath))
		erro := cmd.Run()
		fmt.Println(erro)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	case "shutdown":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, "```yaml\n"+
			"Zylane selfbot\n"+
			"Shutting down...\n"+
			"```\n")

		fmt.Println("[-] Used command: shutdown")
		os.Exit(0)
	case "dog":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		link := "https://random.dog/woof.json"
		response, err := http.Get(link)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		bytevalue, _ := ioutil.ReadAll(response.Body)
		var result map[string]interface{}
		json.Unmarshal(bytevalue, &result)

		msg, err := s.ChannelMessageSend(m.ChannelID, result["url"].(string))
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Println("[-] Used command: dog")
		time.Sleep(delaytime * time.Second)
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
	case "cat":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		link := "https://aws.random.cat/meow"
		response, err := http.Get(link)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		bytevalue, _ := ioutil.ReadAll(response.Body)
		var result map[string]interface{}
		json.Unmarshal(bytevalue, &result)
		msg, err := s.ChannelMessageSend(m.ChannelID, result["file"].(string))
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Println("[-] Used command: cat")
		time.Sleep(delaytime * time.Second)
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
	case "fox":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		link := "https://randomfox.ca/floof/"
		response, err := http.Get(link)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		bytevalue, _ := ioutil.ReadAll(response.Body)
		var result map[string]interface{}
		json.Unmarshal(bytevalue, &result)
		msg, err := s.ChannelMessageSend(m.ChannelID, result["image"].(string))
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Println("[-] Used command: fox")
		time.Sleep(delaytime * time.Second)
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
	case "raccoon":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		link := "https://some-random-api.ml/animal/raccoon"
		response, err := http.Get(link)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		bytevalue, _ := ioutil.ReadAll(response.Body)
		var result map[string]interface{}
		json.Unmarshal(bytevalue, &result)
		msg, err := s.ChannelMessageSend(m.ChannelID, result["image"].(string))
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("[-] Used command: raccoon")
		time.Sleep(delaytime * time.Second)
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
	case "redpanda":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		link := "https://some-random-api.ml/animal/red_panda"
		response, err := http.Get(link)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		bytevalue, _ := ioutil.ReadAll(response.Body)
		var result map[string]interface{}
		json.Unmarshal(bytevalue, &result)
		msg, err := s.ChannelMessageSend(m.ChannelID, result["image"].(string))
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("[-] Used command: redpanda")
		time.Sleep(delaytime * time.Second)
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
	case "msgspam":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		amm, err := strconv.Atoi(args[1])
		fmt.Println("[-] Used command: msgspam")
		if err != nil {
			fmt.Println(err.Error())
		}
		msg := strings.Trim(m.Content, prefix+"msgspam"+args[1]+" ")
		for i := 0; i < amm; i++ {
			s.ChannelMessageSend(m.ChannelID, msg)
		}

	}
}
