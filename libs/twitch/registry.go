package twitch

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix"
)

func (runtime *twitchRuntime) registerBroadcaster(ctx context.Context, userID string) error {
	return runtime.register(ctx, broadcasterCredential, userID, nil)
}

func (runtime *twitchRuntime) registerBot(ctx context.Context, botID string) error {
	return runtime.register(ctx, botCredential, botID, nil)
}

func (runtime *twitchRuntime) registerChannelBot(ctx context.Context, botID string, channelID string) error {
	if botID == "" || channelID == "" {
		return fmt.Errorf("channel bot and channel IDs are required")
	}
	runtime.registrationMu.Lock()
	defer runtime.registrationMu.Unlock()
	if owner, found := runtime.channelBots[channelID]; found && owner != botID {
		return ErrChannelBotConflict
	}
	if err := runtime.registerLocked(ctx, botCredential, botID, []helix.Intent{channelBotIntent(channelID)}); err != nil {
		return err
	}
	runtime.channelBots[channelID] = botID
	return nil
}

func (runtime *twitchRuntime) removeChannelBot(ctx context.Context, botID string, channelID string) error {
	runtime.registrationMu.Lock()
	defer runtime.registrationMu.Unlock()
	if runtime.channelBots[channelID] != botID {
		return ErrChannelBotNotRegistered
	}
	value, found := runtime.registrations.Load(botID)
	if !found {
		return ErrChannelBotNotRegistered
	}
	registered := value.(*registeredCredential)
	remaining := copyIntents(registered.intents)
	delete(remaining, channelBotIntent(channelID))
	if err := runtime.replaceRegistrationLocked(ctx, botID, registered.kind, remaining); err != nil {
		return err
	}
	delete(runtime.channelBots, channelID)
	return nil
}

func (runtime *twitchRuntime) register(ctx context.Context, kind credentialKind, subjectID string, intents []helix.Intent) error {
	runtime.registrationMu.Lock()
	defer runtime.registrationMu.Unlock()
	return runtime.registerLocked(ctx, kind, subjectID, intents)
}

func (runtime *twitchRuntime) registerLocked(ctx context.Context, kind credentialKind, subjectID string, intents []helix.Intent) error {
	if subjectID == "" {
		return fmt.Errorf("Twitch credential ID is required")
	}
	value, found := runtime.registrations.Load(subjectID)
	if !found {
		return runtime.addRegistrationLocked(ctx, kind, subjectID, intents)
	}
	registered := value.(*registeredCredential)
	if registered.kind != kind {
		return ErrRegistrationConflict
	}
	merged := copyIntents(registered.intents)
	for _, intent := range intents {
		merged[intent] = struct{}{}
	}
	if len(merged) == len(registered.intents) {
		return nil
	}
	return runtime.replaceRegistrationLocked(ctx, subjectID, kind, merged)
}

func (runtime *twitchRuntime) addRegistrationLocked(ctx context.Context, kind credentialKind, subjectID string, intents []helix.Intent) error {
	binding, err := newCredentialBinding(ctx, runtime, kind, subjectID)
	if err != nil {
		return err
	}
	if err := runtime.registry.AddCoordinatedUser(ctx, subjectID, binding.loader, binding.hook, intents...); err != nil {
		return fmt.Errorf("add coordinated Twitch credential: %w", err)
	}
	runtime.registrations.Store(subjectID, &registeredCredential{kind: kind, intents: intentsToSet(intents)})
	return nil
}

func (runtime *twitchRuntime) replaceRegistrationLocked(ctx context.Context, subjectID string, kind credentialKind, intents map[helix.Intent]struct{}) error {
	if err := runtime.registry.RemoveUser(subjectID); err != nil {
		return fmt.Errorf("remove previous Twitch credential registration: %w", err)
	}
	runtime.registrations.Delete(subjectID)
	return runtime.addRegistrationLocked(ctx, kind, subjectID, intentsToSlice(intents))
}

func channelBotIntent(channelID string) helix.Intent {
	return helix.Intent("bot:channel:" + channelID)
}

func intentsToSet(intents []helix.Intent) map[helix.Intent]struct{} {
	result := make(map[helix.Intent]struct{}, len(intents))
	for _, intent := range intents {
		result[intent] = struct{}{}
	}
	return result
}

func copyIntents(intents map[helix.Intent]struct{}) map[helix.Intent]struct{} {
	result := make(map[helix.Intent]struct{}, len(intents))
	for intent := range intents {
		result[intent] = struct{}{}
	}
	return result
}

func intentsToSlice(intents map[helix.Intent]struct{}) []helix.Intent {
	result := make([]helix.Intent, 0, len(intents))
	for intent := range intents {
		result = append(result, intent)
	}
	return result
}
