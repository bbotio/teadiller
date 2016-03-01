package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"teadiller/Godeps/_workspace/src/github.com/jasonwinn/geocoder"
	"teadiller/botflow"
	"teadiller/excel"
	"teadiller/flows"
	"teadiller/models"
	"teadiller/payment"
	"teadiller/web"
)

func main() {

	// Init dao
	excelItemDao := excel.ExcelItemDao{FilePath: os.Getenv("TELEGRAM_BOT_DB_PATH")}
	excelOrderDao := excel.ExcelOrderDao{FilePath: os.Getenv("TELEGRAM_BOT_DB_PATH")}

	excelItemDao.ParseItems()
	excelOrderDao.ParseOrders()

	models.InitItemDao(excelItemDao)
	models.InitOrderDao(excelOrderDao)

	mode, _ := strconv.ParseBool(os.Getenv("PAYPAL_SANDBOX"))
	payment.InitClient(
		os.Getenv("PAYPAL_USERNAME"), os.Getenv("PAYPAL_PASSWORD"),
		os.Getenv("PAYPAL_SIGNATURE"), mode,
		os.Getenv("PAYPAL_OK_URL"), os.Getenv("PAYPAL_CANCEL_URL"))
	geocoder.SetAPIKey(os.Getenv("GEOCODER_API_KEY"))

	initFlow := botflow.Flow{Command: "", Handler: flows.Default}
	initFlow.Bind("/about", flows.About)
	categoriesFlow := initFlow.Bind("/categories", flows.Categories)
	showFlow := categoriesFlow.Bind("/show", flows.Show)

	buyFlow := showFlow.Bind("/buy", flows.Buy)
	countFlow := buyFlow.Bind("\\d+", flows.Count)
	deliveryType := countFlow.Bind(".*", flows.DeliveryType)
	locationFlow := deliveryType.Bind(".*", flows.Location)
	locationFlow.Nexts = initFlow.Nexts

	done := make(chan bool)
	err := botflow.StartBot(os.Getenv("TELEGRAMM_TOKEN"), os.Getenv("TELEGRAM_BOT_NAME"), initFlow, done)
	if err != nil {
		panic(err)
	}

	web.StartServer(os.Getenv("TELEGRAM_BOT_WEB_PORT"))

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-signalChannel
	switch sig {
	case os.Interrupt:
		log.Printf("Got Interupt")
	case syscall.SIGTERM:
		log.Printf("Got SIGTERM")
	}
	log.Printf("Stop bot")
	done <- true
	log.Printf("Bye! Bye!")
}
