package shared

type Player struct {
	SessionID  string
	XPosititon int
	YPosition  int
}
type PacketType rune

const (
	PacketTypeConnect    PacketType = 'C'
	PacketTypeMap        PacketType = 'M'
	PacketTypeMove       PacketType = 'L'
	PacketTypeDisconnect PacketType = 'X'
	PacketTypeErr        PacketType = 'E'
)

func NewPlayer() *Player {
	return &Player{XPosititon: 0, YPosition: 0}
}
