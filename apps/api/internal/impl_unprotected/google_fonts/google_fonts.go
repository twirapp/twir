package google_fonts

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/libs/grpc/generated/api/google_fonts_unprotected"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoogleFonts struct {
	*impl_deps.Deps
}

type googleResponseFont struct {
	Family  string            `json:"family"`
	Files   map[string]string `json:"files"`
	SubSets []string          `json:"subsets"`
}

type googleResponse struct {
	Fonts []googleResponseFont `json:"items"`
}

func (c *GoogleFonts) GetGoogleFonts(
	ctx context.Context,
	_ *emptypb.Empty,
) (*google_fonts_unprotected.Fonts, error) {
	response := googleResponse{}

	r, err := req.
		SetSuccessResult(&response).
		Get("https://www.googleapis.com/webfonts/v1/webfonts?key=" + c.Config.GoogleFontsApiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get google fonts: %w", err)
	}
	if !r.IsSuccessState() {
		return nil, fmt.Errorf("failed to get google fonts: %s", r.String())
	}

	fonts := make([]*google_fonts_unprotected.Font, 0, len(response.Fonts))
	for _, font := range response.Fonts {
		fontFiles := make([]*google_fonts_unprotected.Font_File, 0, len(font.Files))
		for key, value := range font.Files {
			fontFiles = append(
				fontFiles, &google_fonts_unprotected.Font_File{
					Name: key,
					Url:  value,
				},
			)
		}

		fonts = append(
			fonts, &google_fonts_unprotected.Font{
				Family:  font.Family,
				Files:   fontFiles,
				Subsets: font.SubSets,
			},
		)
	}

	return &google_fonts_unprotected.Fonts{
		Fonts: fonts,
	}, nil
}
