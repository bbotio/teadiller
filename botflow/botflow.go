package botflow

import (
    "fmt"
    "log"
    "strings"
    "teadiller/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
    "regexp"
)

type Context map[string]interface{}

type Flow struct {
    Nexts []*Flow // list of next steps. So list of flows use here because it depends on context
    Command string // Default handler has this field is nil otherwise command name
    Handler func(updateMessage tgbotapi.Message, requestContext Context) ([]tgbotapi.MessageConfig, error) // hundler function
}

type UserRuntime struct {
    UserContext Context
    UserFlow *Flow
}

// this function bind to the current frow new one and return binded flow
func (f *Flow) Bind(command string, handler func(tgbotapi.Message, Context) ([]tgbotapi.MessageConfig, error)) *Flow {
    bindedFlow := Flow{nil, command, handler}
    f.Nexts = append(f.Nexts, &bindedFlow)
    return &bindedFlow
}

func StartBot(token string, botname string, initFlow Flow, done chan bool)  error {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return err
    }

    runtime := make(map[string]*UserRuntime)

    updateConfiguration := tgbotapi.NewUpdate(0)
    updateConfiguration.Timeout = 60

    updates, err := bot.GetUpdatesChan(updateConfiguration)

    go func() {
        for {
            select {
                case update := <-updates:
                    log.Printf("Got an update %s", update)
                    log.Printf("Message text %s", update.Message.Text)
                    userRuntimeId := fmt.Sprintf("%s_%s", update.Message.From.ID, update.Message.Chat.ID)
                    var userRuntime UserRuntime = UserRuntime{}

                    userRuntimePointer, ok := runtime[userRuntimeId]
                    if (!ok) {
                        userRuntime = UserRuntime{Context{}, &initFlow}
                        runtime[userRuntimeId] = &userRuntime
                        log.Printf("Init runtime %s", userRuntime)
                    } else {
                        userRuntime = *userRuntimePointer
                        log.Printf("Found runtime %s", userRuntime)
                    }

                    text := update.Message.Text
                    if strings.HasPrefix(text, botname) {
                        text = text[len(botname):]
                    }

                    if strings.HasSuffix(text, botname) {
                        text = text[0:len(text) - len(botname)]
                    }
                    text = strings.TrimSpace(text)

                    foundFlow := initFlow
                    nextFlows := userRuntime.UserFlow.Nexts

                    // If flow was ended reset it to init flow
                    if len(nextFlows) == 0 {
                        log.Printf("Reset flow to init")
                        nextFlows = initFlow.Nexts
                    } else {
                        nextFlows = append(nextFlows, initFlow.Nexts...)
                    }

                    log.Printf("Handle the following text '%s'", text)
                    for _, flow := range nextFlows {
                        log.Printf("Check the following flow: %s", flow)
                        if flow.Command == "" { // found default flow
                            log.Printf("Set default flow")
                            foundFlow = *flow
                        } else if match,_ := regexp.MatchString(flow.Command, text); match { // found command flow
                            log.Printf("Found command handler %s", flow)
                            foundFlow = *flow
                            break
                        }
                    }

                    responses, err := foundFlow.Handler(update.Message, userRuntime.UserContext)
                    if err != nil {
                        errorMessage := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
                        bot.Send(errorMessage)
                    } else {
                        *userRuntime.UserFlow = foundFlow
                        for _, response := range responses {
                            bot.Send(response)
                        }
                    }
                case <- done:
                    log.Printf("Handling of bot updates were stoped")
                    return
            }
        }
    }()

    return nil
}
