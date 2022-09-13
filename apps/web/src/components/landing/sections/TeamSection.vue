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
        <li v-for="(member, memberId) in teamMembers" :key="memberId" class="flex">
          <TeamMemberCard
            :isFounder="member.isFounder"
            :role="teamMembersLocale[memberId]"
            :name="member.name"
            :socials="member.socials"
          />
        </li>
      </ul>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import TeamMemberCard from '@/components/landing/TeamMemberCard.vue';
import { teamMembers, TeamMemberLocale } from '@/data/landing/team.js';
import useTranslation from '@/hooks/useTranslation';

const t = useTranslation<'landing'>();
const { tm } = useI18n();

const teamMembersLocale = computed(() => tm('sections.team.members') as TeamMemberLocale);
</script>

<style lang="postcss">
.member-list {
  @apply inline-grid grid-flow-col w-full;

  & > :not(:last-child) {
    @apply border-r border-r-gray-30;
  }
}
</style>
