package matchers

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
)

// Java bytecode and Mach-O binaries share the same magic number.
// More info here https://github.com/threatstack/libmagic/blob/master/magic/Magdir/cafebabe
func classOrMachOFat(in []byte) bool {
	// There should be at least 8 bytes for both of them because the only way to
	// quickly distinguish them is by comparing byte at position 7
	if len(in) < 8 {
		return false
	}

	return bytes.HasPrefix(in, []byte{0xCA, 0xFE, 0xBA, 0xBE})
}

// Class matches a java class file.
func Class(in []byte, _ uint32) bool {
	return classOrMachOFat(in) && in[7] > 30
}

// Lnk matches Microsoft lnk binary format.
func Lnk(in []byte, _ uint32) bool {
	return bytes.HasPrefix(in, []byte{0x4C, 0x00, 0x00, 0x00, 0x01, 0x14, 0x02, 0x00})
}

// MachO matches Mach-O binaries format.
func MachO(in []byte, _ uint32) bool {
	if classOrMachOFat(in) && in[7] < 20 {
		return true
	}

	if len(in) < 4 {
		return false
	}

	be := binary.BigEndian.Uint32(in)
	le := binary.LittleEndian.Uint32(in)

	return be == macho.Magic32 || le == macho.Magic32 || be == macho.Magic64 || le == macho.Magic64
}

// Swf matches an Adobe Flash swf file.
func Swf(in []byte, _ uint32) bool {
	return bytes.HasPrefix(in, []byte("CWS")) ||
		bytes.HasPrefix(in, []byte("FWS")) ||
		bytes.HasPrefix(in, []byte("ZWS"))
}

// Wasm matches a web assembly File Format file.
func Wasm(in []byte, _ uint32) bool {
	return bytes.HasPrefix(in, []byte{0x00, 0x61, 0x73, 0x6D})
}

// Dbf matches a dBase file.
// https://www.dbase.com/Knowledgebase/INT/db7_file_fmt.htm
func Dbf(in []byte, _ uint32) bool {
	if len(in) < 4 {
		return false
	}

	// 3rd and 4th bytes contain the last update month and day of month
	if !(0 < in[2] && in[2] < 13 && 0 < in[3] && in[3] < 32) {
		return false
	}

	// dbf type is dictated by the first byte
	dbfTypes := []byte{
		0x02, 0x03, 0x04, 0x05, 0x30, 0x31, 0x32, 0x42, 0x62, 0x7B, 0x82,
		0x83, 0x87, 0x8A, 0x8B, 0x8E, 0xB3, 0xCB, 0xE5, 0xF5, 0xF4, 0xFB,
	}
	for _, b := range dbfTypes {
		if in[0] == b {
			return true
		}
	}

	return false
}

// Exe matches a Windows/DOS executable file.
func Exe(in []byte, _ uint32) bool {
	return bytes.HasPrefix(in, []byte{0x4D, 0x5A})
}

// Elf matches an Executable and Linkable Format file.
func Elf(in []byte, _ uint32) bool {
	return bytes.HasPrefix(in, []byte{0x7F, 0x45, 0x4C, 0x46})
}

// ElfObj matches an object file.
func ElfObj(in []byte, _ uint32) bool {
	return len(in) > 17 && ((in[16] == 0x01 && in[17] == 0x00) ||
		(in[16] == 0x00 && in[17] == 0x01))
}

// ElfExe matches an executable file.
func ElfExe(in []byte, _ uint32) bool {
	return len(in) > 17 && ((in[16] == 0x02 && in[17] == 0x00) ||
		(in[16] == 0x00 && in[17] == 0x02))
}

// ElfLib matches a shared library file.
func ElfLib(in []byte, _ uint32) bool {
	return len(in) > 17 && ((in[16] == 0x03 && in[17] == 0x00) ||
		(in[16] == 0x00 && in[17] == 0x03))
}

// ElfDump matches a core dump file.
func ElfDump(in []byte, _ uint32) bool {
	return len(in) > 17 && ((in[16] == 0x04 && in[17] == 0x00) ||
		(in[16] == 0x00 && in[17] == 0x04))
}

// Dcm matches a DICOM medical format file.
func Dcm(in []byte, _ uint32) bool {
	return len(in) > 131 &&
		bytes.Equal(in[128:132], []byte{0x44, 0x49, 0x43, 0x4D})
}

// Nes matches a Nintendo Entertainment system ROM file.
func Nes(in []byte, _ uint32) bool {
	return bytes.HasPrefix(in, []byte{0x4E, 0x45, 0x53, 0x1A})
}

// Marc matches a MARC21 (MAchine-Readable Cataloging) file.
func Marc(in []byte, _ uint32) bool {
	// File is at least 24 bytes ("leader" field size)
	if len(in) < 24 {
		return false
	}

	// Fixed bytes at offset 20
	if !bytes.Equal(in[20:24], []byte("4500")) {
		return false
	}

	// First 5 bytes are ASCII digits
	for i := 0; i < 5; i++ {
		if in[i] < '0' || in[i] > '9' {
			return false
		}
	}

	// Field terminator is present
	return bytes.Contains(in, []byte{0x1E})
}

// TzIf matches a Time Zone Information Format (TZif) file.
func TzIf(in []byte, _ uint32) bool {
	return len(in) > 4 && bytes.HasPrefix(in, []byte("TZif"))
}
