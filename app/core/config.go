package core

type Config struct {
	dir        string
	dbfilename string
}

func (r *Config) SetDir(val string) {
	r.dir = val
}

func (r *Config) SetDbFileName(val string) {
	r.dbfilename = val
}

func (r *Config) GetDir() string {
	return r.dir
}

func (r *Config) GetDbFileName() string {
	return r.dbfilename
}
