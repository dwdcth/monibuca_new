package config

type DB struct {
	DBType string `default:"sqlite" desc:"数据库类型"`
	DSN    string `default:"cascade.db" desc:"数据库文件路径"`
}