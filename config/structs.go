package config

type Configuration struct {
	Profiles map[string]string `yaml:"profiles"`
	Groups   []ServerGroup     `yaml:"groups"`
}

type ServerGroup struct {
	Name       string              `yaml:"name"`
	Defaults   map[string]string   `yaml:"defaults"`
	Servers    []map[string]string `yaml:"servers"`
	JsonSource *string             `yaml:"servers_source"`
}

func (c *Configuration) GetServerByName(name string) *map[string]string {
	serverConfiguration := make(map[string]string)
	var serverFound bool
	for _, group := range c.Groups {
		serverConfiguration = group.Defaults
		for _, server := range group.Servers {
			if _, ok := server["name"]; !ok {
				continue
			}
			if server["name"] == name {
				for k, v := range server {
					serverConfiguration[k] = v
				}
				serverFound = true
				break
			}
		}
	}
	if !serverFound {
		return nil
	}
	return &serverConfiguration
}
