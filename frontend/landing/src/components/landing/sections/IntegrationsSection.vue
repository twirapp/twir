<template>
  <section
    class="
      relative
      min-lg:pt-[130px] min-lg:pb-[100px] min-xl:px-12 min-lg:px-16
      bg-black-15
      pt-32
      pb-16
      min-md:px-10 min-sm:px-8
      px-5
    "
  >
    <div class="container max-w-[1180px] relative">
      <div class="flex flex-col justify-between items-center relative z-0">
        <Cube class="absolute left-0 right-0 mx-auto -z-[1] -top-[80px]" />
        <h2
          class="
            min-xl:text-[48px] min-sm:text-[44px]
            text-[36px]
            tracking-tight
            font-semibold
            max-w-[600px] max-sm:px-4
            leading-[125%]
            text-center
            mx-auto
            mb-9
          "
        >
          {{ t('sections.integrations.title') }}
        </h2>
        <ul class="flex flex-wrap justify-center">
          <li
            v-for="(item, index) in integrationLogos"
            :key="index"
            class="min-md:px-4 px-2 min-md:py-[11px] py-2"
          >
            <component :is="item" class="min-md:h-[48px] w-auto h-[36px]" />
          </li>
        </ul>
      </div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import Cube from '@/assets/Cube.svg?component';
import DonationAlerts from '@/assets/logos/DonationAlerts.svg?component';
import Dota2 from '@/assets/logos/Dota2.svg?component';
import Faceit from '@/assets/logos/Faceit.svg?component';
import LastFm from '@/assets/logos/LastFm.svg?component';
import OBS from '@/assets/logos/OBS.svg?component';
import Spotify from '@/assets/logos/Spotify.svg?component';
import StreamElements from '@/assets/logos/StreamElements.svg?component';
import Streamlabs from '@/assets/logos/Streamlabs.svg?component';
import TwitchSoundtrack from '@/assets/logos/TwitchSoundtrack.svg?component';
import VK from '@/assets/logos/VK.svg?component';
import { useTranslation } from '@/services/locale';

const integrationLogos = [
  Spotify,
  TwitchSoundtrack,
  Dota2,
  VK,
  DonationAlerts,
  LastFm,
  StreamElements,
  Streamlabs,
  OBS,
  Faceit,
];

const { t } = useTranslation<'landing'>();
</script>
