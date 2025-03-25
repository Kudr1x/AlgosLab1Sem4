package pipeline

import "AlgosLab1Sem4/internal/algorithms"

type Pipeline struct {
	transforms []algorithms.Algorithm
}

func NewPipeline(transforms ...algorithms.Algorithm) *Pipeline {
	return &Pipeline{transforms: transforms}
}

func (p *Pipeline) Encode(data []byte) ([]byte, error) {
	var err error
	for _, t := range p.transforms {
		data, err = t.Encode(data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (p *Pipeline) Decode(data []byte) ([]byte, error) {
	var err error
	for i := len(p.transforms) - 1; i >= 0; i-- {
		data, err = p.transforms[i].Decode(data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}
