package service

type customerService struct {
	serverURL string
}

func NewCustomerService(serverURL string) *customerService {
	return &customerService{serverURL: serverURL}
}

func (s customerService) CreateCustomerURL() string {
	return s.serverURL + "/services.CustomerService/CreateCustomer"
}

func (s customerService) CreateCustomerRevertURL() string {
	return s.serverURL + "/services.CustomerService/CreateCustomerRevert"
}
