export type Donate = {
	twitchUserId: string;
	amount: number | string;
	currency: string;
	message?: string | null;
	userName?: string | null;
}

export const onDonation = (donate: Donate) => Promise<void>;
