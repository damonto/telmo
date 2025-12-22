package msisdn

type Runner interface {
	Run(data []byte) error
	Select() ([]byte, error)
}
