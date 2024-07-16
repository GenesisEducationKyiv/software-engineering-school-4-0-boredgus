package service

type dispatchService struct {
	serverURL string
}

func NewDispatchService(serverURL string) *dispatchService {
	return &dispatchService{serverURL: serverURL}
}

func (s dispatchService) SubscribeForDispatchURL() string {
	return s.serverURL + "/services.DispatchService/SubscribeForDispatch"
}

func (s dispatchService) SubscribeForDispatchRevertURL() string {
	return s.serverURL + "/services.DispatchService/SubscribeForDispatchRevert"
}

func (s dispatchService) UnsubscribeFromDispatch() string {
	return s.serverURL + "/services.DispatchService/UnsubscribeFromDispatch"
}

func (s dispatchService) UnsubscribeFromDispatchRevert() string {
	return s.serverURL + "/services.DispatchService/UnsubscribeFromDispatchRevert"
}
