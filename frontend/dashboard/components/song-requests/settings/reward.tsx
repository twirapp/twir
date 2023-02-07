import { Avatar, Group, Text } from '@mantine/core';
import { forwardRef } from 'react';

export interface RewardItemProps extends React.ComponentPropsWithoutRef<'div'> {
  label: string;
  description: string;
  image?: string
  value: string
}

export const RewardItem = forwardRef<HTMLDivElement, RewardItemProps>(
  ({ label, description, image, ...others }: RewardItemProps, ref) => (
    <div ref={ref} {...others}>
      <Group noWrap>
          {image && <Avatar src={image} />}
        <div>

          <Text size="sm">{label}</Text>
          <Text size="xs" opacity={0.65}>
            {description}
          </Text>
        </div>
      </Group>
    </div>
  ),
);