package testdata

const USD_UAH_DISPATCH_ID = "f669a90d-d4aa-4285-bbce-6b14c6ff9065"

type SetupData struct {
	Filename   string
	Email      string
	DispatchID string
}

var CancelledSubscriptionData = SetupData{
	Filename:   "add_cancelled_subscription",
	Email:      "cancelled@gmail.com",
	DispatchID: USD_UAH_DISPATCH_ID,
}

var NewSubscriptionData = SetupData{
	Filename:   "add_new_subscription",
	Email:      "created@gmail.com",
	DispatchID: USD_UAH_DISPATCH_ID,
}
