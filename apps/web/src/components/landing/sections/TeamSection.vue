<template>
  <section class="bg-black-15 min-lg:py-24 min-md:py-20 py-14 min-md:px-10 min-sm:px-8 px-5">
    <div class="relative container max-w-[1020px]">
      <div
        class="inline-grid w-full min-lg:grid-flow-col gap-x-12 justify-between min-md:mb-10 mb-8"
      >
        <h2
          class="
            min-lg:text-[48px] min-md:text-[44px]
            text-[42px]
            font-semibold
            leading-[130%]
            max-lg:mb-4
            whitespace-nowrap
            tracking-tight
          "
        >
          {{ t('sections.team.title') }}
        </h2>
        <p class="min-lg:text-[17px] text-[16px] text-gray-70 max-w-[600px]">
          {{ t('sections.team.description') }}
        </p>
      </div>

      <ul class="member-list min-md:pt-10 min-md:border-t border-t-gray-30">
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
  @apply inline-flex grid-flow-col max-lg:grid-flow-row w-full grid-cols-1 min-md:grid-cols-2 max-lg:flex-wrap min-lg:gap-0 min-md:gap-y-6;

  & > * {
    @apply min-lg:flex-1 min-md:flex-[1_0_50%] flex-[1_0_100%] border-gray-30 max-md:py-5;
  }

  & > :not(:last-child) {
    @apply min-lg:border-r max-md:border-b;
  }

  @media screen and (max-width: 1200px) {
    & > :last-child > div {
      @apply pr-0;
    }

    & > :first-child > div {
      @apply pl-0;
    }
  }
}
</style>
