import { createFormContext } from '@mantine/form';
import { YouTubeSettings } from '@tsuwari/types/api';

export const [YouTubeSettingsFormProvider, useYouTubeSettingsFormContext, useYouTubeSettingsForm] = createFormContext<YouTubeSettings>();