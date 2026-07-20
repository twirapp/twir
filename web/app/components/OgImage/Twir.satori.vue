<script setup lang="ts">
interface Props {
	title?: string
	description?: string
}

withDefaults(defineProps<Props>(), {
	title: 'Twir',
	description: '',
})
</script>

<template>
	<div
		:style="{
			width: '1200px',
			height: '630px',
			display: 'flex',
			flexDirection: 'column',
			justifyContent: 'space-between',
			backgroundColor: '#09090B',
			padding: '64px',
			fontFamily: 'Inter, sans-serif',
			position: 'relative',
			overflow: 'hidden',
		}"
	>
		<div
			:style="{
				position: 'absolute',
				top: '-260px',
				left: '200px',
				right: '200px',
				height: '480px',
				borderRadius: '9999px',
				background:
					'radial-gradient(closest-side, rgba(93, 88, 245, 0.45), rgba(93, 88, 245, 0))',
			}"
		/>
		<div :style="{ display: 'flex', alignItems: 'center', gap: '20px' }">
			<img src="/logo-og.png" width="72" height="72" alt="" />
			<span
				:style="{
					color: 'white',
					fontSize: '40px',
					fontWeight: 700,
					letterSpacing: '0.04em',
					textTransform: 'uppercase',
				}"
			>
				Twir
			</span>
		</div>
		<div :style="{ display: 'flex', flexDirection: 'column', gap: '24px' }">
			<h1
				:style="{
					color: 'white',
					fontSize: title.length > 60 ? '52px' : '68px',
					fontWeight: 800,
					lineHeight: 1.1,
					margin: 0,
					letterSpacing: '-0.02em',
				}"
			>
				{{ title }}
			</h1>
			<p
				v-if="description"
				:style="{
					color: '#ADB0B8',
					fontSize: '30px',
					lineHeight: 1.4,
					margin: 0,
					maxWidth: '960px',
				}"
			>
				{{ description }}
			</p>
		</div>
		<div :style="{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }">
			<span :style="{ color: '#B0ADFF', fontSize: '28px', fontWeight: 600 }">twir.app</span>
			<span :style="{ color: '#6b6e75', fontSize: '24px' }">Twitch &amp; Kick chat bot</span>
		</div>
	</div>
</template>
