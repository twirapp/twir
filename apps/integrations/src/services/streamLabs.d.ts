export type StreamLabsEvent = {
	type: 'donation';
	message: StreamLabsMessage[];
	for: string;
	event_id: string;
};

export type StreamLabsMessage = {
	name: string;
	isTest: boolean;
	formatted_amount: string;
	amount: number;
	message: string | null;
	currency: string;
	to: { name: string };
	from: string;
	from_user_id: number;
	_id: string;
	priority: number;
};
