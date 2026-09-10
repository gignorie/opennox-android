package server

import (
	"image/color"
	"reflect"

	"github.com/fxamacker/cbor/v2"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/noxnet/netxfer"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netstr"
	"github.com/opennox/opennox/v1/server/netlib"
)

func init() {
	RegisterNetXferExt(&XferObjectSetLabel{})
}

const (
	NetXferMOTD           = netxfer.Action(1)
	NetXferSavedata       = netxfer.Action(2)
	NetXferSaveServer     = netxfer.Action(3)
	NetXferMOTDType       = "MOTD"
	NetXferSavedataType   = "SAVEDATA"
	NetXferSaveServerType = "SAVE_SERVER"
)

// OpenNox extensions
const (
	NetXferExtOpenNox  = netxfer.Action(10)
	NetXferExtCBORType = "CBOR"
)

var (
	netXferCBORTypes = make(map[string]reflect.Type)
)

type XferConn struct {
	Conn netlib.SendStream
}

func (c XferConn) SendReliable(m netmsg.Message) error {
	_, err := c.Conn.SendReliableMsg(m)
	return err
}

func (c XferConn) SendUnreliable(m netmsg.Message) error {
	_, err := c.Conn.SendUnreliableMsg(m, true)
	return err
}

func RegisterNetXferExt(obj NetXferExt) {
	typ := obj.CBORType()
	if _, ok := netXferCBORTypes[typ]; ok {
		panic("already registered")
	}
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Pointer {
		panic("must be a pointer")
	}
	rt = rt.Elem()
	netXferCBORTypes[typ] = rt
}

func (s *Server) NetXferSendHost(data netxfer.Data, onDone netxfer.DoneFunc, onAbort netxfer.AbortFunc) bool {
	return s.netXferSend(HostPlayerIndex, data, false, onDone, onAbort)
}
func (s *Server) NetXferSend(pli ntype.PlayerInd, data netxfer.Data, onDone netxfer.DoneFunc, onAbort netxfer.AbortFunc) bool {
	return s.netXferSend(pli, data, true, onDone, onAbort)
}
func (s *Server) netXferSend(pli ntype.PlayerInd, data netxfer.Data, remote bool, onDone netxfer.DoneFunc, onAbort netxfer.AbortFunc) bool {
	if len(data.Data) == 0 {
		return false
	}
	var conn *netstr.Conn
	if !noxflags.HasGame(noxflags.GameHost) {
		conn = s.ClientConn().Conn
	} else {
		if !noxflags.HasGame(noxflags.GameClient) {
			return false
		}
		if !remote || pli == HostPlayerIndex {
			s.NetXferLocal(data)
			return true
		}
		conn = s.NetStr.ConnByPlayerInd(pli)
	}
	_, ok := s.NetXfer.StartSend(XferConn{conn}, data, onDone, onAbort)
	return ok
}

type NetXferExt interface {
	CBORType() string
}

type xferCBOR struct {
	_    struct{} `cbor:",toarray"`
	Type string
	Data cbor.RawMessage
}

func (s *Server) NetXferSendExt(pli ntype.PlayerInd, arr ...NetXferExt) bool {
	if len(arr) == 0 {
		return true
	}
	batch := make([]xferCBOR, 0, len(arr))
	for _, m := range arr {
		inner, err := cbor.Marshal(m)
		if err != nil {
			panic(err)
		}
		batch = append(batch, xferCBOR{Type: m.CBORType(), Data: inner})
	}
	data, err := cbor.Marshal(batch)
	if err != nil {
		panic(err)
	}
	return s.NetXferSend(pli, netxfer.Data{
		Action: NetXferExtOpenNox,
		Type:   NetXferExtCBORType,
		Data:   data,
	}, nil, nil)
}

func (s *Server) NetXferHandle(ind ntype.PlayerInd, data netxfer.Data) bool {
	switch data.Action {
	case NetXferExtOpenNox:
		switch data.Type {
		case NetXferExtCBORType:
			s.netXferHandleCBOR(ind, data.Data)
			return true // do not process further
		}
	}
	return false
}

func (s *Server) netXferHandleCBOR(ind ntype.PlayerInd, data []byte) {
	var arr []xferCBOR
	if err := cbor.Unmarshal(data, &arr); err != nil {
		s.Log.Error("cannot decode CBOR xfer", "err", err)
		return
	}
	for _, m := range arr {
		rt, ok := netXferCBORTypes[m.Type]
		if !ok {
			s.Log.Error("unsupported CBOR xfer", "type", m.Type)
			continue
		}
		obj := reflect.New(rt).Interface().(NetXferExt)
		if err := cbor.Unmarshal(m.Data, obj); err != nil {
			s.Log.Error("cannot decode CBOR xfer", "type", m.Type, "err", err)
			continue
		}
		s.handleCBORXfer(ind, obj)
	}
}

func (s *Server) handleCBORXfer(ind ntype.PlayerInd, m NetXferExt) {
	switch m := m.(type) {
	case *XferObjectSetLabel:
		if u := s.ObjectByNetCode(m.Object); u != nil {
			u.setDisplayName(m.Label, m.Color.Color())
		}
	}
	for _, fnc := range s.onXferExt {
		fnc(ind, m)
	}
}

func (s *Server) OnXferExt(fnc func(ind ntype.PlayerInd, m NetXferExt)) {
	s.onXferExt = append(s.onXferExt, fnc)
}

func (s *Server) appendIncomingXfer(dst []NetXferExt, obj *ObjectExt) []NetXferExt {
	netcode := s.GetUnitNetCode(obj)
	if obj.Label != "" || obj.LabelColor != nil {
		dst = append(dst, &XferObjectSetLabel{
			Object: netcode,
			Label:  obj.Label,
			Color:  AsXferColor(obj.LabelColor),
		})
	}
	return dst
}

func AsXferColor(cl color.Color) *XferColor {
	if cl == nil {
		return nil
	}
	rgba := color.NRGBAModel.Convert(cl).(color.NRGBA)
	return &XferColor{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A}
}

type XferColor struct {
	_          struct{} `cbor:",toarray"`
	R, G, B, A byte
}

func (cl *XferColor) Color() color.Color {
	if cl == nil {
		return nil
	}
	return color.NRGBA{R: cl.R, G: cl.G, B: cl.B, A: cl.A}
}

type XferObjectSetLabel struct {
	Object int        `cbor:"1,keyasint"`
	Label  string     `cbor:"2,keyasint,omitempty"`
	Color  *XferColor `cbor:"3,keyasint,omitempty"`
}

func (m *XferObjectSetLabel) CBORType() string {
	return "Obj.SetLabel"
}

func (s *Server) netSendObjectLabel(obj *Object, label string, cl color.Color) {
	m := &XferObjectSetLabel{
		Object: s.GetUnitNetCode(obj),
		Label:  label,
		Color:  AsXferColor(cl),
	}
	for _, pl := range s.Players.List() {
		s.NetXferSendExt(pl.PlayerIndex(), m)
	}
}
