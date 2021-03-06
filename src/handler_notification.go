package main

import (
	"time"
	"gitlab.com/irezaa/calendar/src/model"
	"fmt"
)

const (
	TOPIC = "/topics/someTopic"
)

func startNotificationHandler(app *App, startNow bool) {

	ticker := time.NewTicker(24 * time.Hour)

	if startNow {
		go sendPersonalEventsNotificationTask(app)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				sendPersonalEventsNotificationTask(app)
			}
		}
	}()

}

func sendPersonalEventsNotificationTask(app *App) {

	now := time.Now()

	personalEvents, err := app.EventRepository.FindTodayPersonalEvents(int(now.Month()), now.Year(), now.Day())

	if err == nil {
		if err == nil {
			fmt.Printf("Event Founded for today : %d \n", len(personalEvents))
			for index, event := range personalEvents {
				fmt.Printf("Event Index : %d \n", index)
				go processEventNotification(app, &event)
				time.Sleep(1 * time.Second)
			}
		} else {
			print(err.Error())
		}
	} else {
		print(err.Error())
	}
}

func processEventNotification(app *App, event *model.Event) {
	userSessions, err := app.SessionRepository.FindByUser(event.UserID)

	if err == nil {
		fmt.Printf("Sessions Founded for userID %s : %d \n", event.UserID, len(userSessions))
		for index, session := range userSessions {
			fmt.Printf("Session Index : %d \n", index)
			// todo send notification
			app.PushService.PushToDevice([]string{session.FcmToken}, map[string]string{
				"title": "Event :D",
				"msg":   event.EventTitle,
			})

			time.Sleep(1 * time.Second)
		}
	} else {
		print(err.Error())
	}
}
