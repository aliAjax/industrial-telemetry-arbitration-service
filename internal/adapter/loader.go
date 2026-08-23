package adapter

type Config struct {
	Vendor string
	Labels map[string]string
	Strict bool
}

func Load(values map[string]string) Config {
	var labels map[string]string
	if values != nil {
		labels = make(map[string]string)
	}
	config := Config{Labels: labels, Strict: true}
	for key, value := range values {
		switch key {
		case "vendor":
			config.Vendor = value
		case "strict":
			config.Strict = value != "false"
		default:
			config.Labels[key] = value
		}
	}
	return config
}

func (c Config) Clone() Config {
	out := c
	out.Labels = make(map[string]string, len(c.Labels))
	for key, value := range c.Labels {
		out.Labels[key] = value
	}
	return out
}
