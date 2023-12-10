package imaging

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traQ/testdata/gif"
	"github.com/traPtitech/traQ/utils"
)

const testdataFolder = "../../testdata/images/"

func mustOpen(path string) *os.File {
	fp, err := os.Open(testdataFolder + path)
	if err != nil {
		panic(err)
	}
	return fp
}

func setup() (Processor, *os.File) {
	processor := NewProcessor(Config{
		MaxPixels:        500 * 500,
		Concurrency:      1,
		ThumbnailMaxSize: image.Point{50, 50},
	})
	return processor, mustOpen("test.png")
}

func assertImg(t *testing.T, actualImg image.Image, expectedFilePath string) {
	actualImgBytesBuffer := &bytes.Buffer{}
	err := png.Encode(actualImgBytesBuffer, actualImg)
	if err != nil {
		panic(err)
	}
	actualImgBytes := actualImgBytesBuffer.Bytes()

	fpExpected := mustOpen(expectedFilePath)
	expectedImgBytes, err := io.ReadAll(fpExpected)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, expectedImgBytes, actualImgBytes)
}

func TestProcessorDefault_Thumbnail(t *testing.T) {
	t.Parallel()

	processor, fpActual := setup()
	defer fpActual.Close()
	actualImg, err := processor.Thumbnail(fpActual)
	assert.Nil(t, err)
	assertImg(t, actualImg, "test_thumbnail.png")
}

func TestProcessorDefault_Fit(t *testing.T) {
	t.Parallel()

	processor, fp := setup()
	defer fp.Close()
	actualImg, err := processor.Fit(fp, 100, 100)
	assert.Nil(t, err)
	assertImg(t, actualImg, "test_fit.png")
}

func setupCustomConc(conc int) Processor {
	processor := NewProcessor(Config{
		MaxPixels:        500 * 500,
		Concurrency:      conc,
		ThumbnailMaxSize: image.Point{X: 50, Y: 50},
	})
	return processor
}

func mustOpenGif(name string) fs.File {
	f, err := gif.FS.Open(name)
	if err != nil {
		panic(err)
	}
	return f
}

func TestProcessorDefault_FitAnimationGIF(t *testing.T) {
	t.Parallel()

	testConc := []int{1, 2, 3, 4, 5}

	test := []struct {
		name   string
		file   string
		reader io.Reader
		want   []byte
		err    error
	}{
		{
			name:   "invalid (empty)",
			reader: bytes.NewBufferString(""),
			want:   nil,
			err:    ErrInvalidImageSrc,
		},
		{
			name:   "invalid (invalid gif)",
			reader: io.LimitReader(mustOpenGif("cube.gif"), 10),
			want:   nil,
			err:    ErrInvalidImageSrc,
		},
		{
			name: "success (cube 正方形、透明ピクセルあり)",
			file: "cube.gif",
			want: utils.MustIoReaderToBytes(mustOpenGif("cube_resized.gif")),
			err:  nil,
		},
		{
			name: "success (miku 横長、差分最適化)",
			file: "miku.gif",
			want: utils.MustIoReaderToBytes(mustOpenGif("miku_resized.gif")),
			err:  nil,
		},
		{
			name: "success (parapara 正方形、差分最適化)",
			file: "parapara.gif",
			want: utils.MustIoReaderToBytes(mustOpenGif("parapara_resized.gif")),
			err:  nil,
		},
		{
			name: "success (miku2 縦長、差分最適化)",
			file: "miku2.gif",
			want: utils.MustIoReaderToBytes(mustOpenGif("miku2_resized.gif")),
			err:  nil,
		},
		{
			name: "success (rabbit 小サイズ)",
			file: "rabbit.gif",
			want: utils.MustIoReaderToBytes(mustOpenGif("rabbit_resized.gif")),
			err:  nil,
		},
	}

	for _, tt := range test {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, conc := range testConc {
				tt := tt
				conc := conc

				t.Run(fmt.Sprintf("%s, conc=%d", tt.name, conc), func(t *testing.T) {
					t.Parallel()

					processor := setupCustomConc(conc)
					if tt.file != "" { // ファイルはこのタイミングで開かないと正常なデータにならない
						tt.reader = mustOpenGif(tt.file)
					}

					actual, err := processor.FitAnimationGIF(tt.reader, 256, 256)
					if tt.err != nil {
						assert.Equal(t, tt.err, err)
					} else {
						assert.Nil(t, err)
						assert.Equal(t, tt.want, utils.MustIoReaderToBytes(actual))
					}
				})
			}
		})
	}
}
