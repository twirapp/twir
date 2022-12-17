import { Grid } from '@mantine/core';
import {
  IconLambda,
  IconLetterCaseUpper,
  IconLink,
  IconMoodSmile,
  IconPlaylistX,
  IconTextWrapDisabled,
} from '@tabler/icons';

import { ModerationCard } from '../components/moderation/card';

export default function () {
  return (
    <Grid justify="center">
      <Grid.Col xs={12} sm={12} md={5} lg={5} xl={5}>
        <ModerationCard title="Links" icon={IconLink} iconColor="cyan">
          123
        </ModerationCard>
      </Grid.Col>
      <Grid.Col xs={12} sm={12} md={5} lg={5} xl={5}>
        <ModerationCard title="Caps" icon={IconLetterCaseUpper} iconColor="orange">
          Caps
        </ModerationCard>
      </Grid.Col>
      <Grid.Col xs={12} sm={12} md={5} lg={5} xl={5}>
        <ModerationCard title="Emotes" icon={IconMoodSmile} iconColor="yellow">
          123
        </ModerationCard>
      </Grid.Col>
      <Grid.Col xs={12} sm={12} md={5} lg={5} xl={5}>
        <ModerationCard title="LongMessage" icon={IconTextWrapDisabled}>
          123
        </ModerationCard>
      </Grid.Col>
      <Grid.Col xs={12} sm={12} md={5} lg={5} xl={5} color="red">
        <ModerationCard title="Blacklists" icon={IconPlaylistX}>
          123
        </ModerationCard>
      </Grid.Col>
      <Grid.Col xs={12} sm={12} md={5} lg={5} xl={5}>
        <ModerationCard title="Symbols" icon={IconLambda} iconColor="green">
          123
        </ModerationCard>
      </Grid.Col>
    </Grid>
  );
}
