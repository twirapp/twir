export type DonationAlertsMessage = {
	id: number;
	name: string;
	username?: string | null;
	message: string | null;
	message_type: 'text' | 'audio';
	payin_system: null | any;
	amount: number;
	currency: string;
	amount_in_user_currency: number;
	recipient_name: string;
	recipient: {
		user_id: number;
		code: string;
		name: string;
		avatar: string;
	};
	created_at: string;
	shown_at: null | any;
	reason: string;
};
