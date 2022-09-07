<template>
  <section class="bg-black-15 py-24">
    <div class="relative container max-w-[1020px]">
      <div class="flex justify-between mb-10">
        <h2 class="text-5xl font-semibold leading-[130%] text-center">
          {{ t('sections.team.title') }}
        </h2>
        <p class="text-[17px] text-gray-70 max-w-[600px]">
          {{ t('sections.team.description') }}
        </p>
      </div>

      <ul class="member-list pt-10 border-t border-t-gray-30">
        <li v-for="member in teamMembers" :key="member.id" class="flex">
          <TeamMemberCard
            :isFounder="member.isFounder"
            :role="member.role"
            :name="member.name"
            :socials="member.socials"
          />
        </li>
      </ul>
    </div>
  </section>
</template>

<script lang="ts" setup>
import TeamMemberCard from '@/components/landing/TeamMemberCard.vue';
import type { TeamMember } from '@/types/teamMember';
import { useTranslation } from '@/utils/locales.js';

defineProps<{ teamMembers: TeamMember[] }>();

const t = useTranslation<'landing'>();
</script>

<style lang="postcss">
.member-list {
  @apply inline-grid grid-flow-col w-full;

  & > :not(:last-child) {
    @apply border-r border-r-gray-30;
  }
}
</style>
