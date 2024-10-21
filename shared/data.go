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
type Cell int32

type Universe struct {
	Map    []Cell
	Width  int
	Height int
}
func NewPlayer() *Player {
	return &Player{XPosititon: 0, YPosition: 0}
}
