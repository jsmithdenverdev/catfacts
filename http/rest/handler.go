package rest

import (
	"catfacts/subscription"
	"net/http"
)

type subscriberService interface {
	Create(contact string) error
	Delete(contact string) error
}

func NewHandler(service subscriberService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/subscription", twilioCallback(service))

	return mux
}

func twilioCallback(service subscriberService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only allow posts to the endpoint
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body := r.PostForm.Get("Body")
		phone := r.PostForm.Get("From")

		action := subscription.GetAction(body)

		switch action {
		case subscription.ActionSubscribe:
			err := service.Create(phone)
			if err != nil {
				if err == subscription.ErrSubscriptionExists {
					err = WriteTwimlResponse(w, subscription.ReplyAlreadyExists)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = WriteTwimlResponse(w, subscription.ReplyWelcome)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case subscription.ActionUnsubscribe:
			err := service.Delete(phone)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = WriteTwimlResponse(w, subscription.ReplyGoodbye)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case subscription.ActionHelp:
			err = WriteTwimlResponse(w, subscription.ReplyHelp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case subscription.ActionUnknown:
			err = WriteTwimlResponse(w, subscription.ReplyUnknown)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
