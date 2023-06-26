import { createFormContext } from '@mantine/form';
import { YouTubeSettings } from '@twir/types/api';

export const [YouTubeSettingsFormProvider, useYouTubeSettingsFormContext, useYouTubeSettingsForm] = createFormContext<YouTubeSettings>();
