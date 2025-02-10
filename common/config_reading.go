package API

type Config struct {
	System struct {
		Type string
		Port int
	}
	Remote struct {
		Host string
		Port int
	}
}
