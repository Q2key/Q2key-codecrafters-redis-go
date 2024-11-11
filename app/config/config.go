package config

type Config struct {
	dir        string
	dbfilename string
}

func NewConfig(dir, dbfilename string) *Config {
	return &Config{
		dir:        dir,
		dbfilename: dbfilename,
	}
}

func (r *Config) SetDir(val string) {
	r.dir = val
}

func (r *Config) SetDbFileName(val string) {
	r.dbfilename = val
}

func (r *Config) GetDir() string { return r.dir }

func (r *Config) GetDbFileName() string { return r.dbfilename }
