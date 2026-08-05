package shortenedurls

import (
	"context"
	"testing"

	shortlinksbanneduapresetpatternsrepository "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_preset_patterns"
	shortlinkslinkbannedusaragentsrepository "github.com/twirapp/twir/libs/repositories/short_links_link_banned_user_agents"
	shortlinkslinkpresetsrepository "github.com/twirapp/twir/libs/repositories/short_links_link_presets"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
)

type linkBannedUserAgentsStub struct {
	patterns         []shortlinkslinkbannedusaragentsrepository.BannedUserAgent
	getByLinkIDCalls []string
	createCalls      int
}

func (s *linkBannedUserAgentsStub) GetByLinkID(
	_ context.Context,
	linkID string,
) ([]shortlinkslinkbannedusaragentsrepository.BannedUserAgent, error) {
	s.getByLinkIDCalls = append(s.getByLinkIDCalls, linkID)
	return s.patterns, nil
}

func (s *linkBannedUserAgentsStub) Create(
	_ context.Context,
	input shortlinkslinkbannedusaragentsrepository.CreateInput,
) (shortlinkslinkbannedusaragentsrepository.BannedUserAgent, error) {
	s.createCalls++
	return shortlinkslinkbannedusaragentsrepository.BannedUserAgent{
		ID:      "ban-id",
		LinkID:  input.LinkID,
		Pattern: input.Pattern,
	}, nil
}

func (s *linkBannedUserAgentsStub) Delete(_ context.Context, _, _ string) error {
	return nil
}

type linkPresetsStub struct {
	presets          []shortlinkslinkpresetsrepository.LinkPreset
	getByLinkIDCalls []string
}

func (s *linkPresetsStub) GetByLinkID(
	_ context.Context,
	linkID string,
) ([]shortlinkslinkpresetsrepository.LinkPreset, error) {
	s.getByLinkIDCalls = append(s.getByLinkIDCalls, linkID)
	return s.presets, nil
}

func (s *linkPresetsStub) GetByPresetID(
	_ context.Context,
	_ string,
) ([]shortlinkslinkpresetsrepository.LinkPreset, error) {
	return nil, nil
}

func (s *linkPresetsStub) GetLinksByPresetID(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

func (s *linkPresetsStub) Create(
	_ context.Context,
	_ shortlinkslinkpresetsrepository.CreateInput,
) (shortlinkslinkpresetsrepository.LinkPreset, error) {
	return shortlinkslinkpresetsrepository.LinkPreset{}, nil
}

func (s *linkPresetsStub) Delete(_ context.Context, _ string) error {
	return nil
}

func (s *linkPresetsStub) DeleteByLinkAndPreset(_ context.Context, _, _ string) error {
	return nil
}

type presetPatternsStub struct {
	patternsByPreset map[string][]shortlinksbanneduapresetpatternsrepository.Pattern
}

func (s *presetPatternsStub) GetByPresetID(
	_ context.Context,
	presetID string,
) ([]shortlinksbanneduapresetpatternsrepository.Pattern, error) {
	return s.patternsByPreset[presetID], nil
}

func (s *presetPatternsStub) Create(
	_ context.Context,
	_ shortlinksbanneduapresetpatternsrepository.CreateInput,
) (shortlinksbanneduapresetpatternsrepository.Pattern, error) {
	return shortlinksbanneduapresetpatternsrepository.Pattern{}, nil
}

func (s *presetPatternsStub) Delete(_ context.Context, _, _ string) error {
	return nil
}

func newBanCheckService(
	direct *linkBannedUserAgentsStub,
	linkPresets *linkPresetsStub,
	presetPatterns *presetPatternsStub,
) *Service {
	return &Service{
		linkBannedUserAgentsRepository: direct,
		linkPresetsRepository:          linkPresets,
		presetPatternsRepository:       presetPatterns,
	}
}

func TestIsUserAgentBannedMatchesCaseInsensitive(t *testing.T) {
	direct := &linkBannedUserAgentsStub{
		patterns: []shortlinkslinkbannedusaragentsrepository.BannedUserAgent{
			{Pattern: "chatterino.+|ffzbot.+"},
		},
	}
	service := newBanCheckService(direct, &linkPresetsStub{}, &presetPatternsStub{})

	banned, err := service.IsUserAgentBanned(
		context.Background(),
		model.ShortenedUrl{ID: "link-pk", ShortID: "onlyfans"},
		"Chatterino/2.5.3 (https://chatterino.com)",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !banned {
		t.Fatal("expected mixed-case Chatterino user agent to be banned by lowercase pattern")
	}
}

func TestIsUserAgentBannedUsesLinkPrimaryKey(t *testing.T) {
	direct := &linkBannedUserAgentsStub{}
	linkPresets := &linkPresetsStub{}
	service := newBanCheckService(direct, linkPresets, &presetPatternsStub{})

	link := model.ShortenedUrl{ID: "uuid-primary-key", ShortID: "h7Hs1"}
	if _, err := service.IsUserAgentBanned(context.Background(), link, "Mozilla/5.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, got := range direct.getByLinkIDCalls {
		if got != link.ID {
			t.Fatalf("direct patterns lookup used %q, want link primary key %q", got, link.ID)
		}
	}
	for _, got := range linkPresets.getByLinkIDCalls {
		if got != link.ID {
			t.Fatalf("link presets lookup used %q, want link primary key %q", got, link.ID)
		}
	}
	if len(direct.getByLinkIDCalls) == 0 || len(linkPresets.getByLinkIDCalls) == 0 {
		t.Fatal("expected both repositories to be queried")
	}
}

func TestIsUserAgentBannedMatchesPresetPatterns(t *testing.T) {
	linkPresets := &linkPresetsStub{
		presets: []shortlinkslinkpresetsrepository.LinkPreset{
			{LinkID: "link-pk", PresetID: "preset-1"},
		},
	}
	presetPatterns := &presetPatternsStub{
		patternsByPreset: map[string][]shortlinksbanneduapresetpatternsrepository.Pattern{
			"preset-1": {{Pattern: "twitchlib"}},
		},
	}
	service := newBanCheckService(&linkBannedUserAgentsStub{}, linkPresets, presetPatterns)

	banned, err := service.IsUserAgentBanned(
		context.Background(),
		model.ShortenedUrl{ID: "link-pk", ShortID: "abc12"},
		"TwitchLib/1.0",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !banned {
		t.Fatal("expected user agent to be banned via preset pattern")
	}
}

func TestIsUserAgentBannedSkipsInvalidPatterns(t *testing.T) {
	direct := &linkBannedUserAgentsStub{
		patterns: []shortlinkslinkbannedusaragentsrepository.BannedUserAgent{
			{Pattern: "(["},
			{Pattern: "chatterino"},
		},
	}
	service := newBanCheckService(direct, &linkPresetsStub{}, &presetPatternsStub{})

	banned, err := service.IsUserAgentBanned(
		context.Background(),
		model.ShortenedUrl{ID: "link-pk"},
		"Chatterino/2.5.3",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !banned {
		t.Fatal("expected invalid pattern to be skipped and valid one to match")
	}
}

func TestIsUserAgentBannedEmptyUserAgent(t *testing.T) {
	direct := &linkBannedUserAgentsStub{}
	linkPresets := &linkPresetsStub{}
	service := newBanCheckService(direct, linkPresets, &presetPatternsStub{})

	banned, err := service.IsUserAgentBanned(
		context.Background(),
		model.ShortenedUrl{ID: "link-pk"},
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if banned {
		t.Fatal("expected empty user agent to never be banned")
	}
	if len(direct.getByLinkIDCalls) != 0 || len(linkPresets.getByLinkIDCalls) != 0 {
		t.Fatal("expected no repository calls for empty user agent")
	}
}

func TestIsUserAgentBannedWithoutLinkPrimaryKey(t *testing.T) {
	direct := &linkBannedUserAgentsStub{}
	linkPresets := &linkPresetsStub{}
	service := newBanCheckService(direct, linkPresets, &presetPatternsStub{})

	banned, err := service.IsUserAgentBanned(
		context.Background(),
		model.ShortenedUrl{ShortID: "abc12"},
		"Chatterino/2.5.3",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if banned {
		t.Fatal("expected link without primary key to never be banned")
	}
	if len(direct.getByLinkIDCalls) != 0 || len(linkPresets.getByLinkIDCalls) != 0 {
		t.Fatal("expected no repository calls without link primary key")
	}
}

func TestCreateLinkBannedUserAgentRejectsInvalidPattern(t *testing.T) {
	direct := &linkBannedUserAgentsStub{}
	service := newBanCheckService(direct, &linkPresetsStub{}, &presetPatternsStub{})

	_, err := service.CreateLinkBannedUserAgent(
		context.Background(),
		shortlinkslinkbannedusaragentsrepository.CreateInput{
			LinkID:  "link-pk",
			Pattern: "([",
		},
	)
	if err == nil {
		t.Fatal("expected invalid pattern to be rejected")
	}
	if direct.createCalls != 0 {
		t.Fatalf("expected no repository calls, got %d", direct.createCalls)
	}
}
