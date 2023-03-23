import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { EventOperation } from './EventOperation';

export enum FilterType {
  EQUALS = 'EQUALS',
  NOT_EQUALS = 'NOT_EQUALS',
  CONTAINS = 'CONTAINS',
  NOT_CONTAINS = 'NOT_CONTAINS',
  STARTS_WITH = 'STARTS_WITH',
  ENDS_WITH = 'ENDS_WITH',
  GREATER_THAN = 'GREATER_THAN',
  LESS_THAN = 'LESS_THAN',
  GREATER_THAN_OR_EQUALS = 'GREATER_THAN_OR_EQUALS',
  LESS_THAN_OR_EQUALS = 'LESS_THAN_OR_EQUALS',
  IS_EMPTY = 'IS_EMPTY',
  IS_NOT_EMPTY = 'IS_NOT_EMPTY',
}

@Entity('channels_events_operations_filters')
export class EventOperationFilter {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => EventOperation, _ => _.filters, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'operationId' })
  operation?: EventOperation;

  @Column('uuid')
  operationId: string;

  @Column('simple-enum', { enum: FilterType })
  type: FilterType;

  @Column('text')
  left: string;

  @Column('text')
  right: string;
}