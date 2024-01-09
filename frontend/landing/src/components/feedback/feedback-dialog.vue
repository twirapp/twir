<script lang="ts" setup>
import type { Profile } from '@twir/grpc/generated/api/api/auth';
import { ref, watch } from 'vue';

import type { ReviewBody } from '@/pages/functions/feedback.js';

const dialog = ref<HTMLDialogElement>(null);

const props = defineProps<{
	profile?: Profile
}>();

function showDialog() {
	dialog.value.showModal();
	const body = document.body;
	body.style.overflowY = 'hidden';
}

watch(dialog, (v) => {
	if (!v) return;

	const body = document.body;

	const closeListener = () => {
		body.style.overflowY = 'auto';
	};

	v.addEventListener('close', closeListener);

	return () => {
		v.removeEventListener('close', closeListener);
	};
});

const reviewForm = ref<ReviewBody>({
	author: '',
	message: '',
});

const feedBackError = ref<string | null>(null);

async function sendFeedBack() {
	const { author, message } = reviewForm.value;

	const computedAuthor = props.profile?.login || author;

	if (!computedAuthor || !message) {
		feedBackError.value = 'Author and message required';
		return;
	}

	const req = await fetch('/functions/feedback', {
		method: 'POST',
		body: JSON.stringify({
			author: computedAuthor,
			message,
			profile: props.profile,
		}),
	});

	if (!req.ok) {
		const response = await req.json();
		if (response.error) {
			feedBackError.value = response.error;
		}
	} else {
		dialog.value?.close();
		reviewForm.value.author = '';
		reviewForm.value.message = '';
	}
}

function configForm(event: MouseEvent) {
	event.preventDefault();
	sendFeedBack();
}
</script>

<template>
	<button
		class="inline-flex px-6 py-3 font-semibold text-white rounded-lg bg-[#5D58F5] hover:bg-[#6964FF] transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4"
		@click="showDialog"
	>
		Leave feedback
	</button>
	<dialog id="feedback-dialog" ref="dialog">
		<div v-if="feedBackError" class="p-3 bg-red-500">
			{{ feedBackError }}
		</div>
		<form class="flex flex-col gap-2 p-5">
			<fieldset class="flex flex-col">
				<label for="author">Username</label>
				<input
					id="author"
					v-model="reviewForm.author"
					placeholder="Enter your twitch username"
					class="p-2 rounded-md shadow-sm outline-none"
					maxlength="25"
					autocomplete="off"
				/>
			</fieldset>
			<fieldset class="flex flex-col relative">
				<label for="message">Message</label>
				<textarea
					id="message"
					v-model="reviewForm.message"
					placeholder="Write your feedback"
					rows="6"
					class="p-2 rounded-md shadow-sm outline-none"
					maxlength="200"
					autocomplete="off"
				/>
				<span class="absolute px-2 py-1 text-xs text-white bg-[#5D58F5] rounded right-2 bottom-2">{{
					reviewForm.message.length
				}}/200</span>
			</fieldset>
			<div class="flex justify-end gap-2">
				<button
					value="cancel"
					formmethod="dialog"
					class="inline-flex px-3 py-1 font-semibold text-white rounded-lg bg-red-500 transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4"
				>
					Cancel
				</button>
				<button
					class="inline-flex px-3 py-1 font-semibold text-white rounded-lg bg-[#5D58F5] hover:bg-[#6964FF] transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4"
					@click="configForm"
				>
					Confirm
				</button>
			</div>
		</form>
	</dialog>
</template>

<style>
#feedback-dialog::backdrop {
	background-color: rgba(0, 0, 0, 9);
	opacity: 0.7;
}

#feedback-dialog {
	box-shadow: 0 4px 5px rgb(0 0 0 / 90%);
	border-radius: 5px;
	min-width: 500px;
}
</style>
