package compressors

import (
	"AlgosLab1Sem4/internal/algorithms"
	pipeline "AlgosLab1Sem4/internal/pipline"
)

type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
	GetName() string
}

type HACompressor struct {
	pipe *pipeline.Pipeline
}

func (c *HACompressor) GetName() string {
	return "HACompressor"
}

func NewHACompressor() *HACompressor {
	return &HACompressor{
		pipe: pipeline.NewPipeline(&algorithms.HA{}),
	}
}

func (c *HACompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *HACompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type RLECompressor struct {
	pipe *pipeline.Pipeline
}

func (c *RLECompressor) GetName() string {
	return "RLECompressor"
}

func NewRLECompressor() *RLECompressor {
	return &RLECompressor{
		pipe: pipeline.NewPipeline(&algorithms.RLE{}),
	}
}

func (c *RLECompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *RLECompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type BWTRLECompressor struct {
	pipe *pipeline.Pipeline
}

func (c *BWTRLECompressor) GetName() string {
	return "BWTRLECompressor"
}

func NewBWTRLECompressor(blockSize int) *BWTRLECompressor {
	return &BWTRLECompressor{
		pipe: pipeline.NewPipeline(&algorithms.BWT{BlockSize: blockSize}, &algorithms.RLE{}),
	}
}

func (c *BWTRLECompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *BWTRLECompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type BWTMTFHACompressor struct {
	pipe *pipeline.Pipeline
}

func (c *BWTMTFHACompressor) GetName() string {
	return "BWTMTFHACompressor"
}

func NewBWTMTFHACompressor(blockSize int) *BWTMTFHACompressor {
	return &BWTMTFHACompressor{
		pipe: pipeline.NewPipeline(&algorithms.BWT{BlockSize: blockSize}, &algorithms.MTF{}, &algorithms.HA{}),
	}
}

func (c *BWTMTFHACompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *BWTMTFHACompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type BWTMTFRLEHACompressor struct {
	pipe *pipeline.Pipeline
}

func (c *BWTMTFRLEHACompressor) GetName() string {
	return "BWTMTFRLEHACompressor"
}

func NewBWTMTFRLEHACompressor(blockSize int) *BWTMTFRLEHACompressor {
	return &BWTMTFRLEHACompressor{
		pipe: pipeline.NewPipeline(&algorithms.BWT{BlockSize: blockSize}, &algorithms.MTF{}, &algorithms.RLE{}, &algorithms.HA{}),
	}
}

func (c *BWTMTFRLEHACompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *BWTMTFRLEHACompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type LZ77Compressor struct {
	pipe *pipeline.Pipeline
}

func (c *LZ77Compressor) GetName() string {
	return "LZ77Compressor"
}

func NewLZ77Compressor(buffer int) *LZ77Compressor {
	return &LZ77Compressor{
		pipe: pipeline.NewPipeline(&algorithms.LZ77{Buffer: buffer}),
	}
}

func (c *LZ77Compressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *LZ77Compressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type LZ77HACompressor struct {
	pipe *pipeline.Pipeline
}

func (c *LZ77HACompressor) GetName() string {
	return "LZ77HACompressor"
}

func NewLZ77HACompressor(buffer int) *LZ77HACompressor {
	return &LZ77HACompressor{
		pipe: pipeline.NewPipeline(&algorithms.LZ77{Buffer: buffer}, &algorithms.HA{}),
	}
}

func (c *LZ77HACompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *LZ77HACompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type LZ78Compressor struct {
	pipe *pipeline.Pipeline
}

func (c *LZ78Compressor) GetName() string {
	return "LZ78Compressor"
}

func NewLZ78Compressor() *LZ78Compressor {
	return &LZ78Compressor{
		pipe: pipeline.NewPipeline(&algorithms.LZ78{}),
	}
}

func (c *LZ78Compressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *LZ78Compressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type LZ78HACompressor struct {
	pipe *pipeline.Pipeline
}

func (c *LZ78HACompressor) GetName() string {
	return "LZ78HACompressor"
}

func NewLZ78HACompressor() *LZ78HACompressor {
	return &LZ78HACompressor{
		pipe: pipeline.NewPipeline(&algorithms.LZ78{}, &algorithms.HA{}),
	}
}

func (c *LZ78HACompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *LZ78HACompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}

type BWTMTFCompressor struct {
	pipe *pipeline.Pipeline
}

func NewBWTMTFCompressor(blockSize int) *BWTMTFCompressor {
	return &BWTMTFCompressor{
		pipe: pipeline.NewPipeline(&algorithms.BWT{BlockSize: blockSize}, &algorithms.HA{}),
	}
}

func (c *BWTMTFCompressor) Compress(data []byte) ([]byte, error) {
	return c.pipe.Encode(data)
}

func (c *BWTMTFCompressor) Decompress(data []byte) ([]byte, error) {
	return c.pipe.Decode(data)
}
