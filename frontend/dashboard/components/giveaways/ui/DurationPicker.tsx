import { TextInput } from '@mantine/core';
import { useState } from 'react';

type Props = {
	onChange: (value: number | null) => void;
	label: string;
	placeholder: string;
	defaultValue?: string;
};

export const DurationPicker = ({ onChange, label, placeholder, defaultValue }: Props) => {
	const [inputString, setInputString] = useState(defaultValue);
	const [error, setError] = useState(false);

	const convertTimeString = (str: string) => {
		setError(false);
		setInputString(str);

		const timeRegex = /(\d+)\s*(s|m|h|d|w|mo|с|мин|ч|дн|н|мес)/g;
		const matches = str.match(timeRegex);

		if (!matches) {
			setError(true);
			onChange(null);
			return;
		}
		let totalMs = 0;

		matches.forEach((match) => {
			const value = parseInt(match, 10);
			const unit = match.slice(-1);

			switch (unit) {
				case 's':
				case 'с':
					totalMs += value * 1000;
					break;
				case 'm':
				case 'м':
					totalMs += value * 60000;
					break;
				case 'h':
				case 'ч':
					totalMs += value * 3600000;
					break;
				case 'd':
				case 'д':
					totalMs += value * 86400000;
					break;
				case 'w':
				case 'н':
					totalMs += value * 604800000;
					break;
				case 'mo':
				case 'мес':
					totalMs += value * 2592000000;
					break;
				default:
					break;
			}
		});

		onChange(totalMs);
	};

	return (
		<TextInput
			label={label}
			error={error}
			placeholder={placeholder}
			value={inputString}
			onChange={(e) => convertTimeString(e.target.value)}
		/>
	);
};
