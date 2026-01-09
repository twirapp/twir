<script setup lang="ts">

import { useRouter } from 'vue-router'

import { useBeRightBackOverlayApi } from '#layers/dashboard/api/overlays-be-right-back'
import KeyboardAndMouseIcon from '~/assets/overlays/keyboard-and-mouse.svg?use'
import Card from '#layers/dashboard/components/overlays/card.vue'

const { t } = useI18n()
const api = useBeRightBackOverlayApi()
const { data: settings, error, fetching } = api.useQueryBeRightBack()

const router = useRouter()
</script>

<template>
	<card
		:icon="KeyboardAndMouseIcon"
		title="AFK overlay"
		:description="t('overlays.brb.description')"
		overlay-path="brb"
		:copy-disabled="!settings || !!error || !!fetching"
		@open-settings="router.push({ name: 'BrbOverlay' })"
	>
	</card>
</template>
