<template>
	<div class="flex flex-col gap-6 items-center mb-3">
		<div class="flex rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,10%)] p-5">
			<Icon name="lucide:image-up" fill="rgba(255,255,255,0.3)" class="w-8 h-8" />
		</div>
		<div class="flex flex-col gap-4 items-center justify-center">
			<div class="flex flex-col gap-0.5 items-center justify-center">
				<span class="text-4xl font-bold"> Twir </span>
				<span class="text-xl text-[hsl(240,11%,65%)]"> {{ $t('uploader.title') }} </span>
			</div>
			<span class="text-xl text-center w-3/4">
				{{ $t('uploader.tagline') }}
			</span>
			<NuxtLink
				to="/uploader/profile"
				class="flex items-center gap-2 px-4 py-2 rounded-lg border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] transition-colors text-sm font-medium"
			>
				<Icon name="lucide:user" class="w-4 h-4" />
				{{ $t('uploader.viewMyUploads') }}
			</NuxtLink>
		</div>
	</div>
</template>
