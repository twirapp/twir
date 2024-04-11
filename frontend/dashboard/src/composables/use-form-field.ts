import { ref, watch, type Ref } from 'vue';

export function useFormField<T>(name: string, initialValue: T) {
	const fieldRef = ref<HTMLInputElement | null>(null);
	const fieldModel = ref(initialValue) as Ref<T>;

	watch(fieldRef, () => {
		if (fieldRef.value) {
			fieldRef.value.name = name;
		}

		setValue(fieldModel.value);
	});

	function setValue(value: T) {
		if (fieldRef.value && typeof value === 'string') {
			fieldRef.value.value = value;
		}

		fieldModel.value = value;
	}

	function reset() {
		setValue(initialValue);
	}

	return {
		fieldRef,
		fieldModel,
		setValue,
		reset,
	};
}
