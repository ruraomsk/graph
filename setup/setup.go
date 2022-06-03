package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath  string   `toml:"logpath"`
	DataBase DataBase `toml:"dataBase"`
}

//DataBase настройки базы данных postresql
type DataBase struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBname   string `toml:"dbname"`
}
