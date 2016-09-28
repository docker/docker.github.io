// ************************************************************
// DO NOT EDIT.
// THIS FILE IS AUTO-GENERATED BY codecgen.
// ************************************************************

package etcd

import (
	"errors"
	"fmt"
	codec1978 "github.com/ugorji/go/codec"
	pkg1_http "net/http"
	"reflect"
	"runtime"
	time "time"
)

const (
	// ----- content types ----
	codecSelferC_UTF81978 = 1
	codecSelferC_RAW1978  = 0
	// ----- value types used ----
	codecSelferValueTypeArray1978 = 10
	codecSelferValueTypeMap1978   = 9
	// ----- containerStateValues ----
	codecSelfer_containerMapKey1978    = 2
	codecSelfer_containerMapValue1978  = 3
	codecSelfer_containerMapEnd1978    = 4
	codecSelfer_containerArrayElem1978 = 6
	codecSelfer_containerArrayEnd1978  = 7
)

var (
	codecSelferBitsize1978                         = uint8(reflect.TypeOf(uint(0)).Bits())
	codecSelferOnlyMapOrArrayEncodeToStructErr1978 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer1978 struct{}

func init() {
	if codec1978.GenVersion != 5 {
		_, file, _, _ := runtime.Caller(0)
		err := fmt.Errorf("codecgen version mismatch: current: %v, need %v. Re-generate file: %v",
			5, codec1978.GenVersion, file)
		panic(err)
	}
	if false { // reference the types, but skip this branch at build/run time
		var v0 pkg1_http.Header
		var v1 time.Time
		_, _ = v0, v1
	}
}

func (x responseType) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	yym1 := z.EncBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.EncExt(x) {
	} else {
		r.EncodeInt(int64(x))
	}
}

func (x *responseType) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym2 := z.DecBinary()
	_ = yym2
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		*((*int)(x)) = int(r.DecodeInt(codecSelferBitsize1978))
	}
}

func (x *RawResponse) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym3 := z.EncBinary()
		_ = yym3
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			yysep4 := !z.EncBinary()
			yy2arr4 := z.EncBasicHandle().StructToArray
			var yyq4 [3]bool
			_, _, _ = yysep4, yyq4, yy2arr4
			const yyr4 bool = false
			var yynn4 int
			if yyr4 || yy2arr4 {
				r.EncodeArrayStart(3)
			} else {
				yynn4 = 3
				for _, b := range yyq4 {
					if b {
						yynn4++
					}
				}
				r.EncodeMapStart(yynn4)
				yynn4 = 0
			}
			if yyr4 || yy2arr4 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				yym6 := z.EncBinary()
				_ = yym6
				if false {
				} else {
					r.EncodeInt(int64(x.StatusCode))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("StatusCode"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				yym7 := z.EncBinary()
				_ = yym7
				if false {
				} else {
					r.EncodeInt(int64(x.StatusCode))
				}
			}
			if yyr4 || yy2arr4 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if x.Body == nil {
					r.EncodeNil()
				} else {
					yym9 := z.EncBinary()
					_ = yym9
					if false {
					} else {
						r.EncodeStringBytes(codecSelferC_RAW1978, []byte(x.Body))
					}
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("Body"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				if x.Body == nil {
					r.EncodeNil()
				} else {
					yym10 := z.EncBinary()
					_ = yym10
					if false {
					} else {
						r.EncodeStringBytes(codecSelferC_RAW1978, []byte(x.Body))
					}
				}
			}
			if yyr4 || yy2arr4 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if x.Header == nil {
					r.EncodeNil()
				} else {
					yym12 := z.EncBinary()
					_ = yym12
					if false {
					} else if z.HasExtensions() && z.EncExt(x.Header) {
					} else {
						h.enchttp_Header((pkg1_http.Header)(x.Header), e)
					}
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("Header"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				if x.Header == nil {
					r.EncodeNil()
				} else {
					yym13 := z.EncBinary()
					_ = yym13
					if false {
					} else if z.HasExtensions() && z.EncExt(x.Header) {
					} else {
						h.enchttp_Header((pkg1_http.Header)(x.Header), e)
					}
				}
			}
			if yyr4 || yy2arr4 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd1978)
			} else {
				z.EncSendContainerState(codecSelfer_containerMapEnd1978)
			}
		}
	}
}

func (x *RawResponse) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym14 := z.DecBinary()
	_ = yym14
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		yyct15 := r.ContainerType()
		if yyct15 == codecSelferValueTypeMap1978 {
			yyl15 := r.ReadMapStart()
			if yyl15 == 0 {
				z.DecSendContainerState(codecSelfer_containerMapEnd1978)
			} else {
				x.codecDecodeSelfFromMap(yyl15, d)
			}
		} else if yyct15 == codecSelferValueTypeArray1978 {
			yyl15 := r.ReadArrayStart()
			if yyl15 == 0 {
				z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
			} else {
				x.codecDecodeSelfFromArray(yyl15, d)
			}
		} else {
			panic(codecSelferOnlyMapOrArrayEncodeToStructErr1978)
		}
	}
}

func (x *RawResponse) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yys16Slc = z.DecScratchBuffer() // default slice to decode into
	_ = yys16Slc
	var yyhl16 bool = l >= 0
	for yyj16 := 0; ; yyj16++ {
		if yyhl16 {
			if yyj16 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		z.DecSendContainerState(codecSelfer_containerMapKey1978)
		yys16Slc = r.DecodeBytes(yys16Slc, true, true)
		yys16 := string(yys16Slc)
		z.DecSendContainerState(codecSelfer_containerMapValue1978)
		switch yys16 {
		case "StatusCode":
			if r.TryDecodeAsNil() {
				x.StatusCode = 0
			} else {
				x.StatusCode = int(r.DecodeInt(codecSelferBitsize1978))
			}
		case "Body":
			if r.TryDecodeAsNil() {
				x.Body = nil
			} else {
				yyv18 := &x.Body
				yym19 := z.DecBinary()
				_ = yym19
				if false {
				} else {
					*yyv18 = r.DecodeBytes(*(*[]byte)(yyv18), false, false)
				}
			}
		case "Header":
			if r.TryDecodeAsNil() {
				x.Header = nil
			} else {
				yyv20 := &x.Header
				yym21 := z.DecBinary()
				_ = yym21
				if false {
				} else if z.HasExtensions() && z.DecExt(yyv20) {
				} else {
					h.dechttp_Header((*pkg1_http.Header)(yyv20), d)
				}
			}
		default:
			z.DecStructFieldNotFound(-1, yys16)
		} // end switch yys16
	} // end for yyj16
	z.DecSendContainerState(codecSelfer_containerMapEnd1978)
}

func (x *RawResponse) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj22 int
	var yyb22 bool
	var yyhl22 bool = l >= 0
	yyj22++
	if yyhl22 {
		yyb22 = yyj22 > l
	} else {
		yyb22 = r.CheckBreak()
	}
	if yyb22 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.StatusCode = 0
	} else {
		x.StatusCode = int(r.DecodeInt(codecSelferBitsize1978))
	}
	yyj22++
	if yyhl22 {
		yyb22 = yyj22 > l
	} else {
		yyb22 = r.CheckBreak()
	}
	if yyb22 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.Body = nil
	} else {
		yyv24 := &x.Body
		yym25 := z.DecBinary()
		_ = yym25
		if false {
		} else {
			*yyv24 = r.DecodeBytes(*(*[]byte)(yyv24), false, false)
		}
	}
	yyj22++
	if yyhl22 {
		yyb22 = yyj22 > l
	} else {
		yyb22 = r.CheckBreak()
	}
	if yyb22 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.Header = nil
	} else {
		yyv26 := &x.Header
		yym27 := z.DecBinary()
		_ = yym27
		if false {
		} else if z.HasExtensions() && z.DecExt(yyv26) {
		} else {
			h.dechttp_Header((*pkg1_http.Header)(yyv26), d)
		}
	}
	for {
		yyj22++
		if yyhl22 {
			yyb22 = yyj22 > l
		} else {
			yyb22 = r.CheckBreak()
		}
		if yyb22 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem1978)
		z.DecStructFieldNotFound(yyj22-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
}

func (x *Response) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym28 := z.EncBinary()
		_ = yym28
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			yysep29 := !z.EncBinary()
			yy2arr29 := z.EncBasicHandle().StructToArray
			var yyq29 [6]bool
			_, _, _ = yysep29, yyq29, yy2arr29
			const yyr29 bool = false
			yyq29[2] = x.PrevNode != nil
			var yynn29 int
			if yyr29 || yy2arr29 {
				r.EncodeArrayStart(6)
			} else {
				yynn29 = 5
				for _, b := range yyq29 {
					if b {
						yynn29++
					}
				}
				r.EncodeMapStart(yynn29)
				yynn29 = 0
			}
			if yyr29 || yy2arr29 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				yym31 := z.EncBinary()
				_ = yym31
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81978, string(x.Action))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("action"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				yym32 := z.EncBinary()
				_ = yym32
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81978, string(x.Action))
				}
			}
			if yyr29 || yy2arr29 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if x.Node == nil {
					r.EncodeNil()
				} else {
					x.Node.CodecEncodeSelf(e)
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("node"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				if x.Node == nil {
					r.EncodeNil()
				} else {
					x.Node.CodecEncodeSelf(e)
				}
			}
			if yyr29 || yy2arr29 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq29[2] {
					if x.PrevNode == nil {
						r.EncodeNil()
					} else {
						x.PrevNode.CodecEncodeSelf(e)
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq29[2] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("prevNode"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					if x.PrevNode == nil {
						r.EncodeNil()
					} else {
						x.PrevNode.CodecEncodeSelf(e)
					}
				}
			}
			if yyr29 || yy2arr29 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				yym36 := z.EncBinary()
				_ = yym36
				if false {
				} else {
					r.EncodeUint(uint64(x.EtcdIndex))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("etcdIndex"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				yym37 := z.EncBinary()
				_ = yym37
				if false {
				} else {
					r.EncodeUint(uint64(x.EtcdIndex))
				}
			}
			if yyr29 || yy2arr29 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				yym39 := z.EncBinary()
				_ = yym39
				if false {
				} else {
					r.EncodeUint(uint64(x.RaftIndex))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("raftIndex"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				yym40 := z.EncBinary()
				_ = yym40
				if false {
				} else {
					r.EncodeUint(uint64(x.RaftIndex))
				}
			}
			if yyr29 || yy2arr29 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				yym42 := z.EncBinary()
				_ = yym42
				if false {
				} else {
					r.EncodeUint(uint64(x.RaftTerm))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("raftTerm"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				yym43 := z.EncBinary()
				_ = yym43
				if false {
				} else {
					r.EncodeUint(uint64(x.RaftTerm))
				}
			}
			if yyr29 || yy2arr29 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd1978)
			} else {
				z.EncSendContainerState(codecSelfer_containerMapEnd1978)
			}
		}
	}
}

func (x *Response) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym44 := z.DecBinary()
	_ = yym44
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		yyct45 := r.ContainerType()
		if yyct45 == codecSelferValueTypeMap1978 {
			yyl45 := r.ReadMapStart()
			if yyl45 == 0 {
				z.DecSendContainerState(codecSelfer_containerMapEnd1978)
			} else {
				x.codecDecodeSelfFromMap(yyl45, d)
			}
		} else if yyct45 == codecSelferValueTypeArray1978 {
			yyl45 := r.ReadArrayStart()
			if yyl45 == 0 {
				z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
			} else {
				x.codecDecodeSelfFromArray(yyl45, d)
			}
		} else {
			panic(codecSelferOnlyMapOrArrayEncodeToStructErr1978)
		}
	}
}

func (x *Response) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yys46Slc = z.DecScratchBuffer() // default slice to decode into
	_ = yys46Slc
	var yyhl46 bool = l >= 0
	for yyj46 := 0; ; yyj46++ {
		if yyhl46 {
			if yyj46 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		z.DecSendContainerState(codecSelfer_containerMapKey1978)
		yys46Slc = r.DecodeBytes(yys46Slc, true, true)
		yys46 := string(yys46Slc)
		z.DecSendContainerState(codecSelfer_containerMapValue1978)
		switch yys46 {
		case "action":
			if r.TryDecodeAsNil() {
				x.Action = ""
			} else {
				x.Action = string(r.DecodeString())
			}
		case "node":
			if r.TryDecodeAsNil() {
				if x.Node != nil {
					x.Node = nil
				}
			} else {
				if x.Node == nil {
					x.Node = new(Node)
				}
				x.Node.CodecDecodeSelf(d)
			}
		case "prevNode":
			if r.TryDecodeAsNil() {
				if x.PrevNode != nil {
					x.PrevNode = nil
				}
			} else {
				if x.PrevNode == nil {
					x.PrevNode = new(Node)
				}
				x.PrevNode.CodecDecodeSelf(d)
			}
		case "etcdIndex":
			if r.TryDecodeAsNil() {
				x.EtcdIndex = 0
			} else {
				x.EtcdIndex = uint64(r.DecodeUint(64))
			}
		case "raftIndex":
			if r.TryDecodeAsNil() {
				x.RaftIndex = 0
			} else {
				x.RaftIndex = uint64(r.DecodeUint(64))
			}
		case "raftTerm":
			if r.TryDecodeAsNil() {
				x.RaftTerm = 0
			} else {
				x.RaftTerm = uint64(r.DecodeUint(64))
			}
		default:
			z.DecStructFieldNotFound(-1, yys46)
		} // end switch yys46
	} // end for yyj46
	z.DecSendContainerState(codecSelfer_containerMapEnd1978)
}

func (x *Response) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj53 int
	var yyb53 bool
	var yyhl53 bool = l >= 0
	yyj53++
	if yyhl53 {
		yyb53 = yyj53 > l
	} else {
		yyb53 = r.CheckBreak()
	}
	if yyb53 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.Action = ""
	} else {
		x.Action = string(r.DecodeString())
	}
	yyj53++
	if yyhl53 {
		yyb53 = yyj53 > l
	} else {
		yyb53 = r.CheckBreak()
	}
	if yyb53 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		if x.Node != nil {
			x.Node = nil
		}
	} else {
		if x.Node == nil {
			x.Node = new(Node)
		}
		x.Node.CodecDecodeSelf(d)
	}
	yyj53++
	if yyhl53 {
		yyb53 = yyj53 > l
	} else {
		yyb53 = r.CheckBreak()
	}
	if yyb53 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		if x.PrevNode != nil {
			x.PrevNode = nil
		}
	} else {
		if x.PrevNode == nil {
			x.PrevNode = new(Node)
		}
		x.PrevNode.CodecDecodeSelf(d)
	}
	yyj53++
	if yyhl53 {
		yyb53 = yyj53 > l
	} else {
		yyb53 = r.CheckBreak()
	}
	if yyb53 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.EtcdIndex = 0
	} else {
		x.EtcdIndex = uint64(r.DecodeUint(64))
	}
	yyj53++
	if yyhl53 {
		yyb53 = yyj53 > l
	} else {
		yyb53 = r.CheckBreak()
	}
	if yyb53 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.RaftIndex = 0
	} else {
		x.RaftIndex = uint64(r.DecodeUint(64))
	}
	yyj53++
	if yyhl53 {
		yyb53 = yyj53 > l
	} else {
		yyb53 = r.CheckBreak()
	}
	if yyb53 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.RaftTerm = 0
	} else {
		x.RaftTerm = uint64(r.DecodeUint(64))
	}
	for {
		yyj53++
		if yyhl53 {
			yyb53 = yyj53 > l
		} else {
			yyb53 = r.CheckBreak()
		}
		if yyb53 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem1978)
		z.DecStructFieldNotFound(yyj53-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
}

func (x *Node) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym60 := z.EncBinary()
		_ = yym60
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			yysep61 := !z.EncBinary()
			yy2arr61 := z.EncBasicHandle().StructToArray
			var yyq61 [8]bool
			_, _, _ = yysep61, yyq61, yy2arr61
			const yyr61 bool = false
			yyq61[1] = x.Value != ""
			yyq61[2] = x.Dir != false
			yyq61[3] = x.Expiration != nil
			yyq61[4] = x.TTL != 0
			yyq61[5] = len(x.Nodes) != 0
			yyq61[6] = x.ModifiedIndex != 0
			yyq61[7] = x.CreatedIndex != 0
			var yynn61 int
			if yyr61 || yy2arr61 {
				r.EncodeArrayStart(8)
			} else {
				yynn61 = 1
				for _, b := range yyq61 {
					if b {
						yynn61++
					}
				}
				r.EncodeMapStart(yynn61)
				yynn61 = 0
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				yym63 := z.EncBinary()
				_ = yym63
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81978, string(x.Key))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1978)
				r.EncodeString(codecSelferC_UTF81978, string("key"))
				z.EncSendContainerState(codecSelfer_containerMapValue1978)
				yym64 := z.EncBinary()
				_ = yym64
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81978, string(x.Key))
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq61[1] {
					yym66 := z.EncBinary()
					_ = yym66
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81978, string(x.Value))
					}
				} else {
					r.EncodeString(codecSelferC_UTF81978, "")
				}
			} else {
				if yyq61[1] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("value"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					yym67 := z.EncBinary()
					_ = yym67
					if false {
					} else {
						r.EncodeString(codecSelferC_UTF81978, string(x.Value))
					}
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq61[2] {
					yym69 := z.EncBinary()
					_ = yym69
					if false {
					} else {
						r.EncodeBool(bool(x.Dir))
					}
				} else {
					r.EncodeBool(false)
				}
			} else {
				if yyq61[2] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("dir"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					yym70 := z.EncBinary()
					_ = yym70
					if false {
					} else {
						r.EncodeBool(bool(x.Dir))
					}
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq61[3] {
					if x.Expiration == nil {
						r.EncodeNil()
					} else {
						yym72 := z.EncBinary()
						_ = yym72
						if false {
						} else if yym73 := z.TimeRtidIfBinc(); yym73 != 0 {
							r.EncodeBuiltin(yym73, x.Expiration)
						} else if z.HasExtensions() && z.EncExt(x.Expiration) {
						} else if yym72 {
							z.EncBinaryMarshal(x.Expiration)
						} else if !yym72 && z.IsJSONHandle() {
							z.EncJSONMarshal(x.Expiration)
						} else {
							z.EncFallback(x.Expiration)
						}
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq61[3] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("expiration"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					if x.Expiration == nil {
						r.EncodeNil()
					} else {
						yym74 := z.EncBinary()
						_ = yym74
						if false {
						} else if yym75 := z.TimeRtidIfBinc(); yym75 != 0 {
							r.EncodeBuiltin(yym75, x.Expiration)
						} else if z.HasExtensions() && z.EncExt(x.Expiration) {
						} else if yym74 {
							z.EncBinaryMarshal(x.Expiration)
						} else if !yym74 && z.IsJSONHandle() {
							z.EncJSONMarshal(x.Expiration)
						} else {
							z.EncFallback(x.Expiration)
						}
					}
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq61[4] {
					yym77 := z.EncBinary()
					_ = yym77
					if false {
					} else {
						r.EncodeInt(int64(x.TTL))
					}
				} else {
					r.EncodeInt(0)
				}
			} else {
				if yyq61[4] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("ttl"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					yym78 := z.EncBinary()
					_ = yym78
					if false {
					} else {
						r.EncodeInt(int64(x.TTL))
					}
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq61[5] {
					if x.Nodes == nil {
						r.EncodeNil()
					} else {
						x.Nodes.CodecEncodeSelf(e)
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq61[5] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("nodes"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					if x.Nodes == nil {
						r.EncodeNil()
					} else {
						x.Nodes.CodecEncodeSelf(e)
					}
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq61[6] {
					yym81 := z.EncBinary()
					_ = yym81
					if false {
					} else {
						r.EncodeUint(uint64(x.ModifiedIndex))
					}
				} else {
					r.EncodeUint(0)
				}
			} else {
				if yyq61[6] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("modifiedIndex"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					yym82 := z.EncBinary()
					_ = yym82
					if false {
					} else {
						r.EncodeUint(uint64(x.ModifiedIndex))
					}
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1978)
				if yyq61[7] {
					yym84 := z.EncBinary()
					_ = yym84
					if false {
					} else {
						r.EncodeUint(uint64(x.CreatedIndex))
					}
				} else {
					r.EncodeUint(0)
				}
			} else {
				if yyq61[7] {
					z.EncSendContainerState(codecSelfer_containerMapKey1978)
					r.EncodeString(codecSelferC_UTF81978, string("createdIndex"))
					z.EncSendContainerState(codecSelfer_containerMapValue1978)
					yym85 := z.EncBinary()
					_ = yym85
					if false {
					} else {
						r.EncodeUint(uint64(x.CreatedIndex))
					}
				}
			}
			if yyr61 || yy2arr61 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd1978)
			} else {
				z.EncSendContainerState(codecSelfer_containerMapEnd1978)
			}
		}
	}
}

func (x *Node) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym86 := z.DecBinary()
	_ = yym86
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		yyct87 := r.ContainerType()
		if yyct87 == codecSelferValueTypeMap1978 {
			yyl87 := r.ReadMapStart()
			if yyl87 == 0 {
				z.DecSendContainerState(codecSelfer_containerMapEnd1978)
			} else {
				x.codecDecodeSelfFromMap(yyl87, d)
			}
		} else if yyct87 == codecSelferValueTypeArray1978 {
			yyl87 := r.ReadArrayStart()
			if yyl87 == 0 {
				z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
			} else {
				x.codecDecodeSelfFromArray(yyl87, d)
			}
		} else {
			panic(codecSelferOnlyMapOrArrayEncodeToStructErr1978)
		}
	}
}

func (x *Node) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yys88Slc = z.DecScratchBuffer() // default slice to decode into
	_ = yys88Slc
	var yyhl88 bool = l >= 0
	for yyj88 := 0; ; yyj88++ {
		if yyhl88 {
			if yyj88 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		z.DecSendContainerState(codecSelfer_containerMapKey1978)
		yys88Slc = r.DecodeBytes(yys88Slc, true, true)
		yys88 := string(yys88Slc)
		z.DecSendContainerState(codecSelfer_containerMapValue1978)
		switch yys88 {
		case "key":
			if r.TryDecodeAsNil() {
				x.Key = ""
			} else {
				x.Key = string(r.DecodeString())
			}
		case "value":
			if r.TryDecodeAsNil() {
				x.Value = ""
			} else {
				x.Value = string(r.DecodeString())
			}
		case "dir":
			if r.TryDecodeAsNil() {
				x.Dir = false
			} else {
				x.Dir = bool(r.DecodeBool())
			}
		case "expiration":
			if r.TryDecodeAsNil() {
				if x.Expiration != nil {
					x.Expiration = nil
				}
			} else {
				if x.Expiration == nil {
					x.Expiration = new(time.Time)
				}
				yym93 := z.DecBinary()
				_ = yym93
				if false {
				} else if yym94 := z.TimeRtidIfBinc(); yym94 != 0 {
					r.DecodeBuiltin(yym94, x.Expiration)
				} else if z.HasExtensions() && z.DecExt(x.Expiration) {
				} else if yym93 {
					z.DecBinaryUnmarshal(x.Expiration)
				} else if !yym93 && z.IsJSONHandle() {
					z.DecJSONUnmarshal(x.Expiration)
				} else {
					z.DecFallback(x.Expiration, false)
				}
			}
		case "ttl":
			if r.TryDecodeAsNil() {
				x.TTL = 0
			} else {
				x.TTL = int64(r.DecodeInt(64))
			}
		case "nodes":
			if r.TryDecodeAsNil() {
				x.Nodes = nil
			} else {
				yyv96 := &x.Nodes
				yyv96.CodecDecodeSelf(d)
			}
		case "modifiedIndex":
			if r.TryDecodeAsNil() {
				x.ModifiedIndex = 0
			} else {
				x.ModifiedIndex = uint64(r.DecodeUint(64))
			}
		case "createdIndex":
			if r.TryDecodeAsNil() {
				x.CreatedIndex = 0
			} else {
				x.CreatedIndex = uint64(r.DecodeUint(64))
			}
		default:
			z.DecStructFieldNotFound(-1, yys88)
		} // end switch yys88
	} // end for yyj88
	z.DecSendContainerState(codecSelfer_containerMapEnd1978)
}

func (x *Node) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj99 int
	var yyb99 bool
	var yyhl99 bool = l >= 0
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.Key = ""
	} else {
		x.Key = string(r.DecodeString())
	}
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.Value = ""
	} else {
		x.Value = string(r.DecodeString())
	}
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.Dir = false
	} else {
		x.Dir = bool(r.DecodeBool())
	}
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		if x.Expiration != nil {
			x.Expiration = nil
		}
	} else {
		if x.Expiration == nil {
			x.Expiration = new(time.Time)
		}
		yym104 := z.DecBinary()
		_ = yym104
		if false {
		} else if yym105 := z.TimeRtidIfBinc(); yym105 != 0 {
			r.DecodeBuiltin(yym105, x.Expiration)
		} else if z.HasExtensions() && z.DecExt(x.Expiration) {
		} else if yym104 {
			z.DecBinaryUnmarshal(x.Expiration)
		} else if !yym104 && z.IsJSONHandle() {
			z.DecJSONUnmarshal(x.Expiration)
		} else {
			z.DecFallback(x.Expiration, false)
		}
	}
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.TTL = 0
	} else {
		x.TTL = int64(r.DecodeInt(64))
	}
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.Nodes = nil
	} else {
		yyv107 := &x.Nodes
		yyv107.CodecDecodeSelf(d)
	}
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.ModifiedIndex = 0
	} else {
		x.ModifiedIndex = uint64(r.DecodeUint(64))
	}
	yyj99++
	if yyhl99 {
		yyb99 = yyj99 > l
	} else {
		yyb99 = r.CheckBreak()
	}
	if yyb99 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1978)
	if r.TryDecodeAsNil() {
		x.CreatedIndex = 0
	} else {
		x.CreatedIndex = uint64(r.DecodeUint(64))
	}
	for {
		yyj99++
		if yyhl99 {
			yyb99 = yyj99 > l
		} else {
			yyb99 = r.CheckBreak()
		}
		if yyb99 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem1978)
		z.DecStructFieldNotFound(yyj99-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd1978)
}

func (x Nodes) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym110 := z.EncBinary()
		_ = yym110
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			h.encNodes((Nodes)(x), e)
		}
	}
}

func (x *Nodes) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym111 := z.DecBinary()
	_ = yym111
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		h.decNodes((*Nodes)(x), d)
	}
}

func (x codecSelfer1978) enchttp_Header(v pkg1_http.Header, e *codec1978.Encoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeMapStart(len(v))
	for yyk112, yyv112 := range v {
		z.EncSendContainerState(codecSelfer_containerMapKey1978)
		yym113 := z.EncBinary()
		_ = yym113
		if false {
		} else {
			r.EncodeString(codecSelferC_UTF81978, string(yyk112))
		}
		z.EncSendContainerState(codecSelfer_containerMapValue1978)
		if yyv112 == nil {
			r.EncodeNil()
		} else {
			yym114 := z.EncBinary()
			_ = yym114
			if false {
			} else {
				z.F.EncSliceStringV(yyv112, false, e)
			}
		}
	}
	z.EncSendContainerState(codecSelfer_containerMapEnd1978)
}

func (x codecSelfer1978) dechttp_Header(v *pkg1_http.Header, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv115 := *v
	yyl115 := r.ReadMapStart()
	yybh115 := z.DecBasicHandle()
	if yyv115 == nil {
		yyrl115, _ := z.DecInferLen(yyl115, yybh115.MaxInitLen, 40)
		yyv115 = make(map[string][]string, yyrl115)
		*v = yyv115
	}
	var yymk115 string
	var yymv115 []string
	var yymg115 bool
	if yybh115.MapValueReset {
		yymg115 = true
	}
	if yyl115 > 0 {
		for yyj115 := 0; yyj115 < yyl115; yyj115++ {
			z.DecSendContainerState(codecSelfer_containerMapKey1978)
			if r.TryDecodeAsNil() {
				yymk115 = ""
			} else {
				yymk115 = string(r.DecodeString())
			}

			if yymg115 {
				yymv115 = yyv115[yymk115]
			} else {
				yymv115 = nil
			}
			z.DecSendContainerState(codecSelfer_containerMapValue1978)
			if r.TryDecodeAsNil() {
				yymv115 = nil
			} else {
				yyv117 := &yymv115
				yym118 := z.DecBinary()
				_ = yym118
				if false {
				} else {
					z.F.DecSliceStringX(yyv117, false, d)
				}
			}

			if yyv115 != nil {
				yyv115[yymk115] = yymv115
			}
		}
	} else if yyl115 < 0 {
		for yyj115 := 0; !r.CheckBreak(); yyj115++ {
			z.DecSendContainerState(codecSelfer_containerMapKey1978)
			if r.TryDecodeAsNil() {
				yymk115 = ""
			} else {
				yymk115 = string(r.DecodeString())
			}

			if yymg115 {
				yymv115 = yyv115[yymk115]
			} else {
				yymv115 = nil
			}
			z.DecSendContainerState(codecSelfer_containerMapValue1978)
			if r.TryDecodeAsNil() {
				yymv115 = nil
			} else {
				yyv120 := &yymv115
				yym121 := z.DecBinary()
				_ = yym121
				if false {
				} else {
					z.F.DecSliceStringX(yyv120, false, d)
				}
			}

			if yyv115 != nil {
				yyv115[yymk115] = yymv115
			}
		}
	} // else len==0: TODO: Should we clear map entries?
	z.DecSendContainerState(codecSelfer_containerMapEnd1978)
}

func (x codecSelfer1978) encNodes(v Nodes, e *codec1978.Encoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeArrayStart(len(v))
	for _, yyv122 := range v {
		z.EncSendContainerState(codecSelfer_containerArrayElem1978)
		if yyv122 == nil {
			r.EncodeNil()
		} else {
			yyv122.CodecEncodeSelf(e)
		}
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd1978)
}

func (x codecSelfer1978) decNodes(v *Nodes, d *codec1978.Decoder) {
	var h codecSelfer1978
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv123 := *v
	yyh123, yyl123 := z.DecSliceHelperStart()
	var yyc123 bool
	if yyl123 == 0 {
		if yyv123 == nil {
			yyv123 = []*Node{}
			yyc123 = true
		} else if len(yyv123) != 0 {
			yyv123 = yyv123[:0]
			yyc123 = true
		}
	} else if yyl123 > 0 {
		var yyrr123, yyrl123 int
		var yyrt123 bool
		if yyl123 > cap(yyv123) {

			yyrg123 := len(yyv123) > 0
			yyv2123 := yyv123
			yyrl123, yyrt123 = z.DecInferLen(yyl123, z.DecBasicHandle().MaxInitLen, 8)
			if yyrt123 {
				if yyrl123 <= cap(yyv123) {
					yyv123 = yyv123[:yyrl123]
				} else {
					yyv123 = make([]*Node, yyrl123)
				}
			} else {
				yyv123 = make([]*Node, yyrl123)
			}
			yyc123 = true
			yyrr123 = len(yyv123)
			if yyrg123 {
				copy(yyv123, yyv2123)
			}
		} else if yyl123 != len(yyv123) {
			yyv123 = yyv123[:yyl123]
			yyc123 = true
		}
		yyj123 := 0
		for ; yyj123 < yyrr123; yyj123++ {
			yyh123.ElemContainerState(yyj123)
			if r.TryDecodeAsNil() {
				if yyv123[yyj123] != nil {
					*yyv123[yyj123] = Node{}
				}
			} else {
				if yyv123[yyj123] == nil {
					yyv123[yyj123] = new(Node)
				}
				yyw124 := yyv123[yyj123]
				yyw124.CodecDecodeSelf(d)
			}

		}
		if yyrt123 {
			for ; yyj123 < yyl123; yyj123++ {
				yyv123 = append(yyv123, nil)
				yyh123.ElemContainerState(yyj123)
				if r.TryDecodeAsNil() {
					if yyv123[yyj123] != nil {
						*yyv123[yyj123] = Node{}
					}
				} else {
					if yyv123[yyj123] == nil {
						yyv123[yyj123] = new(Node)
					}
					yyw125 := yyv123[yyj123]
					yyw125.CodecDecodeSelf(d)
				}

			}
		}

	} else {
		yyj123 := 0
		for ; !r.CheckBreak(); yyj123++ {

			if yyj123 >= len(yyv123) {
				yyv123 = append(yyv123, nil) // var yyz123 *Node
				yyc123 = true
			}
			yyh123.ElemContainerState(yyj123)
			if yyj123 < len(yyv123) {
				if r.TryDecodeAsNil() {
					if yyv123[yyj123] != nil {
						*yyv123[yyj123] = Node{}
					}
				} else {
					if yyv123[yyj123] == nil {
						yyv123[yyj123] = new(Node)
					}
					yyw126 := yyv123[yyj123]
					yyw126.CodecDecodeSelf(d)
				}

			} else {
				z.DecSwallow()
			}

		}
		if yyj123 < len(yyv123) {
			yyv123 = yyv123[:yyj123]
			yyc123 = true
		} else if yyj123 == 0 && yyv123 == nil {
			yyv123 = []*Node{}
			yyc123 = true
		}
	}
	yyh123.End()
	if yyc123 {
		*v = yyv123
	}
}
