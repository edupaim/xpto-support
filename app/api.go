package app

type Api struct {
}

func initializeApi() (*Api, error) {
	return nil, nil
}

func (api *Api) run() <-chan error {
	apiErrorChan := make(chan error)
	go func() {
	}()
	return apiErrorChan
}

func (api *Api) shutdown() {
}
