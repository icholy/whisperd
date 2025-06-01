package wav

import (
	"encoding/binary"
	"io"
)

// Header represents the header of a WAV file.
type Header struct {
	ChunkID       [4]byte // "RIFF"
	ChunkSize     uint32  // 36 + Subchunk2Size
	Format        [4]byte // "WAVE"
	Subchunk1ID   [4]byte // "fmt "
	Subchunk1Size uint32  // 16 for PCM
	AudioFormat   uint16  // 1 for PCM
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32 // SampleRate * NumChannels * BitsPerSample/8
	BlockAlign    uint16 // NumChannels * BitsPerSample/8
	BitsPerSample uint16
	Subchunk2ID   [4]byte // "data"
	Subchunk2Size uint32  // NumSamples * NumChannels * BitsPerSample/8
}

// Options specifies the WAV encoding parameters.
type Options struct {
	NumChannels int
	SampleRate  int
}

// Write writes the given PCM data as a WAV file to the provided writer using the specified options.
func Write(w io.Writer, pcm []byte, opt Options) error {
	h := Header{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + uint32(len(pcm)),
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   uint16(opt.NumChannels),
		SampleRate:    uint32(opt.SampleRate),
		ByteRate:      uint32(opt.SampleRate * opt.NumChannels * 2),
		BlockAlign:    uint16(opt.NumChannels * 2),
		BitsPerSample: 16,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: uint32(len(pcm)),
	}
	if err := binary.Write(w, binary.LittleEndian, h); err != nil {
		return err
	}
	_, err := w.Write(pcm)
	return err
}
