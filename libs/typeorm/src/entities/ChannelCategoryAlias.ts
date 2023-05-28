import {
	Column,
	Entity,
	Index,
	JoinColumn,
	ManyToOne,
	PrimaryGeneratedColumn,
	Unique,
} from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Channel } from './Channel';

@Entity('channels_categories_aliases', { schema: 'public' })
@Unique(['alias', 'channelId'])
export class ChannelCategoryAlias {
	@PrimaryGeneratedColumn('uuid')
	id: string;

	@Column('text', { name: 'category', nullable: false })
	category: string;

	@Column('text', { name: 'categoryId', nullable: false })
	categoryId: string;

	@Column('text', { name: 'alias', nullable: false })
	alias: string;

	@ManyToOne(() => Channel, (_) => _.categoriesAliases, {
		onDelete: 'RESTRICT',
		onUpdate: 'CASCADE',
	})
	@JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
	channel?: Channel;

	@Index()
	@Column('text', { name: 'channelId' })
	channelId: string;
}
