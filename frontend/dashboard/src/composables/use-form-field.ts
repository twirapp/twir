import { type Ref, ref, watch } from 'vue'

export function useFormField<T>(name: string, initialValue: T) {
	const fieldRef = ref<HTMLInputElement | null>(null)
	const fieldModel = ref(initialValue) as Ref<T>

	watch(fieldRef, () => {
		if (!fieldRef.value) return
		fieldRef.value.name = name
	})

	function reset() {
		fieldModel.value = initialValue
	}

	return {
		fieldRef,
		fieldModel,
		reset
	}
}
