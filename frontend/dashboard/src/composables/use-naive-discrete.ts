import {
	useDialog,
	useMessage,
	useNotification,
} from 'naive-ui';


export const useNaiveDiscrete = () => {
	const dialog = useDialog();
	const notification = useNotification();
	const message = useMessage();

	return {
		dialog,
		notification,
		message,
	};
};
