package sd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"math/rand"

	"github.com/seasonjs/hf-hub/api"
	sd "github.com/seasonjs/stable-diffusion"
)

func newClient() *StableDiffusion {
	options := sd.DefaultOptions
	options.Threads = 6
	options.FreeParamsImmediately = false
	// options.GpuEnable = true

	rand.Seed(time.Now().Unix())

	model, err := sd.NewAutoModel(options)
	if err != nil {
		panic(err)
	}

	hapi, err := api.NewApi()
	if err != nil {
		panic(err)
	}

	modelPath, err := hapi.Model("justinpinkney/miniSD").Get("miniSD.ckpt")
	if err != nil {
		panic(err)
	}

	err = model.LoadFromFile(modelPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("model ready")

	return &StableDiffusion{
		model:     model,
		modelPath: modelPath,
		params: sd.FullParams{
			NegativePrompt:   "out of frame, lowers, text, error, cropped, worst quality, low quality, jpeg artifacts, ugly, duplicate, morbid, mutilated, out of frame, extra fingers, mutated hands, poorly drawn hands, poorly drawn face, mutation, deformed, blurry, dehydrated, bad anatomy, bad proportions, extra limbs, cloned face, disfigured, gross proportions, malformed limbs, missing arms, missing legs, extra arms, extra legs, fused fingers, too many fingers, long neck, username, watermark, signature",
			CfgScale:         7.0,
			Width:            64,
			Height:           64,
			SampleMethod:     sd.DPMPP2S_A,
			SampleSteps:      5,
			Strength:         0.4,
			Seed:             13,
			BatchCount:       1,
			OutputsImageType: sd.PNG,
		},
	}
}

type StableDiffusion struct {
	model     *sd.Model
	params    sd.FullParams
	modelPath string
}

func (s *StableDiffusion) Close() {
	fmt.Println(s.model.Close())
}

func (s *StableDiffusion) Inference(ctx context.Context, prompt string) ([]byte, error) {
	var b bytes.Buffer
	params := s.params
	params.Seed = int64(rand.Int())

	if err := s.model.Predict(prompt, params, []io.Writer{bufio.NewWriter(&b)}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
