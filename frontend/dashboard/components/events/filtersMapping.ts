import type { FilterType } from '@twir/typeorm/entities/events/EventOperationFilter';

export const filtersMapping: Record<keyof typeof FilterType, {
  description: string,
  withoutRight?: boolean,
}> = {
  CONTAINS: {
    description: 'Contains',
  },
  NOT_CONTAINS: {
    description: 'Not contains',
  },
  EQUALS: {
    description: 'Equals',
  },
  NOT_EQUALS: {
    description: 'Not equals',
  },
  STARTS_WITH: {
    description: 'Starts with',
  },
  ENDS_WITH: {
    description: 'Ends with',
  },
  GREATER_THAN_OR_EQUALS: {
    description: 'Greater than or equals',
  },
  LESS_THAN_OR_EQUALS: {
    description: 'Less than or equals',
  },
  GREATER_THAN: {
    description: 'Greater than',
  },
  LESS_THAN: {
    description: 'Less than',
  },
  IS_EMPTY: {
    description: 'Is empty',
    withoutRight: true,
  },
  IS_NOT_EMPTY: {
    description: 'Is not empty',
  },
};
