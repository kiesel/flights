package mode_s

type Message struct {
  bits    int
  type    int
  crcOk   int
  crc     uint32

}

func DetectModeA()