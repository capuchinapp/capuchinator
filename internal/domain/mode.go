package domain

type Mode string

const (
	ModeUpdate      Mode = "update"
	ModeInstall     Mode = "install"
	ModeClearPGStat Mode = "clearpgstat"
)

func (s Mode) String() string {
	return string(s)
}
