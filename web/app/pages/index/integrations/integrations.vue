<script setup lang="ts">
import { integrations } from '~/pages/index/integrations/integrations'
</script>

<template>
	<section id="integrations" class="bg-[#17171A]">
		<div class="container mx-auto">
			<div class="px-[32px] pb-[100px] pt-[140px]">
				<div class="flex flex-col flex-wrap items-center text-center gap-[20px] relative">
					<SvgoCube :fontControlled="false" :filled="true" class="absolute -top-[68px] h-[270px] w-[270px] stroke-[1]" aria-hidden />
					<h2 class="font-semibold text-[#B0ADFF] text-base uppercase text-center">
						Integrations
					</h2>
					<h3 class="text-4xl text-white font-bold z-10 max-w-2xl leading-[1.25] tracking-tight">
						Power-up your stream with seamless integrations
					</h3>
					<span class="text-xl text-[#ADB0B8] font-normal leading-normal max-w-3xl">
						Connect with your favorite tools and platforms to supercharge your stream's performance
						and engagement. Twir's integrations provide endless possibilities.
					</span>
				</div>
				<div class="flex flex-wrap justify-center mt-[64px] gap-[24px]" role="list">
					<a
						v-for="integration of integrations"
						:key="integration.href"
						:href="integration.href"
						target="_blank"
						rel="noopener noreferrer"
						aria-label="{integration.label}"
						data-microtip-position="bottom"
						role="tooltip"
					>
						<component
							:is="integration.icon"
							role="listitem"
							class="h-[64px] w-[64px] text-[#72757D] hover:text-[#606268] transition-colors m-1"
							:style="{
								width: integration.width ?? 'none',
							}"
							aria-label="{integration.label}"
						/>
					</a>
				</div>
			</div>
		</div>
	</section>
</template>
