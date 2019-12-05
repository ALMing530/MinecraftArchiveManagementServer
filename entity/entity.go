package entity

type Params struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Time   string  `json:"time"`
	Length float64 `json:"length"`
}
type Conf struct {
	Port           int    `port`
	ShellBash      string `shellbash`
	MinecraftDir   string `minecraftdir`
	CurrentArchive string `currentarchive`
	Authentication bool   "authentication"
	User           []User `user`
}
type User struct {
	Username string `username`
	Password string `password`
}
