package handshake

type HandshakeTag uint32

const (
	TAG_CHLO	HandshakeTag = 'C' + 'H' << 8 + 'L' << 16 + 'O' << 24
	TAG_REG		HandshakeTag = 'R' + 'E' << 8 + 'G' << 16
	TAG_SCFG	HandshakeTag = 'S' + 'C' << 8 + 'F' << 16 + 'G' << 24
	TAG_PAD		HandshakeTag = 'P' + 'A' << 8 + 'G' << 16
	TAG_SNI		HandshakeTag = 'S' + 'N' << 8 + 'I' << 16
	TAG_VER		HandshakeTag = 'V' + 'E' << 8 + 'R' << 16
	TAG_CCS		HandshakeTag = 'C' + 'C' << 8 + 'S' << 16
	TAG_CCRT	HandshakeTag = 'C' + 'C' << 8 + 'R' << 16 + 'T' << 24
	TAG_MSPC	HandshakeTag = 'M' + 'S' << 8 + 'P' << 16 + 'C' << 24
	TAG_MIDS	HandshakeTag = 'M' + 'I' << 8 + 'D' << 16 + 'S' << 24
	TAG_UAID	HandshakeTag = 'U' + 'A' << 8 + 'I' << 16 + 'D' << 24
	TAG_TCID	HandshakeTag = 'T' + 'C' << 8 + 'I' << 16 + 'D' << 24
	TAG_PDMD	HandshakeTag = 'P' + 'D' << 8 + 'M' << 16 + 'D' << 24
	TAG_SRBF	HandshakeTag = 'S' + 'R' << 8 + 'B' << 16 + 'F' << 24
	TAG_ICSL	HandshakeTag = 'I' + 'C' << 8 + 'S' << 16 + 'L' << 24
	TAG_NONP	HandshakeTag = 'N' + 'O' << 8 + 'N' << 16 + 'P' << 24
	TAG_SCLS	HandshakeTag = 'S' + 'C' << 8 + 'L' << 16 + 'S' << 24
	TAG_CSCT	HandshakeTag = 'C' + 'S' << 8 + 'C' << 16 + 'T' << 24
	TAG_COPT	HandshakeTag = 'C' + 'O' << 8 + 'P' << 16 + 'T' << 24
	TAG_CFCW	HandshakeTag = 'C' + 'F' << 8 + 'C' << 16 + 'W' << 24
	TAG_SFCW	HandshakeTag = 'S' + 'F' << 8 + 'C' << 16 + 'W' << 24
	TAG_FHL2	HandshakeTag = 'F' + 'H' << 8 + 'L' << 16 + '2' << 24
	TAG_NSTP	HandshakeTag = 'N' + 'S' << 8 + 'T' << 16 + 'P' << 24
	TAG_STK		HandshakeTag = 'S' + 'T' << 8 + 'K' << 16
	TAG_SNO		HandshakeTag = 'S' + 'N' << 8 + 'O' << 16
	TAG_PROF	HandshakeTag = 'P' + 'R' << 8 + 'O' << 16 + 'F' << 24
	TAG_NONC	HandshakeTag = 'N' + 'O' << 8 + 'N' << 16 + 'C' << 24
	TAG_XLCT	HandshakeTag = 'X' + 'L' << 8 + 'C' << 16 + 'T' << 24
	TAG_SCID	HandshakeTag = 'S' + 'C' << 8 + 'I' << 16 + 'D' << 24
	TAG_KEXS	HandshakeTag = 'K' + 'E' << 8 + 'X' << 16 + 'S' << 24
	TAG_AEAD	HandshakeTag = 'A' + 'E' << 8 + 'A' << 16 + 'D' << 24
	TAG_PUBS	HandshakeTag = 'P' + 'U' << 8 + 'B' << 16 + 'S' << 24
	TAG_OBIT	HandshakeTag = 'O' + 'B' << 8 + 'I' << 16 + 'T' << 24
	TAG_EXPY	HandshakeTag = 'E' + 'X' << 8 + 'P' << 16 + 'Y' << 24
	TAG_CERT	HandshakeTag = 'C' + 'E' << 8 + 'R' << 16 + 'T' << 24
	TAG_SHLO	HandshakeTag = 'S' + 'H' << 8 + 'L' << 16 + 'O' << 24
	TAG_PRST	HandshakeTag = 'P' + 'R' << 8 + 'S' << 16 + 'T' << 24
	TAG_RSEQ	HandshakeTag = 'R' + 'S' << 8 + 'E' << 16 + 'Q' << 24
	TAG_RNON	HandshakeTag = 'R' + 'N' << 8 + 'O' << 16 + 'N' << 24
)