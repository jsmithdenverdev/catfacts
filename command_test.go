package catfacts

import (
	"github.com/chuckpreslar/emission"
	"testing"
)

func TestHandleSubscribe(t *testing.T) {
	e := emission.NewEmitter()

	HandleSubscribe(nil, e)

	e.On(EventSubscriptionCreated, func() {
		return
	}).On(EventError, func(e error) {
		t.Fatal(e)
	}).Emit(CommandSubscribe, Subscription{Contact: "1234567890"})

}
